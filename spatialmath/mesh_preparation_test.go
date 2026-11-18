package spatialmath

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"
)

func TestMeshUnmarshalInvalidatesPreparation(t *testing.T) {
	m := NewMesh(NewZeroPose(), []*Triangle{NewTriangle(r3.Vector{}, r3.Vector{X: 1}, r3.Vector{Y: 1})}, "mesh")
	before := GeometrySetHash([]Geometry{m})
	oldState := m.state
	oldBVH := m.ensureBVH()
	m.ensureUniqueVertices()
	m.ensurePoseCache()
	replacement := NewMesh(NewPoseFromPoint(r3.Vector{Z: 20}),
		[]*Triangle{NewTriangle(r3.Vector{}, r3.Vector{X: 3}, r3.Vector{Y: 4})}, "mesh")
	data, err := json.Marshal(replacement)
	test.That(t, err, test.ShouldBeNil)
	// Compare with a fresh decode: mesh serialization may quantize vertices.
	var decoded Mesh
	test.That(t, json.Unmarshal(data, &decoded), test.ShouldBeNil)
	test.That(t, json.Unmarshal(data, m), test.ShouldBeNil)
	test.That(t, m.state == oldState, test.ShouldBeFalse)
	test.That(t, m.ensureBVH() == oldBVH, test.ShouldBeFalse)
	test.That(t, GeometrySetHash([]Geometry{m}), test.ShouldNotEqual, before)
	test.That(t, PoseAlmostEqual(m.pose, replacement.pose), test.ShouldBeTrue)
	test.That(t, m.ensurePoseCache(), test.ShouldResemble, decoded.ensurePoseCache())
	points := map[r3.Vector]bool{}
	for _, p := range m.ensureUniqueVertices() {
		points[p] = true
	}
	test.That(t, len(points), test.ShouldEqual, len(decoded.ensureUniqueVertices()))
	for _, p := range decoded.ensureUniqueVertices() {
		test.That(t, points[p], test.ShouldBeTrue)
	}
	test.That(t, GeometrySetHash([]Geometry{m}), test.ShouldEqual, GeometrySetHash([]Geometry{&decoded}))
}

func TestMeshCacheResetConcurrent(t *testing.T) {
	mesh := NewMesh(NewZeroPose(), []*Triangle{NewTriangle(r3.Vector{X: 20}, r3.Vector{Y: 20}, r3.Vector{Z: 20})}, "mesh")
	sphere, err := NewSphere(NewPoseFromPoint(r3.Vector{X: 30}), 2, "sphere")
	test.That(t, err, test.ShouldBeNil)
	var wg sync.WaitGroup
	for worker := 0; worker < 3; worker++ {
		wg.Add(1)
		go func(reset bool) {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				if reset {
					mesh.ResetCache()
				} else {
					_, _, _ = mesh.CollidesWith(sphere, .1)
				}
			}
		}(worker == 0)
	}
	wg.Wait()
}
