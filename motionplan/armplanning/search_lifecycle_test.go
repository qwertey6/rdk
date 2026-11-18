package armplanning

import (
	"context"
	"testing"

	"go.uber.org/goleak"

	"go.viam.com/rdk/logging"
)

func TestNativeSearchWorkerLifetime(t *testing.T) {
	testNativeSearchWorkerLifetime(t)
}

func TestNativeSearchWorkerOverrideLifetime(t *testing.T) {
	previous := cbirrtRaceAttempts
	cbirrtRaceAttempts = maxTotalAttempts + 6
	t.Cleanup(func() { cbirrtRaceAttempts = previous })
	testNativeSearchWorkerLifetime(t)
}

func testNativeSearchWorkerLifetime(t *testing.T) {
	t.Helper()
	for _, canceled := range []bool{false, true} {
		name := "solution"
		if canceled {
			name = "canceled"
		}
		t.Run(name, func(t *testing.T) {
			defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
			req := nativeFixture(t)
			logger := logging.NewBlankLogger("search-lifetime")
			meta := &PlanMeta{}
			pc, err := NewPlanContext(context.Background(), logger, req, meta)
			if err != nil {
				t.Fatal(err)
			}
			start := req.StartState.LinearConfiguration()
			goalPoses, err := req.Goals[0].ComputePoses(context.Background(), req.FrameSystem)
			if err != nil {
				t.Fatal(err)
			}
			psc, err := NewPlanSegmentContext(context.Background(), pc, start, goalPoses)
			if err != nil {
				t.Fatal(err)
			}
			goal := &node{inputs: req.Goals[0].LinearConfiguration()}
			maps := &rrtMaps{startMap: rrtMap{&node{inputs: start}: nil}, goalMap: rrtMap{goal: nil}, optNode: goal}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if canceled {
				cancel()
			}
			pm := &planManager{pc: pc, request: req, logger: logger}
			steps, _, err := pm.raceCBiRRT(ctx, psc, maps, nil, false)
			if !canceled && (err != nil || len(steps) == 0) {
				t.Fatalf("search failed: %v", err)
			}
			if canceled && err == nil {
				t.Fatal("canceled search returned a solution")
			}
		})
	}
}
