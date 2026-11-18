package spatialmath

import (
	"math"
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/num/dualquat"
	"gonum.org/v1/gonum/num/quat"
)

// Frozen arithmetic from HEAD, independent of the normalization fast path.
//
//nolint:dupl // Keep an independent copy of the original arithmetic as a compatibility oracle.
func originalTransformationForNormalizationTest(q *DualQuaternion, by dualquat.Number) dualquat.Number {
	var newReal quat.Number

	if q.Real.Real == 1 {
		// Since we're working with unit quaternions, if either Real is 1, then that quat is an identity quat
		newReal = by.Real
	} else if by.Real.Real == 1 {
		newReal = q.Real
	} else {
		newReal = quat.Mul(q.Real, by.Real)
	}
	// Multiplication is faster than division. Thus if this is hit, it is faster to divide once and multiply four times.
	// However, if this is not hit, it may be a wash. The branch predictor won't save us as the following lines rely on newReal.
	if vecLen := 1 / quat.Abs(newReal); vecLen-1 > 1e-10 || vecLen-1 < -1e-10 {
		newReal.Real *= vecLen
		newReal.Imag *= vecLen
		newReal.Jmag *= vecLen
		newReal.Kmag *= vecLen
	}

	if q.Dual.Real == 0 && q.Dual.Imag == 0 && q.Dual.Jmag == 0 && q.Dual.Kmag == 0 {
		return dualquat.Number{
			Real: newReal,
			Dual: quat.Mul(q.Real, by.Dual),
		}
	} else if by.Dual.Real == 0 && by.Dual.Imag == 0 && by.Dual.Jmag == 0 && by.Dual.Kmag == 0 {
		return dualquat.Number{
			Real: newReal,
			Dual: quat.Mul(q.Dual, by.Real),
		}
	}

	return dualquat.Number{
		Real: newReal,
		// Equivalent to but faster than quat.Add(quat.Mul(q.Real, by.Dual), quat.Mul(q.Dual, by.Real))
		Dual: quat.Number{
			Real: q.Real.Real*by.Dual.Real - q.Real.Imag*by.Dual.Imag - q.Real.Jmag*by.Dual.Jmag - q.Real.Kmag*by.Dual.Kmag +
				q.Dual.Real*by.Real.Real - q.Dual.Imag*by.Real.Imag - q.Dual.Jmag*by.Real.Jmag - q.Dual.Kmag*by.Real.Kmag,
			Imag: q.Real.Real*by.Dual.Imag + q.Real.Imag*by.Dual.Real + q.Real.Jmag*by.Dual.Kmag - q.Real.Kmag*by.Dual.Jmag +
				q.Dual.Real*by.Real.Imag + q.Dual.Imag*by.Real.Real + q.Dual.Jmag*by.Real.Kmag - q.Dual.Kmag*by.Real.Jmag,
			Jmag: q.Real.Real*by.Dual.Jmag - q.Real.Imag*by.Dual.Kmag + q.Real.Jmag*by.Dual.Real + q.Real.Kmag*by.Dual.Imag +
				q.Dual.Real*by.Real.Jmag - q.Dual.Imag*by.Real.Kmag + q.Dual.Jmag*by.Real.Real + q.Dual.Kmag*by.Real.Imag,
			Kmag: q.Real.Real*by.Dual.Kmag + q.Real.Imag*by.Dual.Jmag - q.Real.Jmag*by.Dual.Imag + q.Real.Kmag*by.Dual.Real +
				q.Dual.Real*by.Real.Kmag + q.Dual.Imag*by.Real.Jmag - q.Dual.Jmag*by.Real.Imag + q.Dual.Kmag*by.Real.Real,
		},
	}
}

func TestTransformationNormalizationParity(t *testing.T) {
	rng := rand.New(rand.NewSource(817))
	compare := func(a, b dualquat.Number) {
		t.Helper()
		q := DualQuaternion{Number: a}
		got, want := q.Transformation(b), originalTransformationForNormalizationTest(&q, b)
		parts := func(v dualquat.Number) [8]float64 {
			return [8]float64{v.Real.Real, v.Real.Imag, v.Real.Jmag, v.Real.Kmag, v.Dual.Real, v.Dual.Imag, v.Dual.Jmag, v.Dual.Kmag}
		}
		g, w := parts(got), parts(want)
		for i := range g {
			if math.Float64bits(g[i]) != math.Float64bits(w[i]) {
				t.Fatalf("component %d differs: %x vs %x for %v * %v", i, math.Float64bits(g[i]), math.Float64bits(w[i]), a, b)
			}
		}
	}
	random := func() quat.Number {
		return quat.Number{Real: rng.NormFloat64(), Imag: rng.NormFloat64(), Jmag: rng.NormFloat64(), Kmag: rng.NormFloat64()}
	}
	scales := []float64{
		0, math.Copysign(0, -1), 1, -1,
		1 - 1e-13, 1 + 1e-13, 1 - 1e-12, 1 + 1e-12,
		1 - 1e-10, 1 + 1e-10, 1e-300, 1e300, math.Inf(1), math.NaN(),
	}
	for n := range 100000 {
		a, b := random(), random()
		a = quat.Scale(scales[n%len(scales)]/quat.Abs(a), a)
		b = quat.Scale(scales[(n/len(scales))%len(scales)]/quat.Abs(b), b)
		compare(dualquat.Number{Real: a, Dual: random()}, dualquat.Number{Real: b, Dual: random()})
	}
	// Dense probes around both the conservative gate and the original normalization
	// threshold. Identity shortcuts with nonzero vector parts retain native behavior.
	for _, center := range []float64{1 - 1e-10, 1 + 1e-10, 1 - 5e-13, 1 + 5e-13, 1} {
		for i := -1000; i <= 1000; i++ {
			scale := center + float64(i)*(math.Nextafter(1, 2)-1)
			for _, a := range []quat.Number{
				{Real: scale},
				{Real: scale * .5, Imag: scale * .5, Jmag: scale * .5, Kmag: scale * .5},
				{Real: 1, Imag: scale * 1e-9},
			} {
				compare(dualquat.Number{Real: a}, dualquat.Number{Real: quat.Number{Real: 1}})
			}
		}
	}
}

var normalizationBenchmarkResult dualquat.Number

func BenchmarkTransformationNormalization(b *testing.B) {
	q := DualQuaternion{Number: dualquat.Number{
		Real: quat.Number{Real: .5, Imag: .5, Jmag: .5, Kmag: .5},
		Dual: quat.Number{Imag: 23, Jmag: -44, Kmag: 11},
	}}
	by := q.Number
	for _, original := range []bool{true, false} {
		name := "candidate"
		if original {
			name = "original"
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				if original {
					normalizationBenchmarkResult = originalTransformationForNormalizationTest(&q, by)
				} else {
					normalizationBenchmarkResult = q.Transformation(by)
				}
			}
		})
	}
}
