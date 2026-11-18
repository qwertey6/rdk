package armplanning

import (
	"context"
	"math"
	"math/rand"
	"path/filepath"
	"sync"
	"testing"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

func TestNativeGoalCloudFallbackAndPositionOnly(t *testing.T) {
	req, err := ReadRequestFromFile("data/native-dual-xarm6-crossing.json")
	if err != nil {
		t.Fatal(err)
	}
	goals, err := req.Goals[0].ComputePoses(context.Background(), req.FrameSystem)
	if err != nil {
		t.Fatal(err)
	}
	pc, err := NewPlanContext(context.Background(), logging.NewBlankLogger("fallback"), req, &PlanMeta{})
	if err != nil {
		t.Fatal(err)
	}
	q := req.StartState.LinearConfiguration().GetLinearizedInputs()
	pc.planOpts.GoalMetricType = motionplan.PositionOnly
	cost, _ := pc.goalCost(goals)
	want := pc.LinearizeFSMetric(pc.planOpts.GetGoalMetric(goals))(context.Background(), q)
	if got := cost(context.Background(), q); math.Abs(got-want) > 1e-8*math.Max(1, want) {
		t.Fatal("position-only metric differs", got, want)
	}
	for _, g := range goals {
		g.GoalCloud = &referenceframe.PoseCloud{X: 100, Y: 100, Z: 100}
	}
	cost, batch := pc.goalCost(goals)
	if batch != nil {
		t.Fatal("goal clouds must retain original evaluator")
	}
	want = pc.LinearizeFSMetric(pc.planOpts.GetGoalMetric(goals))(context.Background(), q)
	if got := cost(context.Background(), q); got != want {
		t.Fatal("goal cloud metric changed", got, want)
	}
}

// A winning IK worker cancels its peers. Their in-flight NLopt calls must still
// receive finite objectives while winding down; infinities produce NaN gradients
// and turn millisecond plans into waits for the optimizer's evaluation limit.
func TestNativeGoalCostAfterCancellation(t *testing.T) {
	req, err := ReadRequestFromFile("data/native-dual-xarm6-crossing.json")
	if err != nil {
		t.Fatal(err)
	}
	goals, err := req.Goals[0].ComputePoses(context.Background(), req.FrameSystem)
	if err != nil {
		t.Fatal(err)
	}
	pc, err := NewPlanContext(context.Background(), logging.NewBlankLogger("cancellation"), req, &PlanMeta{})
	if err != nil {
		t.Fatal(err)
	}
	q := append([]float64(nil), req.StartState.LinearConfiguration().GetLinearizedInputs()...)
	for i := range q {
		q[i] += .123
	}
	flat := append(append([]float64(nil), q...), q...)
	flat[len(q)] += 1e-8
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	cost, batch := pc.goalCost(goals)
	if batch == nil {
		t.Fatalf("did not prepare an evaluator")
	}
	want, got := make([]float64, 2), make([]float64, 2)
	batch(context.Background(), flat, 2, want)
	batch(canceled, flat, 2, got)
	for i := range got {
		if math.IsNaN(got[i]) || math.IsInf(got[i], 0) || got[i] != want[i] {
			t.Fatalf("changed objective after cancellation: %v vs %v", got, want)
		}
	}
	if got := cost(canceled, q); got != want[0] {
		t.Fatalf("changed scalar objective after cancellation: %g vs %g", got, want[0])
	}
	var workers sync.WaitGroup
	for range 8 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			var scores [2]float64
			for range 100 {
				batch(canceled, flat, 2, scores[:])
				if scores[0] != want[0] || scores[1] != want[1] {
					t.Errorf("concurrent batch objective changed")
				}
				if got := cost(canceled, q); got != want[0] {
					t.Errorf("concurrent scalar objective changed")
				}
			}
		}()
	}
	workers.Wait()
}

func TestAutomaticGoalPreservesNativeArithmetic(t *testing.T) {
	req, err := ReadRequestFromFile("data/native-dual-xarm6-crossing.json")
	if err != nil {
		t.Fatal(err)
	}
	goals, err := req.Goals[0].ComputePoses(context.Background(), req.FrameSystem)
	if err != nil {
		t.Fatal(err)
	}
	pc, err := NewPlanContext(context.Background(), logging.NewBlankLogger("global-fk"), req, &PlanMeta{})
	if err != nil {
		t.Fatal(err)
	}
	cost, batch := pc.goalCost(goals)
	if batch == nil {
		t.Fatal("automatic IK must preserve native arithmetic")
	}
	q := req.StartState.LinearConfiguration().GetLinearizedInputs()
	reference := pc.LinearizeFSMetric(pc.planOpts.GetGoalMetric(goals))
	if got, want := cost(context.Background(), q), reference(context.Background(), q); got != want {
		t.Fatalf("global objective changed: got %g want %g", got, want)
	}
}

func TestNativeGoalKinematicsParity(t *testing.T) {
	for _, file := range []string{"data/native-dual-xarm6-crossing.json"} {
		t.Run(filepath.Base(file), func(t *testing.T) {
			req, err := ReadRequestFromFile(file)
			if err != nil {
				t.Fatal(err)
			}
			goals, err := req.Goals[0].ComputePoses(context.Background(), req.FrameSystem)
			if err != nil {
				t.Fatal(err)
			}
			pc, err := NewPlanContext(context.Background(), logging.NewBlankLogger("parity"), req, &PlanMeta{})
			if err != nil {
				t.Fatal(err)
			}
			reference := pc.LinearizeFSMetric(pc.planOpts.GetGoalMetric(goals))

			cost, batch := pc.goalCost(goals)
			rng := rand.New(rand.NewSource(17))
			limits := pc.lis.GetLimits()
			for sample := 0; sample < 40; sample++ {
				q := append([]float64(nil), req.StartState.LinearConfiguration().GetLinearizedInputs()...)
				if sample > 0 {
					for i, l := range limits {
						lo, hi := math.Max(l.Min, -3), math.Min(l.Max, 3)
						q[i] = lo + rng.Float64()*(hi-lo)
					}
				}
				want, got := reference(context.Background(), q), cost(context.Background(), q)

				inputs, err := pc.lis.FloatsToInputs(q)
				if err != nil {
					t.Fatal(err)
				}
				for name, goal := range goals {
					transform, err := pc.fs.PrepareTransform(pc.lis, name, goal.Parent())
					if err != nil {
						continue // Deliberate fallback for unsupported custom frames.
					}
					nativePose, err := pc.fs.TransformToDQ(inputs, name, goal.Parent())
					if err != nil {
						t.Fatal(err)
					}
					preparedPose, err := transform.Evaluate(q)
					if err != nil || preparedPose.Number != nativePose.Number {
						t.Fatalf("prepared FK arithmetic differs for %s at sample %d: got %v want %v, %v",
							name, sample, preparedPose.Number, nativePose.Number, err)
					}
				}

				if math.Abs(want-got) > 1e-8*math.Max(1, math.Abs(want)) {
					t.Fatalf("metric got %.12g want %.12g", got, want)
				}
				if batch == nil {
					continue
				}
				flat := append(append([]float64(nil), q...), q...)
				scores := make([]float64, 2)
				for j := range q {
					if q[j]+1e-8 >= limits[j].Max {
						continue
					}
					flat[len(q)+j] += 1e-8
					batch(context.Background(), flat, 2, scores)
					wantD := (reference(context.Background(), flat[len(q):]) - want) / 1e-8
					gotD := (scores[1] - scores[0]) / 1e-8
					// Arithmetic order changes roundoff in an objective divided by
					// a 1e-8 probe. Bound that noise explicitly in objective ULPs.
					noise := 64 * (math.Nextafter(math.Abs(want), math.Inf(1)) - math.Abs(want)) / 1e-8
					if math.Abs(wantD-gotD) > math.Max(.2, noise)+2e-4*math.Abs(wantD) {
						t.Fatalf("gradient joint=%d got=%g want=%g", j, gotD, wantD)
					}
					flat[len(q)+j] = q[j]
				}
			}
		})
	}
}
