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

func adaptiveArms(t *testing.T) (*ConstraintChecker, *SegmentFS) {
	t.Helper()
	fs := referenceframe.NewEmptyFrameSystem("arms")
	moving := map[string]bool{}
	for _, name := range []string{"left", "right"} {
		joint, err := referenceframe.NewRotationalFrame(name, spatialmath.R4AA{RZ: 1}, referenceframe.Limit{Min: -2, Max: 2})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, fs.AddFrame(joint, fs.World()), test.ShouldBeNil)
		g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000}), 1, name+"-sphere")
		test.That(t, err, test.ShouldBeNil)
		link, err := referenceframe.NewStaticFrameWithGeometry(name+"-link", spatialmath.NewZeroPose(), g)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, fs.AddFrame(link, joint), test.ShouldBeNil)
		moving[name], moving[link.Name()] = true, true
	}
	from := referenceframe.FrameSystemInputs{"left": {.7}, "right": {1.3}}.ToLinearInputs()
	to := referenceframe.FrameSystemInputs{"left": {1.3}, "right": {.7}}.ToLinearInputs()
	geoms, err := referenceframe.FrameSystemGeometriesLinearInputs(fs, from)
	test.That(t, err, test.ShouldBeNil)
	var robot []spatialmath.Geometry
	for _, gif := range geoms {
		robot = append(robot, gif.Geometries()...)
	}
	c, err := NewConstraintChecker(1e-8, nil, nil, nil, fs, robot, nil, moving, from, nil,
		logging.NewTestLogger(t), NewCollisionCache())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, c.motionBounds, test.ShouldNotBeNil)
	return c, &SegmentFS{FS: fs, StartConfiguration: from, EndConfiguration: to}
}

func TestAdaptiveRelativeMotionAndPrefix(t *testing.T) {
	c, edge := adaptiveArms(t)
	prefix, err := c.CheckStateConstraintsAcrossSegmentFS(context.Background(), edge, 1, true)
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, prefix, test.ShouldNotBeNil)
	c.motionBounds = nil
	dense, err := c.CheckStateConstraintsAcrossSegmentFS(context.Background(), edge, 1, true)
	test.That(t, err, test.ShouldNotBeNil)
	test.That(t, prefix.EndConfiguration.GetLinearizedInputs(), test.ShouldResemble, dense.EndConfiguration.GetLinearizedInputs())
}

func TestAdaptiveKeepsEveryTopologyCheck(t *testing.T) {
	c, edge := adaptiveArms(t)
	edge.EndConfiguration.Get("left")[0], edge.EndConfiguration.Get("right")[0] = .8, 1.2
	checks := 0
	c.topoConstraint = func(*StateFS) error {
		checks++
		if checks == 5 {
			return fmt.Errorf("topology violation")
		}
		return nil
	}
	prefix, err := c.CheckStateConstraintsAcrossSegmentFS(context.Background(), edge, 1, true)
	test.That(t, err.Error(), test.ShouldEqual, "topology violation")
	test.That(t, checks, test.ShouldEqual, 5)
	samples, err := InterpolateSegmentFS(edge, 1)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, prefix.EndConfiguration.GetLinearizedInputs(), test.ShouldResemble, samples[3].GetLinearizedInputs())
}

func TestAdaptiveCustomCallbacksFallBack(t *testing.T) {
	c, edge := adaptiveArms(t)
	checks := 0
	c.SetCollisionConstraints(CollisionConstraints{Obstacle: func(*StateFS) (float64, error) {
		checks++
		return 1e9, nil
	}})
	_, err := c.CheckStateConstraintsAcrossSegmentFS(context.Background(), edge, 1, false)
	test.That(t, err, test.ShouldBeNil)
	n, err := SegmentSampleCount(edge, 1)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, checks, test.ShouldEqual, n-1)
}

func TestAdaptiveCertifiesSeparatedInterval(t *testing.T) {
	c, edge := adaptiveArms(t)
	edge.EndConfiguration.Get("left")[0], edge.EndConfiguration.Get("right")[0] = .8, 1.2
	validation, err := c.SegmentValidation(edge, 1)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, validation.certified, test.ShouldBeTrue)
	_, _, err = validation.Check(context.Background(), 0, 0)
	test.That(t, err, test.ShouldBeNil)
	_, _, err = validation.Check(context.Background(), 1, 0)
	test.That(t, err, test.ShouldBeNil)
	// A retained anchor proves the next sample reused the clearance certificate.
	test.That(t, validation.anchor, test.ShouldEqual, 0)
}
