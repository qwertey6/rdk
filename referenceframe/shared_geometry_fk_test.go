package referenceframe

import (
	"math"
	"math/rand"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/spatialmath"
	"go.viam.com/rdk/utils"
)

func TestGeometryEvaluationSharedFK(t *testing.T) {
	model, err := ParseModelJSONFile(utils.ResolveFile("components/arm/kinematics/xarm6.json"), "arm")
	test.That(t, err, test.ShouldBeNil)
	fs := NewEmptyFrameSystem("shared")
	test.That(t, fs.AddFrame(model, fs.World()), test.ShouldBeNil)
	shape, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "tool")
	test.That(t, err, test.ShouldBeNil)
	tool, err := NewStaticFrameWithGeometry("tool", spatialmath.NewZeroPose(), shape)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(tool, model), test.ShouldBeNil)
	p := fs.PrepareGeometries(nil)
	rng := rand.New(rand.NewSource(9021))
	for trial := 0; trial < 200; trial++ {
		q := FrameSystemInputs{"arm": RandomFrameInputs(model, rng)}.ToLinearInputs()
		e := p.Evaluate(q)
		separate := fs.NewFrameToWorld(q)
		for _, name := range fs.FrameNames() {
			wq, wp, err := separate.PoseQT(name)
			test.That(t, err, test.ShouldBeNil)
			gq, gp, err := e.PoseQT(name)
			test.That(t, err, test.ShouldBeNil)
			for _, pair := range [][2]float64{
				{wq.Real, gq.Real},
				{wq.Imag, gq.Imag},
				{wq.Jmag, gq.Jmag},
				{wq.Kmag, gq.Kmag},
				{wp.X, gp.X},
				{wp.Y, gp.Y},
				{wp.Z, gp.Z},
			} {
				test.That(t, math.Float64bits(pair[0]), test.ShouldEqual, math.Float64bits(pair[1]))
			}
		}
		separate.Release()
		untouched := p.Evaluate(q)
		want, err := untouched.Subset(nil)
		test.That(t, err, test.ShouldBeNil)
		got, err := e.Subset(nil)
		test.That(t, err, test.ShouldBeNil)
		test.That(t, len(got), test.ShouldEqual, len(want))
		for i, g := range got {
			test.That(t, g.Label(), test.ShouldEqual, want[i].Label())
			test.That(t, g.Pose().(*spatialmath.DualQuaternion).Number, test.ShouldResemble,
				want[i].Pose().(*spatialmath.DualQuaternion).Number)
		}
		_, _, err = e.PoseQT("missing")
		test.That(t, err, test.ShouldNotBeNil)
		e.Release()
		untouched.Release()
	}
}

func BenchmarkGeometryCoverFKSharing(b *testing.B) {
	model, err := ParseModelJSONFile(utils.ResolveFile("components/arm/kinematics/xarm6.json"), "arm")
	if err != nil {
		b.Fatal(err)
	}
	fs := NewEmptyFrameSystem("shared")
	if err := fs.AddFrame(model, fs.World()); err != nil {
		b.Fatal(err)
	}
	shape, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "tool")
	if err != nil {
		b.Fatal(err)
	}
	tool, err := NewStaticFrameWithGeometry("tool", spatialmath.NewZeroPose(), shape)
	if err != nil {
		b.Fatal(err)
	}
	if err := fs.AddFrame(tool, model); err != nil {
		b.Fatal(err)
	}
	p := fs.PrepareGeometries(nil)
	q := FrameSystemInputs{"arm": make([]Input, len(model.DoF()))}.ToLinearInputs()
	names := fs.FrameNames()
	for _, shared := range []bool{false, true} {
		name := "separate"
		if shared {
			name = "shared"
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				e := p.Evaluate(q)
				if shared {
					for _, name := range names {
						if _, _, err := e.PoseQT(name); err != nil {
							b.Fatal(err)
						}
					}
				} else {
					fk := fs.NewFrameToWorld(q)
					for _, name := range names {
						if _, _, err := fk.PoseQT(name); err != nil {
							b.Fatal(err)
						}
					}
					fk.Release()
				}
				if _, err := e.Subset(nil); err != nil {
					b.Fatal(err)
				}
				e.Release()
			}
		})
	}
}
