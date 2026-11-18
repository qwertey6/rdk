package referenceframe

import (
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"

	"go.viam.com/rdk/spatialmath"
)

func TestMotionBoundNativeIdentityDiscontinuity(t *testing.T) {
	fs := NewEmptyFrameSystem("numerical")
	r, err := NewRotationalFrame("joint", spatialmath.R4AA{RZ: 1}, Limit{Min: -2, Max: 2})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(r, fs.World()), test.ShouldBeNil)
	g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000}), 1, "tip")
	test.That(t, err, test.ShouldBeNil)
	link, err := NewStaticFrameWithGeometry("tip", spatialmath.NewZeroPose(), g)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(link, r), test.ShouldBeNil)
	from := FrameSystemInputs{"joint": {1.9e-8}}.ToLinearInputs()
	to := FrameSystemInputs{"joint": {2.2e-8}}.ToLinearInputs()
	p := fs.PrepareMotionBounds()
	speed, ok := p.maxIdealPointSpeed(from, to)
	test.That(t, ok, test.ShouldBeTrue)
	a, err := FrameSystemGeometriesLinearInputs(fs, from)
	test.That(t, err, test.ShouldBeNil)
	b, err := FrameSystemGeometriesLinearInputs(fs, to)
	test.That(t, err, test.ShouldBeNil)
	displacement := a["tip"].Geometries()[0].Pose().Point().Sub(b["tip"].Geometries()[0].Pose().Point()).Norm()
	// The original ideal-rigid bound alone is not valid for this arithmetic.
	test.That(t, displacement > speed, test.ShouldBeTrue)
	test.That(t, p.InterpolationStable(from, to), test.ShouldBeFalse)
	_, certified := p.MaxPointSpeed(from, to)
	test.That(t, certified, test.ShouldBeFalse)
	from.Get("joint")[0], to.Get("joint")[0] = .7, .8
	test.That(t, p.InterpolationStable(from, to), test.ShouldBeTrue)
}

func TestMotionBoundComposedIdentityFallback(t *testing.T) {
	fs := NewEmptyFrameSystem("cancellation")
	parent := fs.World()
	for _, name := range []string{"first", "second"} {
		r, err := NewRotationalFrame(name, spatialmath.R4AA{RZ: 1}, Limit{Min: -2, Max: 2})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, fs.AddFrame(r, parent), test.ShouldBeNil)
		parent = r
	}
	from := FrameSystemInputs{"first": {.7}, "second": {-.7}}.ToLinearInputs()
	to := FrameSystemInputs{"first": {.8}, "second": {-.8}}.ToLinearInputs()
	p := fs.PrepareMotionBounds()
	test.That(t, p.InterpolationStable(from, to), test.ShouldBeFalse)
	from.Get("second")[0], to.Get("second")[0] = .7, .8
	test.That(t, p.InterpolationStable(from, to), test.ShouldBeTrue)
}

func TestMotionBoundRobotWorldRelativeIdentity(t *testing.T) {
	fs := NewEmptyFrameSystem("relative")
	r, err := NewRotationalFrame("joint", spatialmath.R4AA{RZ: 1}, Limit{Min: -2, Max: 2})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(r, fs.World()), test.ShouldBeNil)
	g, err := spatialmath.NewBox(spatialmath.NewPoseFromPoint(r3.Vector{X: 100}), r3.Vector{X: 1, Y: 1, Z: 1}, "tip")
	test.That(t, err, test.ShouldBeNil)
	tip, err := NewStaticFrameWithGeometry("tip", spatialmath.NewZeroPose(), g)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(tip, r), test.ShouldBeNil)
	obstacle, err := spatialmath.NewBox(
		spatialmath.NewPose(r3.Vector{X: 500}, &spatialmath.EulerAngles{Yaw: .75}), r3.Vector{X: 1, Y: 1, Z: 1}, "world",
	)
	test.That(t, err, test.ShouldBeNil)
	from := FrameSystemInputs{"joint": {.7}}.ToLinearInputs()
	to := FrameSystemInputs{"joint": {.8}}.ToLinearInputs()
	test.That(t, fs.PrepareMotionBounds().InterpolationStable(from, to), test.ShouldBeTrue)
	test.That(t, fs.PrepareMotionBounds(obstacle).InterpolationStable(from, to), test.ShouldBeFalse)
}
