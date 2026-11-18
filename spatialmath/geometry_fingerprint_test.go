package spatialmath

import (
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"
)

func TestGeometrySetHashCompleteShape(t *testing.T) {
	a, err := NewBox(NewZeroPose(), r3.Vector{X: 10, Y: 10, Z: 10}, "box")
	test.That(t, err, test.ShouldBeNil)
	b, err := NewBox(NewZeroPose(), r3.Vector{X: 10.00001, Y: 10, Z: 10}, "box")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, GeometrySetHash([]Geometry{a}), test.ShouldNotEqual, GeometrySetHash([]Geometry{b}))
	triangles := make([]*Triangle, 12)
	for i := range triangles {
		triangles[i] = NewTriangle(r3.Vector{X: float64(i)}, r3.Vector{Y: 1}, r3.Vector{Z: 1})
	}
	m1 := NewMesh(NewZeroPose(), append([]*Triangle(nil), triangles...), "mesh")
	triangles[11] = NewTriangle(r3.Vector{X: 200}, r3.Vector{Y: 1}, r3.Vector{Z: 1})
	m2 := NewMesh(NewZeroPose(), triangles, "mesh")
	test.That(t, GeometrySetHash([]Geometry{m1}), test.ShouldNotEqual, GeometrySetHash([]Geometry{m2}))
	test.That(t, GeometrySetHash([]Geometry{a, m1}), test.ShouldEqual, GeometrySetHash([]Geometry{m1, a}))
}
