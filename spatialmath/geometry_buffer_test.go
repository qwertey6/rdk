package spatialmath

import (
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"
)

func TestGeometryTransformBufferParity(t *testing.T) {
	resetMesh := func(m *Mesh) {
		m.ResetCache()
		if m.state != nil {
			m.state.distAnchors.Clear()
			m.state.geomAnchors.Clear()
		}
	}
	reset := func(g Geometry) {
		switch shape := g.(type) {
		case *Mesh:
			resetMesh(shape)
		case *Cylinder:
			resetMesh(shape.mesh)
		case *box:
			if shape.mesh != nil {
				resetMesh(shape.mesh)
			}
		}
	}
	pose := NewPose(r3.Vector{X: 20, Y: 30}, &EulerAngles{Pitch: .2})
	box, err := NewBox(pose, r3.Vector{X: 40, Y: 50, Z: 60}, "box")
	test.That(t, err, test.ShouldBeNil)
	sphere, err := NewSphere(pose, 20, "sphere")
	test.That(t, err, test.ShouldBeNil)
	capsule, err := NewCapsule(pose, 10, 80, "capsule")
	test.That(t, err, test.ShouldBeNil)
	cylinder, err := NewCylinderWithCapped(pose, 20, 50, false, "tube")
	test.That(t, err, test.ShouldBeNil)
	mesh := NewMesh(pose, []*Triangle{NewTriangle(r3.Vector{X: 20}, r3.Vector{Y: 20}, r3.Vector{Z: 20})}, "mesh")
	shapes := []Geometry{box, sphere, capsule, cylinder, mesh, NewPoint(r3.Vector{X: 10}, "point")}
	for _, shape := range shapes {
		t.Run(shape.Label(), func(t *testing.T) {
			b := NewGeometryTransformBuffer(shape)
			for k := 0; k < 60; k++ {
				parent := NewPose(r3.Vector{X: float64(k), Y: -float64(k)}, &EulerAngles{Yaw: .1 * float64(k)})
				want, got := shape.Transform(parent), b.Transform(parent)
				test.That(t, PoseAlmostEqualEps(got.Pose(), want.Pose(), 1e-12), test.ShouldBeTrue)
				test.That(t, got.Label(), test.ShouldEqual, want.Label())
				for _, other := range shapes {
					// Mesh clearance anchors return different conservative lower
					// bounds on a warm hit. Compare exact paths with equal history.
					reset(want)
					reset(other)
					wc, wd, we := want.CollidesWith(other, .5)
					reset(got)
					reset(other)
					gc, gd, ge := got.CollidesWith(other, .5)
					test.That(t, ge, test.ShouldResemble, we)
					test.That(t, gc, test.ShouldEqual, wc)
					if we == nil {
						test.That(t, gd, test.ShouldAlmostEqual, wd, 1e-8)
					}
				}
			}
		})
	}
}
