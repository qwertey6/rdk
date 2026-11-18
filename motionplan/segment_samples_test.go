package motionplan

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func sampleTestSegment(t *testing.T) *SegmentFS {
	t.Helper()
	r, err := referenceframe.NewRotationalFrame("rotate", spatialmath.R4AA{RZ: 1}, referenceframe.Limit{Min: -10, Max: 10})
	test.That(t, err, test.ShouldBeNil)
	p, err := referenceframe.NewTranslationalFrame("slide", r3.Vector{X: 1}, referenceframe.Limit{Min: -100, Max: 100})
	test.That(t, err, test.ShouldBeNil)
	m, err := referenceframe.NewSerialModel("mixed", []referenceframe.Frame{r, p})
	test.That(t, err, test.ShouldBeNil)
	fs := referenceframe.NewEmptyFrameSystem("samples")
	test.That(t, fs.AddFrame(m, fs.World()), test.ShouldBeNil)
	return &SegmentFS{
		FS:                 fs,
		StartConfiguration: referenceframe.FrameSystemInputs{"mixed": {-0.11, 7}}.ToLinearInputs(),
		EndConfiguration:   referenceframe.FrameSystemInputs{"mixed": {0.15, -11}}.ToLinearInputs(),
	}
}

func TestSegmentSamplesParity(t *testing.T) {
	segment := sampleTestSegment(t)
	checker := NewEmptyConstraintChecker(logging.NewTestLogger(t))
	samples, err := checker.SegmentSamples(segment, 1)
	test.That(t, err, test.ShouldBeNil)
	eager, err := InterpolateSegmentFS(segment, 1)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, samples.Len(), test.ShouldEqual, len(eager))
	for i, want := range eager {
		test.That(t, samples.At(i).Get("mixed"), test.ShouldResemble, want.Get("mixed"))
	}
	// Random access and an owned saved sample survive scratch reuse.
	saved := samples.At(5).Copy()
	samples.At(0)
	test.That(t, saved.Get("mixed"), test.ShouldResemble, eager[5].Get("mixed"))
}

func TestStreamingPrefixAndCallbackOwnership(t *testing.T) {
	segment := sampleTestSegment(t)
	eager, err := InterpolateSegmentFS(segment, 1)
	test.That(t, err, test.ShouldBeNil)
	checker := NewEmptyConstraintChecker(logging.NewTestLogger(t))
	var retained []*StateFS
	checker.SetCollisionConstraints(CollisionConstraints{Obstacle: func(s *StateFS) (float64, error) {
		retained = append(retained, s)
		if len(retained) == 4 {
			return 0, fmt.Errorf("blocked fourth sample")
		}
		return 1, nil
	}})
	prefix, err := checker.CheckStateConstraintsAcrossSegmentFS(context.Background(), segment, 1, true)
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, prefix.EndConfiguration.Get("mixed"), test.ShouldResemble, eager[2].Get("mixed"))
	for i, s := range retained {
		test.That(t, s.Configuration.Get("mixed"), test.ShouldResemble, eager[i].Get("mixed"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = checker.CheckStateConstraintsAcrossSegmentFS(ctx, segment, 1, true)
	test.That(t, err, test.ShouldEqual, context.Canceled)
}

type lateInterpolationError struct{ referenceframe.Frame }

func (f lateInterpolationError) Interpolate(a, b []referenceframe.Input, fraction float64) ([]referenceframe.Input, error) {
	if fraction > .5 {
		return nil, fmt.Errorf("late interpolation failure")
	}
	return f.Frame.Interpolate(a, b, fraction)
}

func TestStreamingCustomInterpolationErrorOrder(t *testing.T) {
	segment := sampleTestSegment(t)
	fs := referenceframe.NewEmptyFrameSystem("custom")
	test.That(t, fs.AddFrame(lateInterpolationError{segment.FS.Frame("mixed")}, fs.World()), test.ShouldBeNil)
	segment.FS = fs
	checker := NewEmptyConstraintChecker(logging.NewTestLogger(t))
	checks := 0
	checker.SetCollisionConstraints(CollisionConstraints{Obstacle: func(*StateFS) (float64, error) {
		checks++
		return 0, fmt.Errorf("earlier collision")
	}})
	_, err := checker.CheckStateConstraintsAcrossSegmentFS(context.Background(), segment, 1, true)
	test.That(t, err.Error(), test.ShouldEqual, "late interpolation failure")
	test.That(t, checks, test.ShouldEqual, 0)
}

func TestStreamingCancellationDuringFinalState(t *testing.T) {
	segment := sampleTestSegment(t)
	n, err := SegmentSampleCount(segment, 1)
	test.That(t, err, test.ShouldBeNil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	checks := 0
	checker := NewEmptyConstraintChecker(logging.NewTestLogger(t))
	checker.SetCollisionConstraints(CollisionConstraints{Obstacle: func(*StateFS) (float64, error) {
		checks++
		if checks == n {
			cancel()
		}
		return 1, nil
	}})
	prefix, err := checker.CheckStateConstraintsAcrossSegmentFS(ctx, segment, 1, true)
	test.That(t, err, test.ShouldEqual, context.Canceled)
	test.That(t, prefix, test.ShouldBeNil)
}
