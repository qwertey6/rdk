package motionplan

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func TestShortcutSampleCoverageAndRejection(t *testing.T) {
	segment := sampleTestSegment(t)
	for _, resolution := range []float64{20, 1, .7, .3} {
		for _, checkFinal := range []bool{false, true} {
			t.Run(fmt.Sprint(resolution, checkFinal), func(t *testing.T) {
				want, err := InterpolateSegmentFS(segment, resolution)
				test.That(t, err, test.ShouldBeNil)
				if !checkFinal {
					want = want[:len(want)-1]
				}
				checker := NewEmptyConstraintChecker(logging.NewTestLogger(t))
				// This test-only pure callback observes the native traversal without
				// changing constraints or retaining borrowed state.
				checker.nativeCallbacks = true
				checker.motionBounds = segment.FS.PrepareMotionBounds()
				test.That(t, checker.motionBounds, test.ShouldNotBeNil)
				var visited []float64
				checker.topoConstraint = func(s *StateFS) error {
					visited = append(visited, s.Configuration.Get("mixed")[0])
					return nil
				}
				test.That(t, checker.CheckShortcutFS(context.Background(), segment, resolution, checkFinal), test.ShouldBeNil)
				slices.Sort(visited)
				values := make([]float64, len(want))
				for i, q := range want {
					values[i] = q.Get("mixed")[0]
				}
				test.That(t, visited, test.ShouldResemble, values)
				// An invalid state anywhere in the required grid must reject,
				// including points missed by the early subdivision probes.
				for _, q := range want {
					bad := q.Get("mixed")[0]
					checker.topoConstraint = func(s *StateFS) error {
						if s.Configuration.Get("mixed")[0] == bad {
							return fmt.Errorf("blocked sample")
						}
						return nil
					}
					test.That(t, checker.CheckShortcutFS(context.Background(), segment, resolution, checkFinal), test.ShouldNotBeNil)
				}
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				test.That(t, checker.CheckShortcutFS(ctx, segment, resolution, checkFinal), test.ShouldEqual, context.Canceled)
			})
		}
	}
}

func TestShortcutCustomOrderAndCancellation(t *testing.T) {
	segment := sampleTestSegment(t)
	want, err := InterpolateSegmentFS(segment, .3)
	test.That(t, err, test.ShouldBeNil)
	checker := NewEmptyConstraintChecker(logging.NewTestLogger(t))
	var seen []*referenceframe.LinearInputs
	checker.SetCollisionConstraints(CollisionConstraints{Obstacle: func(s *StateFS) (float64, error) {
		seen = append(seen, s.Configuration)
		return 1, nil
	}})
	test.That(t, checker.CheckShortcutFS(context.Background(), segment, .3, true), test.ShouldBeNil)
	for i, q := range seen {
		test.That(t, q.Get("mixed"), test.ShouldResemble, want[i].Get("mixed"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	checks := 0
	checker.SetCollisionConstraints(CollisionConstraints{Obstacle: func(*StateFS) (float64, error) {
		checks++
		if checks == len(want) {
			cancel()
		}
		return 1, nil
	}})
	test.That(t, checker.CheckShortcutFS(ctx, segment, .3, true), test.ShouldEqual, context.Canceled)
	fs := referenceframe.NewEmptyFrameSystem("custom")
	test.That(t, fs.AddFrame(lateInterpolationError{segment.FS.Frame("mixed")}, fs.World()), test.ShouldBeNil)
	segment.FS = fs
	checks = 0
	err = checker.CheckShortcutFS(context.Background(), segment, .3, true)
	test.That(t, err.Error(), test.ShouldEqual, "late interpolation failure")
	test.That(t, checks, test.ShouldEqual, 0)
}

type customStaticShortcutFrame struct{ referenceframe.Frame }

func TestShortcutCustomZeroDoFOrder(t *testing.T) {
	segment := sampleTestSegment(t)
	// A custom zero-DoF frame has no entry in the interpolation input map.
	// Native input interpolation alone must not authorize reordering its calls.
	base, err := referenceframe.NewStaticFrame("custom-static", spatialmath.NewZeroPose())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, segment.FS.AddFrame(customStaticShortcutFrame{base}, segment.FS.World()), test.ShouldBeNil)
	checker := NewEmptyConstraintChecker(logging.NewTestLogger(t))
	checker.nativeCallbacks = true
	checker.motionBounds = segment.FS.PrepareMotionBounds()
	test.That(t, checker.motionBounds, test.ShouldBeNil)
	var visited []float64
	checker.topoConstraint = func(s *StateFS) error {
		visited = append(visited, s.Configuration.Get("mixed")[0])
		return nil
	}
	want, err := InterpolateSegmentFS(segment, .3)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, checker.CheckShortcutFS(context.Background(), segment, .3, true), test.ShouldBeNil)
	test.That(t, len(visited), test.ShouldEqual, len(want))
	for i, q := range want {
		test.That(t, visited[i], test.ShouldEqual, q.Get("mixed")[0])
	}
}
