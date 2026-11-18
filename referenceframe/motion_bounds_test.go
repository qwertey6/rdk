package referenceframe

import (
	"math"
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"

	"go.viam.com/rdk/spatialmath"
)

func TestMotionBoundsDistalAndRelative(t *testing.T) {
	fs := NewEmptyFrameSystem("bounds")
	for _, name := range []string{"left", "right"} {
		joint, err := NewRotationalFrame(name, spatialmath.R4AA{RZ: 1}, Limit{Min: -20, Max: 20})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, fs.AddFrame(joint, fs.World()), test.ShouldBeNil)
		g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000}), 10, name+"-sphere")
		test.That(t, err, test.ShouldBeNil)
		link, err := NewStaticFrameWithGeometry(name+"-link", spatialmath.NewZeroPose(), g)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, fs.AddFrame(link, joint), test.ShouldBeNil)
	}
	from := FrameSystemInputs{"left": {0}, "right": {1}}.ToLinearInputs()
	to := FrameSystemInputs{"left": {1}, "right": {0}}.ToLinearInputs()
	bound := fs.PrepareMotionBounds()
	test.That(t, bound, test.ShouldNotBeNil)
	speed, ok := bound.maxIdealPointSpeed(from, to)
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, speed, test.ShouldAlmostEqual, 1010, 1e-5)
	interp, ok := PrepareInterpolation(fs, from, to)
	test.That(t, ok, test.ShouldBeTrue)
	var prev map[string]r3.Vector
	for i := 0; i <= 100; i++ {
		geoms, err := FrameSystemGeometriesLinearInputs(fs, interp.At(float64(i)/100))
		test.That(t, err, test.ShouldBeNil)
		points := map[string]r3.Vector{}
		for name, gif := range geoms {
			if len(gif.Geometries()) == 1 {
				points[name] = gif.Geometries()[0].Pose().Point()
			}
		}
		if prev != nil {
			for name, point := range points {
				test.That(t, point.Sub(prev[name]).Norm() <= speed/100, test.ShouldBeTrue)
			}
			closing := points["left-link"].Sub(points["right-link"]).Sub(prev["left-link"].Sub(prev["right-link"]))
			test.That(t, closing.Norm() <= 2*speed/100, test.ShouldBeTrue)
		}
		prev = points
	}
	to.Get("left")[0] = math.NaN()
	_, ok = bound.maxIdealPointSpeed(from, to)
	test.That(t, ok, test.ShouldBeFalse)
}

func TestMotionBoundsPrismaticReach(t *testing.T) {
	fs := NewEmptyFrameSystem("extension")
	r, err := NewRotationalFrame("rotation", spatialmath.R4AA{RZ: 1}, Limit{Min: -10, Max: 10})
	test.That(t, err, test.ShouldBeNil)
	g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 50}), 10, "tip")
	test.That(t, err, test.ShouldBeNil)
	p, err := NewTranslationalFrameWithGeometry("extension", r3.Vector{X: 1}, Limit{Min: -10000, Max: 10000}, g)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(r, fs.World()), test.ShouldBeNil)
	test.That(t, fs.AddFrame(p, r), test.ShouldBeNil)
	from := FrameSystemInputs{"rotation": {0}, "extension": {-1000}}.ToLinearInputs()
	to := FrameSystemInputs{"rotation": {8}, "extension": {2000}}.ToLinearInputs()
	speed, ok := fs.PrepareMotionBounds().maxIdealPointSpeed(from, to)
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, speed, test.ShouldAlmostEqual, 3000+8*(2000+60), 1e-4)
}

func TestMotionBoundsModelBranchMimicAndTool(t *testing.T) {
	inner := NewEmptyFrameSystem("branches")
	for _, name := range []string{"source", "mimic"} {
		r, err := NewRotationalFrame(name, spatialmath.R4AA{RZ: 1}, Limit{Min: -100, Max: 100})
		test.That(t, err, test.ShouldBeNil)
		test.That(t, inner.AddFrame(r, inner.World()), test.ShouldBeNil)
	}
	g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 1000}), 1, "branch")
	test.That(t, err, test.ShouldBeNil)
	link, err := NewStaticFrameWithGeometry("branch", spatialmath.NewZeroPose(), g)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, inner.AddFrame(link, inner.Frame("mimic")), test.ShouldBeNil)
	model, err := NewModelWithMimics("arm", inner, "source", map[string]*mimicMapping{
		"mimic": {sourceFrameName: "source", valueMultiplier: -5, valueOffset: .5},
	})
	test.That(t, err, test.ShouldBeNil)
	fs := NewEmptyFrameSystem("outer")
	test.That(t, fs.AddFrame(model, fs.World()), test.ShouldBeNil)
	g, err = spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: 2000}), 1, "tool")
	test.That(t, err, test.ShouldBeNil)
	tool, err := NewStaticFrameWithGeometry("tool", spatialmath.NewZeroPose(), g)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(tool, model), test.ShouldBeNil)
	from := FrameSystemInputs{"arm": {0}}.ToLinearInputs()
	to := FrameSystemInputs{"arm": {1}}.ToLinearInputs()
	speed, ok := fs.PrepareMotionBounds().maxIdealPointSpeed(from, to)
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, speed, test.ShouldAlmostEqual, 5005, 1e-4)
	// The derived branch still participates in sampling even though the model's
	// primary output (and the attached tool) move five times more slowly.
	steps, ok := JointInterpolationSteps(model, []Input{0}, []Input{1}, 1, .002)
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, steps, test.ShouldEqual, 2500)
}

type unknownMotionFrame struct{ Frame }

func TestMotionBoundsCustomFrameFallback(t *testing.T) {
	fs := NewEmptyFrameSystem("fallback")
	test.That(t, fs.AddFrame(unknownMotionFrame{NewZeroStaticFrame("custom")}, fs.World()), test.ShouldBeNil)
	test.That(t, fs.PrepareMotionBounds(), test.ShouldBeNil)
}
