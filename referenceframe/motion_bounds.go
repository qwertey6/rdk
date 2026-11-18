package referenceframe

import (
	"math"

	"github.com/golang/geo/r3"
	"gonum.org/v1/gonum/num/quat"

	"go.viam.com/rdk/spatialmath"
)

// MotionBounds bounds the speed of every robot geometry point with respect to
// the edge interpolation fraction, not elapsed time. The immutable tree covers
// all branches, mimics, attached tools and multiple moving components.
// Unsupported frames/shapes return nil and must use dense collision checking.
type MotionBounds struct {
	nativePoses   bool
	nodes         []motionBoundNode
	shapes        []motionBoundShape
	geometryNodes []int
	pairs         []motionBoundPair
}

// HasNativePoses reports whether all fixed frame and geometry poses use native
// dual-quaternion storage. It does not certify collision freedom or permit
// mutating the snapshot while a prepared evaluator is in use.
func (p *MotionBounds) HasNativePoses() bool { return p != nil && p.nativePoses }

type motionBoundPair struct{ a, b, common int }

type motionBoundInput struct {
	frame       string
	index       int
	scale, bias float64
}

type motionBoundNode struct {
	parent      int
	kind        byte
	reach       float64
	input       motionBoundInput
	axis        r3.Vector
	orientation quat.Number
}

type motionBoundShape struct {
	node   int
	radius float64
}

// PrepareMotionBounds prepares a conservative geometric certificate for native
// rigid frames. Rebuild it when the frame-system/model snapshot changes.
func (fs *FrameSystem) PrepareMotionBounds(stationary ...spatialmath.Geometry) *MotionBounds {
	p := &MotionBounds{nodes: []motionBoundNode{{parent: -1}}, nativePoses: true}
	var addFS func(*FrameSystem, int, func(string) (motionBoundInput, bool)) (map[string]int, bool)
	var addFrame func(Frame, int, motionBoundInput) (int, bool)
	addShapes := func(node int, geometries []spatialmath.Geometry) bool {
		for _, g := range geometries {
			center, radius, ok := spatialmath.BoundingSphere(g)
			extent := center.Norm() + radius
			if !ok || !finiteBound(extent) || radius < 0 || extent > 1e5 {
				return false
			}
			p.shapes = append(p.shapes, motionBoundShape{node, extent})
			// Geometry pose composition can also approach an identity rotation.
			pose := g.Pose()
			if _, ok := pose.(*spatialmath.DualQuaternion); !ok {
				p.nativePoses = false
			}
			q := pose.Orientation().Quaternion()
			if !boundedQuaternion(q) {
				return false
			}
			p.nodes = append(p.nodes, motionBoundNode{parent: node, kind: 'f', orientation: q})
			p.geometryNodes = append(p.geometryNodes, len(p.nodes)-1)
		}
		return true
	}
	addFrame = func(frame Frame, parent int, input motionBoundInput) (int, bool) {
		if named, ok := frame.(*namedFrame); ok {
			return addFrame(named.Frame, parent, input)
		}
		node := motionBoundNode{parent: parent, input: input, orientation: quat.Number{Real: 1}}
		shapeNode := parent
		var geometries []spatialmath.Geometry
		switch f := frame.(type) {
		case *SimpleModel:
			bindings := map[string]motionBoundInput{}
			for _, meta := range f.inputSchema.metas {
				binding := input
				binding.index += meta.offset
				bindings[meta.frameName] = binding
			}
			for name, mimic := range f.mimicMappings {
				// A mimic is scalar. Nested non-scalar affine bindings are unsupported.
				if input.scale != 1 || input.bias != 0 {
					return 0, false
				}
				bindings[name] = motionBoundInput{
					input.frame, input.index + mimic.sourceInputIdx,
					mimic.valueMultiplier, mimic.valueOffset,
				}
			}
			indices, ok := addFS(f.internalFS, parent, func(name string) (motionBoundInput, bool) {
				b, exists := bindings[name]
				return b, exists || len(f.internalFS.Frame(name).DoF()) == 0
			})
			out, exists := indices[f.primaryOutputFrame]
			return out, ok && exists
		case *tailGeometryStaticFrame:
			if _, ok := f.transform.(*spatialmath.DualQuaternion); !ok {
				p.nativePoses = false
			}
			node.kind, node.reach = 'f', f.transform.Point().Norm()
			node.orientation = f.transform.Orientation().Quaternion()
			geometries, shapeNode = f.geometries, len(p.nodes)
		case *staticFrame:
			if _, ok := f.transform.(*spatialmath.DualQuaternion); !ok {
				p.nativePoses = false
			}
			node.kind, node.reach, geometries = 'f', f.transform.Point().Norm(), f.geometries
			node.orientation = f.transform.Orientation().Quaternion()
		case *rotationalFrame:
			if !finiteBound(f.rotAxis.Norm()) || math.Abs(f.rotAxis.Norm()-1) > 1e-12 {
				return 0, false
			}
			node.kind, node.axis = 'r', f.rotAxis
		case *translationalFrame:
			node.kind, node.reach = 't', f.transAxis.Norm()
			geometries, shapeNode = f.geometries, len(p.nodes)
		default:
			return 0, false
		}
		if !finiteBound(node.reach) || node.reach > 1e5 || !boundedQuaternion(node.orientation) {
			return 0, false
		}
		p.nodes = append(p.nodes, node)
		index := len(p.nodes) - 1
		return index, addShapes(shapeNode, geometries)
	}
	addFS = func(tree *FrameSystem, root int, binding func(string) (motionBoundInput, bool)) (map[string]int, bool) {
		indices := map[string]int{World: root}
		var visit func(string) (int, bool)
		visit = func(name string) (int, bool) {
			if i, ok := indices[name]; ok {
				return i, true
			}
			frame := tree.Frame(name)
			if frame == nil {
				return 0, false
			}
			parent, ok := visit(tree.parents[name])
			if !ok {
				return 0, false
			}
			b, ok := binding(name)
			if !ok {
				return 0, false
			}
			i, ok := addFrame(frame, parent, b)
			if ok {
				indices[name] = i
			}
			return i, ok
		}
		for _, name := range tree.FrameNames() {
			if _, ok := visit(name); !ok {
				return nil, false
			}
		}
		return indices, true
	}
	if _, ok := addFS(fs, 0, func(name string) (motionBoundInput, bool) {
		return motionBoundInput{frame: name, scale: 1}, true
	}); !ok {
		return nil
	}
	if !addShapes(0, stationary) || len(p.geometryNodes) > 128 {
		return nil
	}
	depths := make([]int, len(p.nodes))
	for i := 1; i < len(p.nodes); i++ {
		depths[i] = depths[p.nodes[i].parent] + 1
	}
	for i, a := range p.geometryNodes {
		for _, b := range p.geometryNodes[i+1:] {
			x, y := a, b
			for x != y {
				if depths[x] >= depths[y] {
					x = p.nodes[x].parent
				} else {
					y = p.nodes[y].parent
				}
			}
			p.pairs = append(p.pairs, motionBoundPair{a, b, x})
		}
	}
	return p
}

func finiteBound(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// MaxPointSpeed returns a rigid-body motion bound only when native quaternion
// arithmetic stays in a stable regime over this entire interpolation. False
// requires dense checking; it is not a planning failure.
func (p *MotionBounds) MaxPointSpeed(from, to *LinearInputs) (float64, bool) {
	if !p.InterpolationStable(from, to) {
		return 0, false
	}
	return p.maxIdealPointSpeed(from, to)
}

// maxIdealPointSpeed bounds ideal rigid-body displacement by speed*abs(fractionB-fractionA).
// A native collision certificate must also require InterpolationStable: native
// quaternion arithmetic has a discontinuity near identity rotations.
// For a child origin: v <= parent.v + parent.omega*maxTranslation +
// translationSpeed. A geometry point adds omega*its local enclosing radius.
// Triangle inequalities and max endpoint prismatic extension hold over the
// entire linear joint interpolation, including rotations greater than pi.
func (p *MotionBounds) maxIdealPointSpeed(from, to *LinearInputs) (float64, bool) {
	if p == nil {
		return 0, false
	}
	type speed struct{ linear, angular float64 }
	speeds := make([]speed, len(p.nodes))
	for i := 1; i < len(p.nodes); i++ {
		n := p.nodes[i]
		v := speeds[n.parent]
		reach, translation, rotation := n.reach, 0., 0.
		if n.kind == 'r' || n.kind == 't' {
			a, b := from.Get(n.input.frame), to.Get(n.input.frame)
			if n.input.index < 0 || n.input.index >= len(a) || n.input.index >= len(b) {
				return 0, false
			}
			x, y := a[n.input.index]*n.input.scale+n.input.bias, b[n.input.index]*n.input.scale+n.input.bias
			if !finiteBound(x) || !finiteBound(y) {
				return 0, false
			}
			if n.kind == 'r' {
				rotation = math.Abs(y - x)
			} else {
				reach = max(math.Abs(x), math.Abs(y)) * n.reach
				translation = math.Abs(y-x) * n.reach
			}
		}
		speeds[i] = speed{v.linear + v.angular*reach + translation, v.angular + rotation}
	}
	bound := 0.
	for _, shape := range p.shapes {
		v := speeds[shape.node]
		bound = max(bound, v.linear+v.angular*shape.radius)
	}
	// Inflate the analytic bound for floating-point composition and subtraction.
	bound = bound*(1+1e-10) + 1e-9
	return bound, finiteBound(bound)
}
