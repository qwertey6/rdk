package referenceframe

import (
	"math"

	"gonum.org/v1/gonum/num/quat"

	"go.viam.com/rdk/spatialmath"
)

func boundedQuaternion(q quat.Number) bool {
	n := quat.Abs(q)
	return finiteBound(n) && math.Abs(n-1) < 1e-12
}

// InterpolationStable rejects edges that could cross the native quaternion
// Real==1 shortcut with a nonzero vector part. Check every contiguous ancestor
// chain, including geometry pose composition: identity can arise by cancellation
// even when individual joint rotations are far from zero. Angular distance is
// Lipschitz with the sum of the chain's absolute joint travel. Unknown numerical
// regimes keep dense checking; this does not reject the planning request.
func (p *MotionBounds) InterpolationStable(from, to *LinearInputs) bool {
	if p == nil {
		return false
	}
	type rotation struct {
		q      quat.Number
		travel float64
	}
	rotations := make([]rotation, len(p.nodes))
	reaches := make([]float64, len(p.nodes))
	world := make([]rotation, len(p.nodes))
	world[0].q = quat.Number{Real: 1}
	for i := 1; i < len(p.nodes); i++ {
		n := p.nodes[i]
		rotations[i].q = n.orientation
		reach := n.reach
		if n.kind == 'r' || n.kind == 't' {
			a, b := from.Get(n.input.frame), to.Get(n.input.frame)
			if n.input.index < 0 || n.input.index >= len(a) || n.input.index >= len(b) {
				return false
			}
			x, y := a[n.input.index]*n.input.scale+n.input.bias, b[n.input.index]*n.input.scale+n.input.bias
			if !finiteBound(x) || !finiteBound(y) {
				return false
			}
			if n.kind == 'r' {
				if max(math.Abs(x), math.Abs(y)) > 100 {
					return false
				}
				aa := spatialmath.R4AA{Theta: x, RX: n.axis.X, RY: n.axis.Y, RZ: n.axis.Z}
				rotations[i] = rotation{aa.Quaternion(), math.Abs(y - x)}
			} else {
				reach = max(math.Abs(x), math.Abs(y)) * n.reach
			}
		}
		world[i] = rotation{quat.Mul(world[n.parent].q, rotations[i].q), world[n.parent].travel + rotations[i].travel}
		reaches[i] = reaches[n.parent] + reach
		if !finiteBound(reaches[i]) || reaches[i] > 1e6 {
			return false
		}
		q, travel, depth := quat.Number{Real: 1}, 0., 0
		for j := i; j > 0; j = p.nodes[j].parent {
			depth++
			if depth > 64 {
				return false
			}
			q = quat.Mul(rotations[j].q, q)
			travel += rotations[j].travel
			vector := math.Sqrt(q.Imag*q.Imag + q.Jmag*q.Jmag + q.Kmag*q.Kmag)
			if travel == 0 {
				// Exact identities and roundoff-sized constant residues are harmless.
				if q.Real == 1 && vector > 1e-12 {
					return false
				}
				continue
			}
			distance := 2 * math.Atan2(vector, math.Abs(q.Real))
			if !finiteBound(distance) || distance <= travel+1e-4 {
				return false
			}
		}
	}
	// Narrow-phase mesh checks also form relative geometry poses. Include
	// inter-arm and robot/world pairs, removing shared ancestor travel so
	// rigidly attached geometry does not force unnecessary fallback.
	for _, pair := range p.pairs {
		travel := world[pair.a].travel + world[pair.b].travel - 2*world[pair.common].travel
		q := quat.Mul(quat.Conj(world[pair.a].q), world[pair.b].q)
		vector := math.Sqrt(q.Imag*q.Imag + q.Jmag*q.Jmag + q.Kmag*q.Kmag)
		if travel <= 1e-14 {
			if q.Real == 1 && vector > 1e-12 {
				return false
			}
			continue
		}
		distance := 2 * math.Atan2(vector, math.Abs(q.Real))
		if !finiteBound(distance) || distance <= travel+1e-4 {
			return false
		}
	}
	return true
}
