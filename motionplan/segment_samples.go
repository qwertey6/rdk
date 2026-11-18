package motionplan

import (
	"context"
	"math"

	"go.viam.com/rdk/referenceframe"
)

// SegmentSamples visits precisely the samples returned by InterpolateSegmentFS.
// At returns borrowed storage; Copy a result that must survive the next call.
// Custom interpolation retains eager evaluation and its error ordering.
type SegmentSamples struct {
	steps    int
	prepared *referenceframe.PreparedInterpolation
	eager    []*referenceframe.LinearInputs
}

// SegmentSamples constructs an independent sampler for an edge.
func (c *ConstraintChecker) SegmentSamples(ci *SegmentFS, resolution float64) (*SegmentSamples, error) {
	if c.streaming {
		if p, ok := referenceframe.PrepareInterpolation(ci.FS, ci.StartConfiguration, ci.EndConfiguration); ok {
			steps, err := segmentStepCount(ci, resolution)
			if err != nil {
				return nil, err
			}
			return &SegmentSamples{steps: steps, prepared: p}, nil
		}
	}
	inputs, err := InterpolateSegmentFS(ci, resolution)
	if err != nil {
		return nil, err
	}
	return &SegmentSamples{steps: len(inputs) - 1, eager: inputs}, nil
}

// Len includes both endpoints.
func (s *SegmentSamples) Len() int { return s.steps + 1 }

// At returns the ith sample, in [0, Len()).
func (s *SegmentSamples) At(i int) *referenceframe.LinearInputs {
	if s.prepared != nil {
		return s.prepared.At(float64(i) / float64(s.steps))
	}
	return s.eager[i]
}

// SegmentSampleCount returns the full-resolution validation work of an edge.
func SegmentSampleCount(ci *SegmentFS, resolution float64) (int, error) {
	n, err := segmentStepCount(ci, resolution)
	return n + 1, err
}

// SegmentValidation keeps a private clearance certificate for one edge. It
// skips collision work only while every geometry pair is provably separated;
// topological and cable constraints still run at every original sample.
type SegmentValidation struct {
	Samples                 *SegmentSamples
	checker                 *ConstraintChecker
	segment                 *SegmentFS
	state                   StateFS
	anchor                  int
	clearance, relativeStep float64
	certified               bool
}

// SegmentValidation constructs an edge validator with independent scratch.
func (c *ConstraintChecker) SegmentValidation(ci *SegmentFS, resolution float64) (*SegmentValidation, error) {
	samples, err := c.SegmentSamples(ci, resolution)
	if err != nil {
		return nil, err
	}
	v := &SegmentValidation{Samples: samples, checker: c, segment: ci}
	if c.motionBounds != nil {
		if speed, ok := c.motionBounds.MaxPointSpeed(ci.StartConfiguration, ci.EndConfiguration); ok {
			// Either side of a self/inter-arm pair can move by speed*dt.
			v.relativeStep = 2 * speed / float64(samples.steps)
			v.certified = true
		}
	}
	return v, nil
}

// Check returns a borrowed configuration and a conservative clearance. The
// floor additionally prevents skipping close-waypoint classification. Copy a
// configuration retained beyond the next Check or Samples.At call.
func (v *SegmentValidation) Check(ctx context.Context, i int, floor float64) (*referenceframe.LinearInputs, float64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	q := v.Samples.At(i)
	if !v.checker.nativeCallbacks {
		q = q.Copy()
	}
	// Clear all per-state caches before reusing scratch.
	v.state = StateFS{FS: v.segment.FS, Configuration: q}
	state := &v.state
	if !v.checker.nativeCallbacks || !v.checker.streaming {
		state = &StateFS{FS: v.segment.FS, Configuration: q}
	}
	lower := v.clearance - math.Abs(float64(i-v.anchor))*v.relativeStep
	if v.certified && lower > max(floor, v.checker.collisionBufferMM)+1e-2 {
		return q, lower, v.checker.CheckStateFSTopoOnly(state)
	}
	distance, err := v.checker.CheckStateFSConstraints(ctx, state)
	if err == nil {
		v.anchor, v.clearance = i, distance
	} else {
		v.clearance = 0
	}
	return q, distance, err
}
