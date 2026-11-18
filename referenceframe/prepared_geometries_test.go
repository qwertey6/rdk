package referenceframe

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"

	"go.viam.com/rdk/spatialmath"
	"go.viam.com/rdk/utils"
)

func TestPreparedGeometriesParity(t *testing.T) {
	for _, name := range []string{"xarm6", "ur5e"} {
		t.Run(name, func(t *testing.T) {
			model, err := ParseModelJSONFile(utils.ResolveFile("components/arm/kinematics/"+name+".json"), "arm")
			test.That(t, err, test.ShouldBeNil)
			fs := NewEmptyFrameSystem("geometry")
			mount, err := NewStaticFrame("mount", spatialmath.NewPose(r3.Vector{X: 100, Y: -20}, &spatialmath.EulerAngles{Yaw: .3}))
			test.That(t, err, test.ShouldBeNil)
			test.That(t, fs.AddFrame(mount, fs.World()), test.ShouldBeNil)
			test.That(t, fs.AddFrame(model, mount), test.ShouldBeNil)
			prepared := fs.PrepareGeometries(nil)
			rng := rand.New(rand.NewSource(42))
			for trial := 0; trial < 100; trial++ {
				q := FrameSystemInputs{"arm": RandomFrameInputs(model, rng)}.ToLinearInputs()
				original, err := FrameSystemGeometriesLinearInputs(fs, q)
				test.That(t, err, test.ShouldBeNil)
				want := map[string]spatialmath.Geometry{}
				for _, gif := range original {
					for _, g := range gif.Geometries() {
						want[g.Label()] = g
					}
				}
				e := prepared.Evaluate(q)
				got, err := e.Subset(nil)
				test.That(t, err, test.ShouldBeNil)
				test.That(t, len(got), test.ShouldEqual, len(want))
				for _, g := range got {
					w := want[g.Label()]
					test.That(t, w, test.ShouldNotBeNil)
					test.That(t, spatialmath.PoseAlmostEqualEps(g.Pose(), w.Pose(), 1e-10), test.ShouldBeTrue)
					// Exercise pose-dependent collision caches before recycling.
					probe, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(g.Pose().Point().Add(r3.Vector{X: float64(trial)})), 5, "probe")
					test.That(t, err, test.ShouldBeNil)
					gc, gd, ge := g.CollidesWith(probe, .1)
					wc, wd, we := w.CollidesWith(probe, .1)
					test.That(t, ge, test.ShouldResemble, we)
					test.That(t, gc, test.ShouldEqual, wc)
					test.That(t, gd, test.ShouldAlmostEqual, wd, 1e-9)
				}
				e.Release()
			}
		})
	}
}

func TestPreparedGeometriesConcurrent(t *testing.T) {
	model, err := ParseModelJSONFile(utils.ResolveFile("components/arm/kinematics/xarm6.json"), "arm")
	test.That(t, err, test.ShouldBeNil)
	fs := NewEmptyFrameSystem("concurrent")
	test.That(t, fs.AddFrame(model, fs.World()), test.ShouldBeNil)
	prepared := fs.PrepareGeometries(nil)
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		q := FrameSystemInputs{"arm": RandomFrameInputs(model, rand.New(rand.NewSource(int64(worker))))}.ToLinearInputs()
		original, err := FrameSystemGeometriesLinearInputs(fs, q)
		test.That(t, err, test.ShouldBeNil)
		want := map[string]spatialmath.Pose{}
		for _, gif := range original {
			for _, g := range gif.Geometries() {
				want[g.Label()] = g.Pose()
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				e := prepared.Evaluate(q)
				got, err := e.Subset(nil)
				if err != nil {
					t.Error(err)
					e.Release()
					return
				}
				for _, g := range got {
					if !spatialmath.PoseAlmostEqualEps(g.Pose(), want[g.Label()], 1e-10) {
						t.Errorf("concurrent pose changed for %s", g.Label())
					}
				}
				e.Release()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkPreparedGeometries(b *testing.B) {
	model, err := ParseModelJSONFile(utils.ResolveFile("components/arm/kinematics/xarm6.json"), "arm")
	if err != nil {
		b.Fatal(err)
	}
	fs := NewEmptyFrameSystem("benchmark")
	if err := fs.AddFrame(model, fs.World()); err != nil {
		b.Fatal(err)
	}
	q := FrameSystemInputs{"arm": make([]Input, len(model.DoF()))}.ToLinearInputs()
	p := fs.PrepareGeometries(nil)
	b.Run("original", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := FrameSystemGeometriesLinearInputs(fs, q); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("reused", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			e := p.Evaluate(q)
			_, err := e.Subset(nil)
			e.Release()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
