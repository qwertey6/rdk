package referenceframe

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/golang/geo/r3"
	"gonum.org/v1/gonum/num/quat"

	"go.viam.com/rdk/spatialmath"
)

// PreparedGeometries retains immutable local geometry and evaluator storage for
// one frame-system snapshot. Construct again after any model/geometry change.
// It supports custom frames through their ordinary Geometries method.
type PreparedGeometries struct {
	fs      *FrameSystem
	frames  []preparedGeometryFrame
	indices map[string]int
	pool    sync.Pool
}

type preparedGeometryFrame struct {
	name     string
	frame    Frame
	parent   string
	local    []spatialmath.Geometry
	fixed    bool
	model    *SimpleModel
	children *PreparedGeometries
}

// PrepareGeometries resolves frame topology and fixed local shapes once.
func (fs *FrameSystem) PrepareGeometries(wanted map[string]bool) *PreparedGeometries {
	p := &PreparedGeometries{fs: fs, indices: map[string]int{}}
	for _, name := range fs.FrameNames() {
		if wanted != nil && !wanted[name] {
			continue
		}
		frame := fs.Frame(name)
		entry := preparedGeometryFrame{name: name, frame: frame, parent: fs.parents[name]}
		inner := frame
		for {
			named, ok := inner.(*namedFrame)
			if !ok {
				break
			}
			inner = named.Frame
		}
		switch f := inner.(type) {
		case *SimpleModel:
			entry.model, entry.children = f, f.internalFS.PrepareGeometries(nil)
		case *staticFrame, *tailGeometryStaticFrame:
			if gif, err := frame.Geometries(nil); err == nil && gif != nil {
				entry.local, entry.fixed = gif.Geometries(), true
			}
		}
		p.indices[name] = len(p.frames)
		p.frames = append(p.frames, entry)
	}
	p.pool.New = func() any {
		return &GeometryEvaluation{prepared: p, frames: make([]geometryFrameEvaluation, len(p.frames))}
	}
	return p
}

type geometryFrameEvaluation struct {
	parentPose spatialmath.DualQuaternion
	done       bool
	sources    []spatialmath.Geometry
	buffers    []*spatialmath.GeometryTransformBuffer
	geometries []spatialmath.Geometry
	child      *GeometryEvaluation
}

// GeometryEvaluation is private scratch for a single state evaluation.
// Subset results borrow its storage; Release invalidates all borrowed geometry.
type GeometryEvaluation struct {
	prepared *PreparedGeometries
	inputs   *LinearInputs
	fk       *lazyFrameToWorld
	frames   []geometryFrameEvaluation
	output   []spatialmath.Geometry
}

// Evaluate acquires independent scratch and a memoized FK pass.
func (p *PreparedGeometries) Evaluate(inputs *LinearInputs) *GeometryEvaluation {
	e := p.pool.Get().(*GeometryEvaluation)
	e.inputs, e.fk = inputs, newLazyFrameToWorld(p.fs, inputs)
	return e
}

// PoseQT returns a frame's world pose using the same memo as Subset. Callers
// sharing it with other computations must use immutable native frames and
// must not mutate the configuration until Release.
func (e *GeometryEvaluation) PoseQT(name string) (quat.Number, r3.Vector, error) {
	dq, err := e.fk.get(name)
	if err != nil {
		return quat.Number{}, r3.Vector{}, err
	}
	d := spatialmath.DualQuaternion{Number: dq}
	return dq.Real, d.Point(), nil
}

// Release recycles private scratch after all collision checks have finished.
func (e *GeometryEvaluation) Release() {
	for i := range e.frames {
		f := &e.frames[i]
		f.done = false
		if f.child != nil {
			f.child.Release()
			f.child = nil
		}
	}
	e.inputs = nil
	e.fk.release()
	e.fk = nil
	e.output = e.output[:0]
	e.prepared.pool.Put(e)
}

func (e *GeometryEvaluation) frame(i int) ([]spatialmath.Geometry, error) {
	f, work := &e.prepared.frames[i], &e.frames[i]
	if work.done {
		return work.geometries, nil
	}
	inputs := e.inputs.Get(f.frame.Name())
	if len(inputs) != len(f.frame.DoF()) {
		return nil, NewIncorrectDoFError(len(inputs), len(f.frame.DoF()))
	}
	local := f.local
	if f.model != nil {
		li, err := f.model.toLinearInputs(inputs)
		if err != nil {
			return nil, err
		}
		if len(f.model.mimicMappings) > 0 {
			li.forkSchema()
		}
		for name, mimic := range f.model.mimicMappings {
			li.Put(name, []Input{inputs[mimic.sourceInputIdx]*mimic.valueMultiplier + mimic.valueOffset})
		}
		work.child = f.children.Evaluate(li)
		local, err = work.child.Subset(nil)
		if err != nil {
			return nil, err
		}
	} else if !f.fixed {
		gif, err := f.frame.Geometries(inputs)
		if err != nil {
			return nil, err
		}
		local = gif.Geometries()
	}
	if cap(work.geometries) < len(local) {
		work.geometries = make([]spatialmath.Geometry, len(local))
		work.buffers = make([]*spatialmath.GeometryTransformBuffer, len(local))
		work.sources = make([]spatialmath.Geometry, len(local))
	} else {
		work.geometries = work.geometries[:len(local)]
		work.buffers = work.buffers[:len(local)]
		work.sources = work.sources[:len(local)]
	}
	if len(local) > 0 {
		dq, err := e.fk.get(f.parent)
		if err != nil {
			return nil, err
		}
		work.parentPose.Number = dq
		for j, g := range local {
			// Custom frames can return a different shape on every call. Rebuild
			// their buffer whenever the source changes, keeping exact semantics.
			if !reflect.TypeOf(g).Comparable() || work.sources[j] != g {
				work.sources[j], work.buffers[j] = g, spatialmath.NewGeometryTransformBuffer(g)
			}
			out := work.buffers[j].Transform(&work.parentPose)
			if f.model != nil {
				out.SetLabel(f.model.name + ":" + g.Label())
			}
			work.geometries[j] = out
		}
	}
	work.done = true
	return work.geometries, nil
}

// Subset returns world geometry in stable frame order. nil selects all frames.
func (e *GeometryEvaluation) Subset(wanted map[string]bool) ([]spatialmath.Geometry, error) {
	e.output = e.output[:0]
	// nil selects all frames below; only named subsets need membership validation.
	if wanted != nil { //nolint:staticcheck // Keep the nil-as-all case explicit.
		for name := range wanted {
			if _, ok := e.prepared.indices[name]; !ok {
				return nil, fmt.Errorf("unprepared geometry frame %q", name)
			}
		}
	}
	for i, frame := range e.prepared.frames {
		if wanted != nil && !wanted[frame.name] {
			continue
		}
		geometries, err := e.frame(i)
		if err != nil {
			return nil, err
		}
		e.output = append(e.output, geometries...)
	}
	return e.output, nil
}
