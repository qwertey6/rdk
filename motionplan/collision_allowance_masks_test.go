package motionplan

import (
	"fmt"
	"math"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/golang/geo/r3"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/spatialmath"
)

func maskTestAllowances(n int) *collisionAllowances {
	var pairs []Collision
	for i := range n {
		a := fmt.Sprintf("g%02d", i)
		// Register every name, including the single-name and 64th-bit cases.
		pairs = append(pairs, Collision{a, a})
		for j := range i {
			if (i+j)%3 == 0 {
				b := fmt.Sprintf("g%02d", j)
				pairs = append(pairs, Collision{a, b}, Collision{b, a}, Collision{a, b})
			}
		}
	}
	return makeAllowedCollisionsLookup(pairs)
}

func TestCollisionAllowanceMasksMembership(t *testing.T) {
	for _, n := range []int{0, 1, 2, 63, 64, 65, 128} {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			a := maskTestAllowances(n)
			if (a.compact != nil) != (n > 0 && n <= 64) {
				t.Fatalf("unexpected compact index for %d names", n)
			}
			if a.compact == nil {
				return
			}
			names := []string{"", "unknown", "λ", unnamedCollisionGeometryPrefix + "0"}
			for i := range n {
				names = append(names, fmt.Sprintf("g%02d", i))
			}
			for _, x := range names {
				for _, y := range names {
					got := a.compact[x].partners&a.compact[y].self != 0
					if want := a.pairs[canonicalPair(x, y)]; got != want {
						t.Fatalf("%q/%q: got %v, map says %v", x, y, got, want)
					}
				}
			}
		})
	}
}

// Compare whole collision calls against the existing map implementation. This
// covers first collision, full enumeration, pair hints, synthetic names, pooled
// scratch reuse and static pre-naming, including exact clearance/error results.
func TestCollisionAllowanceMasksDifferential(t *testing.T) {
	logger := logging.NewTestLogger(t)
	var geoms []spatialmath.Geometry
	for i, label := range []string{"g00", "g01", "g02", "g03", "g63", "g00", "", "", "unknown"} {
		g, err := spatialmath.NewSphere(spatialmath.NewPoseFromPoint(r3.Vector{X: float64(i % 4)}), 1, label)
		if err != nil {
			t.Fatal(err)
		}
		geoms = append(geoms, g)
	}
	for _, n := range []int{0, 1, 8, 64, 65} {
		compact := maskTestAllowances(n)
		plain := &collisionAllowances{pairs: compact.pairs, names: compact.names}
		for _, collectAll := range []bool{false, true} {
			for _, mode := range []string{"self", "separate", "static"} {
				for _, hintNames := range [][2]string{{"g02", "g01"}, {"g00", "g03"}, {"missing", "g63"}} {
					var compactHint, plainHint atomic.Pointer[[2]string]
					compactHint.Store(&hintNames)
					plainHint.Store(&hintNames)
					other := geoms
					if mode != "self" {
						other = geoms[2:]
					}
					var compactPre, plainPre []namedGeom
					if mode == "static" {
						compactPre, plainPre = prenameGeoms(other, compact), prenameGeoms(other, plain)
					}
					got, gd, ge := checkCollisionsHinted(geoms, other, compactPre, compact, .01, collectAll, &compactHint, logger)
					want, wd, we := checkCollisionsHinted(geoms, other, plainPre, plain, .01, collectAll, &plainHint, logger)
					if !reflect.DeepEqual(got, want) || math.Float64bits(gd) != math.Float64bits(wd) ||
						fmt.Sprint(ge) != fmt.Sprint(we) || !reflect.DeepEqual(compactHint.Load(), plainHint.Load()) {
						t.Fatalf("n=%d all=%v mode=%s hint=%v: got %v/%v/%v, want %v/%v/%v",
							n, collectAll, mode, hintNames, got, gd, ge, want, wd, we)
					}
				}
			}
		}
	}
	// Explicit synthetic labels can appear in an allowance. Preserve their
	// original interpretation, even when paired with an unnamed geometry.
	a := makeAllowedCollisionsLookup([]Collision{{unnamedCollisionGeometryPrefix + "0", "g01"}})
	unnamed, named := []spatialmath.Geometry{geoms[6]}, []spatialmath.Geometry{geoms[1]}
	got, _, err := checkCollisionsHinted(unnamed, named, nil, a, .01, true, nil, logger)
	if err != nil || len(got) != 0 {
		t.Fatalf("synthetic allowance changed: %v, %v", got, err)
	}
}

type allowanceErrorGeometry struct{ spatialmath.Geometry }

func (g allowanceErrorGeometry) CollidesWith(spatialmath.Geometry, float64) (bool, float64, error) {
	return false, 0, fmt.Errorf("deliberate unsupported pair")
}

func TestCollisionAllowanceMasksEarlyReturns(t *testing.T) {
	logger := logging.NewTestLogger(t)
	var geoms []spatialmath.Geometry
	for _, label := range []string{"g00", "g01"} {
		g, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, label)
		if err != nil {
			t.Fatal(err)
		}
		geoms = append(geoms, g)
	}
	for _, n := range []int{2, 65} {
		for _, fail := range []bool{false, true} {
			x, y := geoms[0], geoms[1]
			if fail {
				x, y = allowanceErrorGeometry{x}, allowanceErrorGeometry{y}
			}
			got, _, err := checkCollisionsHinted([]spatialmath.Geometry{x}, []spatialmath.Geometry{y}, nil,
				maskTestAllowances(n), .01, false, nil, logger)
			if fail && (err == nil || err.Error() != "deliberate unsupported pair") || !fail && (err != nil || len(got) != 1) {
				t.Fatalf("unexpected early return: %v, %v", got, err)
			}
		}
	}
}
