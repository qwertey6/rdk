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

func shortcutGeometryFixture(t *testing.T, obstacleY float64) (*ConstraintChecker, *SegmentFS) {
	t.Helper()
	fs := referenceframe.NewEmptyFrameSystem("local-bounds")
	joint, err := referenceframe.NewRotationalFrame("joint", spatialmath.R4AA{RZ: 1}, referenceframe.Limit{Min: -2, Max: 2})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(joint, fs.World()), test.ShouldBeNil)
	g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000}), 1, "tip")
	test.That(t, err, test.ShouldBeNil)
	tip, err := referenceframe.NewStaticFrameWithGeometry("tip", spatialmath.NewZeroPose(), g)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(tip, joint), test.ShouldBeNil)
	segment := &SegmentFS{
		FS:                 fs,
		StartConfiguration: referenceframe.FrameSystemInputs{"joint": {-.4}}.ToLinearInputs(),
		EndConfiguration:   referenceframe.FrameSystemInputs{"joint": {.4}}.ToLinearInputs(),
	}
	obstacle, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000, Y: obstacleY}), 1, "obstacle")
	test.That(t, err, test.ShouldBeNil)
	geoms, err := referenceframe.FrameSystemGeometriesLinearInputs(fs, segment.StartConfiguration)
	test.That(t, err, test.ShouldBeNil)
	checker, err := NewConstraintChecker(defaultCollisionBufferMM, nil, nil, nil, fs,
		geoms["tip"].Geometries(), nil, map[string]bool{"tip": true}, segment.StartConfiguration,
		referenceframe.NewGeometriesInFrame(referenceframe.World, []spatialmath.Geometry{obstacle}),
		logging.NewTestLogger(t), NewCollisionCache())
	test.That(t, err, test.ShouldBeNil)
	return checker, segment
}

func TestShortcutOriginalGeometryOracle(t *testing.T) {
	for _, obstacleY := range []float64{0, 1000} {
		t.Run(fmt.Sprint(obstacleY), func(t *testing.T) {
			checker, segment := shortcutGeometryFixture(t, obstacleY)
			validation, err := checker.SegmentValidation(segment, 2)
			test.That(t, err, test.ShouldBeNil)
			topoChecks := 0
			checker.topoConstraint = func(*StateFS) error { topoChecks++; return nil }
			// Independent original materialization at every sample, with no
			// interval certificate or reused geometry, is the verdict oracle.
			obstacle, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000, Y: obstacleY}), 1, "obstacle")
			test.That(t, err, test.ShouldBeNil)
			collisions := 0
			for _, reverse := range []bool{false, true} {
				for j := 0; j < validation.Samples.Len(); j++ {
					i := j
					if reverse {
						i = validation.Samples.Len() - 1 - j
					}
					q, _, got := validation.Check(context.Background(), i, .1)
					geoms, err := referenceframe.FrameSystemGeometriesLinearInputs(segment.FS, q.Copy())
					test.That(t, err, test.ShouldBeNil)
					collision, _, err := geoms["tip"].Geometries()[0].CollidesWith(obstacle, defaultCollisionBufferMM)
					test.That(t, err, test.ShouldBeNil)
					test.That(t, got == nil, test.ShouldEqual, !collision)
					if collision {
						collisions++
					}
				}
			}
			test.That(t, topoChecks, test.ShouldEqual, 2*validation.Samples.Len())
			if obstacleY == 0 {
				test.That(t, collisions, test.ShouldBeGreaterThan, 0)
			} else {
				test.That(t, collisions, test.ShouldEqual, 0)
			}
			err = checker.CheckShortcutFS(context.Background(), segment, 2, true)
			test.That(t, err == nil, test.ShouldEqual, collisions == 0)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			_, _, err = validation.Check(ctx, 0, 0)
			test.That(t, err, test.ShouldEqual, context.Canceled)
		})
	}
}

func TestShortcutIdentityFallback(t *testing.T) {
	checker, segment := shortcutGeometryFixture(t, 0)
	segment.StartConfiguration.Get("joint")[0] = 1.9e-8
	segment.EndConfiguration.Get("joint")[0] = 2.2e-8
	// Force many original samples around the known native discontinuity.
	validation, err := checker.SegmentValidation(segment, 1e-8)
	test.That(t, err, test.ShouldBeNil)
	for i := 0; i < validation.Samples.Len(); i++ {
		_, _, _ = validation.Check(context.Background(), i, 0)
	}
	_, expected := checker.CheckStateConstraintsAcrossSegmentFS(context.Background(), segment, 1e-8, true)
	got := checker.CheckShortcutFS(context.Background(), segment, 1e-8, true)
	test.That(t, got == nil, test.ShouldEqual, expected == nil)
	test.That(t, validation.certified, test.ShouldBeFalse)
}

func TestShortcutBothArmsApproach(t *testing.T) {
	fs := referenceframe.NewEmptyFrameSystem("two-arms")
	for _, name := range []string{"left", "right"} {
		joint, err := referenceframe.NewRotationalFrame(name, spatialmath.R4AA{RZ: 1}, referenceframe.Limit{Min: -2, Max: 2})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, fs.AddFrame(joint, fs.World()), test.ShouldBeNil)
		g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000}), 1, name+"-tip")
		test.That(t, err, test.ShouldBeNil)
		tip, err := referenceframe.NewStaticFrameWithGeometry(name+"-tip", spatialmath.NewZeroPose(), g)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, fs.AddFrame(tip, joint), test.ShouldBeNil)
	}
	segment := &SegmentFS{
		FS:                 fs,
		StartConfiguration: referenceframe.FrameSystemInputs{"left": {-.4}, "right": {.4}}.ToLinearInputs(),
		EndConfiguration:   referenceframe.FrameSystemInputs{"left": {.4}, "right": {-.4}}.ToLinearInputs(),
	}
	geoms, err := referenceframe.FrameSystemGeometriesLinearInputs(fs, segment.StartConfiguration)
	test.That(t, err, test.ShouldBeNil)
	moving := []spatialmath.Geometry{geoms["left-tip"].Geometries()[0], geoms["right-tip"].Geometries()[0]}
	checker, err := NewConstraintChecker(defaultCollisionBufferMM, nil, nil, nil, fs, moving, nil,
		map[string]bool{"left-tip": true, "right-tip": true}, segment.StartConfiguration, nil,
		logging.NewTestLogger(t), NewCollisionCache())
	test.That(t, err, test.ShouldBeNil)
	validation, err := checker.SegmentValidation(segment, 2)
	test.That(t, err, test.ShouldBeNil)
	collisions := 0
	for i := 0; i < validation.Samples.Len(); i++ {
		q, _, got := validation.Check(context.Background(), i, 0)
		geoms, err := referenceframe.FrameSystemGeometriesLinearInputs(fs, q.Copy())
		test.That(t, err, test.ShouldBeNil)
		collision, _, err := geoms["left-tip"].Geometries()[0].CollidesWith(geoms["right-tip"].Geometries()[0], defaultCollisionBufferMM)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, got == nil, test.ShouldEqual, !collision)
		if collision {
			collisions++
		}
	}
	test.That(t, collisions, test.ShouldBeGreaterThan, 0)
	test.That(t, checker.CheckShortcutFS(context.Background(), segment, 2, true), test.ShouldNotBeNil)
}
