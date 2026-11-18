//go:build !windows && !no_cgo

package ik

import (
	"context"
	"errors"
	"testing"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
)

func TestNloptCancellationStopsInFlightSolve(t *testing.T) {
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
	result, err := solver.SolveOnce(ctx, cost, []float64{0}, []referenceframe.Limit{{Min: -1, Max: 1}})
	if !errors.Is(err, context.Canceled) || result != nil {
		t.Fatalf("canceled solve returned %v, %v", result, err)
	}
	if calls != 1 {
		t.Fatalf("%d objective calls after cancellation, want 1", calls)
	}
}
