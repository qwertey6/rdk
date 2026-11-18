//go:build !windows && !no_cgo

package ik

import (
	"context"
	"errors"
	"math"
	"testing"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
)

func TestBatchGradientMatchesScalarAtBounds(t *testing.T) {
	cost := func(_ context.Context, q []float64) float64 { return q[0]*q[0] + 3*q[1]*q[1] + q[0]*q[1] }
	for _, q := range [][]float64{{.3, -.2}, {1, 1}, {0, 0}, {1 - 1e-9, -1}} {
		iterations := 0
		scalar := nloptSeedState{jump: []float64{1e-8, 1e-8}, upperBound: []float64{1, 1}, logger: logging.NewBlankLogger("batch")}
		batch := scalar
		batch.batchInputs = make([]float64, 6)
		batch.batchScores = make([]float64, 3)
		batch.batchCost = func(ctx context.Context, flat []float64, count int, scores []float64) {
			for i := range count {
				scores[i] = cost(ctx, flat[i*2:(i+1)*2])
			}
		}
		before := append([]float64(nil), q...)
		g1, g2 := make([]float64, 2), make([]float64, 2)
		v1 := scalar.getMinFunc(context.Background(), cost, &iterations)(q, g1)
		v2 := batch.getMinFunc(context.Background(), cost, &iterations)(q, g2)
		if math.Abs(v1-v2) > 1e-14 || math.Abs(g1[0]-g2[0]) > 1e-7 || math.Abs(g1[1]-g2[1]) > 1e-7 || math.IsNaN(v2) {
			t.Fatalf("different gradients: %v vs %v", g1, g2)
		}
		if q[0] != before[0] || q[1] != before[1] {
			t.Fatal("gradient mutated input")
		}
	}
}

func TestBatchNloptCancellationStopsInFlightSolve(t *testing.T) {
	for _, batched := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		calls := 0
		cost := func(_ context.Context, q []float64) float64 {
			calls++
			cancel()
			return (q[0]-0.5)*(q[0]-0.5) + 100
		}
		solver, err := CreateNloptSolver(logging.NewBlankLogger("cancel"), 1, false, false, 0)
		if err != nil {
			t.Fatal(err)
		}
		if batched {
			solver.SetBatchCost(func(ctx context.Context, q []float64, count int, scores []float64) {
				for i := range count {
					scores[i] = cost(ctx, q[i:i+1])
				}
			})
		}
		result, err := solver.SolveOnce(ctx, cost, []float64{0}, []referenceframe.Limit{{Min: -1, Max: 1}})
		if !errors.Is(err, context.Canceled) || result != nil {
			t.Fatalf("batched=%v: canceled solve returned %v, %v", batched, result, err)
		}
		limit := 1
		if batched {
			limit = 2 // Finish only the batch already in progress.
		}
		if calls > limit {
			t.Fatalf("batched=%v: %d objective calls after cancellation, limit %d", batched, calls, limit)
		}
	}
}
