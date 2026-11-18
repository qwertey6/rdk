package spatialmath

import (
	"math"
	"math/rand"
	"testing"

	"github.com/golang/geo/r3"
	"gonum.org/v1/gonum/num/dualquat"
	"gonum.org/v1/gonum/num/quat"
)

// Frozen pre-change formula: compare exact output bits, including zero signs,
// nonunit inputs, over/underflow and nonfinite payloads.
func originalDualQuaternionPoint(q *DualQuaternion) r3.Vector {
	p := dualquat.Mul(q.Number, dualquat.Conj(q.Number)).Dual
	return r3.Vector{p.Imag, p.Jmag, p.Kmag}
}

func TestDualQuaternionPointOriginalBits(t *testing.T) {
	check := func(q DualQuaternion) {
		t.Helper()
		want, got := originalDualQuaternionPoint(&q), q.Point()
		for i, pair := range [][2]float64{{want.X, got.X}, {want.Y, got.Y}, {want.Z, got.Z}} {
			if math.Float64bits(pair[0]) != math.Float64bits(pair[1]) {
				t.Fatalf("coordinate %d: want %016x got %016x for %#v", i, math.Float64bits(pair[0]), math.Float64bits(pair[1]), q.Number)
			}
		}
	}
	special := []float64{
		0, math.Copysign(0, -1), 1, -1, math.Nextafter(1, 0), math.Nextafter(1, 2),
		math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, math.MaxFloat64, -math.MaxFloat64,
		1e-150, -1e-150, 1e150, -1e150, math.Inf(1), math.Inf(-1), math.NaN(),
		math.Float64frombits(0x7ff8000000000042), math.Float64frombits(0xfff8000000000123),
	}
	for _, a := range special {
		for _, b := range special {
			check(DualQuaternion{Number: dualquat.Number{
				Real: quat.Number{Real: a, Imag: b, Jmag: .25, Kmag: -.5},
				Dual: quat.Number{Real: b, Imag: -a, Jmag: 3, Kmag: -7},
			}})
		}
	}
	rng := rand.New(rand.NewSource(5271))
	for i := 0; i < 200000; i++ {
		var v [8]float64
		for j := range v {
			if i%2 == 0 {
				v[j] = math.Float64frombits(rng.Uint64())
			} else {
				v[j] = rng.NormFloat64() * 1000
			}
		}
		check(DualQuaternion{Number: dualquat.Number{
			Real: quat.Number{Real: v[0], Imag: v[1], Jmag: v[2], Kmag: v[3]},
			Dual: quat.Number{Real: v[4], Imag: v[5], Jmag: v[6], Kmag: v[7]},
		}})
	}
}

var pointBenchmarkSink r3.Vector

func BenchmarkDualQuaternionPoint(b *testing.B) {
	var poses [256]DualQuaternion
	rng := rand.New(rand.NewSource(5271))
	for i := range poses {
		poses[i].Real = quat.Number{Real: rng.NormFloat64(), Imag: rng.NormFloat64(), Jmag: rng.NormFloat64(), Kmag: rng.NormFloat64()}
		poses[i].Real = quat.Scale(1/quat.Abs(poses[i].Real), poses[i].Real)
		poses[i].SetTranslation(r3.Vector{X: rng.Float64() * 1000, Y: rng.Float64() * 1000, Z: rng.Float64() * 1000})
	}
	b.Run("original", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pointBenchmarkSink = originalDualQuaternionPoint(&poses[i%len(poses)])
		}
	})
	b.Run("candidate", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pointBenchmarkSink = poses[i%len(poses)].Point()
		}
	})
}
