package motionplan

import (
	"math"
	"sync/atomic"
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/spatialmath"
)

func TestCollisionHintReusePreservesVerdictAndClearance(t *testing.T) {
	logger := logging.NewTestLogger(t)
	for _, mesh := range []bool{false, true} {
		for _, self := range []bool{false, true} {
			for _, offset := range []float64{0, .5, 1, 1.000000001, 2, 5, 20} {
				makeGeometry := func(name string, offset float64) spatialmath.Geometry {
					pose := spatialmath.NewPoseFromPoint(r3.Vector{Z: offset})
					if mesh {
						return spatialmath.NewMesh(pose, []*spatialmath.Triangle{
							spatialmath.NewTriangle(r3.Vector{X: -5, Y: -5}, r3.Vector{X: 5, Y: -5}, r3.Vector{Y: 5}),
						}, name)
					}
					g, err := spatialmath.NewBox(pose, r3.Vector{X: 10, Y: 1, Z: 1}, name)
					test.That(t, err, test.ShouldBeNil)
					return g
				}
				a, b := makeGeometry("a", 0), makeGeometry("b", offset)
				gg, other := []spatialmath.Geometry{a}, []spatialmath.Geometry{b}
				if self {
					gg = []spatialmath.Geometry{a, b}
					other = gg
				}
				expected, _, err := checkCollisionsHinted(gg, other, nil, nil, defaultCollisionBufferMM, false, nil, logger)
				test.That(t, err, test.ShouldBeNil)
				directCollision, directDistance, err := a.CollidesWith(b, defaultCollisionBufferMM)
				test.That(t, err, test.ShouldBeNil)
				for _, hintNames := range [][2]string{{"a", "b"}, {"b", "a"}, {"missing", "b"}} {
					var hint atomic.Pointer[[2]string]
					hint.Store(&hintNames)
					got, distance, err := checkCollisionsHinted(gg, other, nil, nil,
						defaultCollisionBufferMM, false, &hint, logger)
					test.That(t, err, test.ShouldBeNil)
					test.That(t, got, test.ShouldResemble, expected)
					if !directCollision {
						// Returned clearance is a lower bound, possibly from the
						// bounding sphere or a cached native mesh witness.
						test.That(t, math.IsNaN(distance), test.ShouldBeFalse)
						test.That(t, distance, test.ShouldBeLessThanOrEqualTo, directDistance+1e-8)
					}
				}
			}
		}
	}
}

type countedHintGeometry struct {
	spatialmath.Geometry
	calls *int
}

func (g countedHintGeometry) CollidesWith(spatialmath.Geometry, float64) (bool, float64, error) {
	*g.calls++
	return false, 1, nil
}

func TestCollisionHintReuseRetainsUnknownGeometryCalls(t *testing.T) {
	a, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "a")
	test.That(t, err, test.ShouldBeNil)
	b, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "b")
	test.That(t, err, test.ShouldBeNil)
	calls := 0
	gg := []spatialmath.Geometry{countedHintGeometry{a, &calls}, countedHintGeometry{b, &calls}}
	var hint atomic.Pointer[[2]string]
	hint.Store(&[2]string{"a", "b"})
	_, _, err = checkCollisionsHinted(gg, gg, nil, nil, defaultCollisionBufferMM, false,
		&hint, logging.NewTestLogger(t))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, calls, test.ShouldEqual, 3)
}

func TestCollisionHintReuseRetainsCollectAll(t *testing.T) {
	a, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "a")
	test.That(t, err, test.ShouldBeNil)
	b, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "b")
	test.That(t, err, test.ShouldBeNil)
	gg := []spatialmath.Geometry{a, b}
	var hint atomic.Pointer[[2]string]
	hint.Store(&[2]string{"a", "b"})
	got, _, err := checkCollisionsHinted(gg, gg, nil, nil, defaultCollisionBufferMM, true,
		&hint, logging.NewTestLogger(t))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, len(got), test.ShouldEqual, 3)
}
