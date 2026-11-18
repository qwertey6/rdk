package motionplan

import (
	"fmt"
	"sync"
	"testing"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/spatialmath"
)

func TestCollisionAllowancesConcurrent(t *testing.T) {
	for _, count := range []int{3, 6, 7, 12, 64, 65} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			geometries := make([]spatialmath.Geometry, count)
			var allows []Collision
			for i := range geometries {
				g, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, fmt.Sprintf("g%02d", i))
				if err != nil {
					t.Fatal(err)
				}
				geometries[i] = g
				for j := 0; j < i; j++ {
					if j != 0 || i != count-1 {
						// Reversed names must still match canonical collision pairs.
						allows = append(allows, Collision{g.Label(), geometries[j].Label()})
					}
				}
			}
			allowance := makeAllowedCollisionsLookup(allows)
			logger := logging.NewTestLogger(t)
			var workers sync.WaitGroup
			for range 8 {
				workers.Add(1)
				go func() {
					defer workers.Done()
					for range 20 {
						got, _, err := checkCollisionsHinted(geometries, geometries, nil, allowance,
							defaultCollisionBufferMM, true, nil, logger)
						want := Collision{geometries[0].Label(), geometries[count-1].Label()}
						if err != nil || len(got) != 1 || !collisionsEqual(got[0], want) {
							t.Errorf("expected only %v, got %v, error %v", want, got, err)
						}
					}
				}()
			}
			workers.Wait()
		})
	}
}
