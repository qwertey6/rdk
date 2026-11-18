package armplanning

import (
	"context"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

func TestPartitionShortcutsValidatedAndBounded(t *testing.T) {
	ctx := context.Background()
	fs, _, goalJoints, req := nudgeBlockedScene(t)
	logger := logging.NewTestLogger(t)
	pc, err := NewPlanContext(ctx, logger, req, &PlanMeta{})
	test.That(t, err, test.ShouldBeNil)
	goals, err := req.Goals[0].ComputePoses(ctx, fs)
	test.That(t, err, test.ShouldBeNil)
	psc, err := NewPlanSegmentContext(ctx, pc, req.StartState.LinearConfiguration(), goals)
	test.That(t, err, test.ShouldBeNil)
	// Fix the goal configuration so IK cannot choose an unobstructed chord.
	goal := referenceframe.FrameSystemInputs{"arm": goalJoints}.ToLinearInputs()
	trajectory := tryNudgedStraightLine(ctx, psc, rrtMap{&node{inputs: goal}: nil}, pc.randseed, logger)
	test.That(t, len(trajectory), test.ShouldBeGreaterThanOrEqualTo, 3)
	steps := []*referenceframe.LinearInputs{trajectory[0]}
	for i := 1; i < len(trajectory); i++ {
		from, to := trajectory[i-1], trajectory[i]
		for sample := 1; sample <= 12; sample++ {
			q, err := referenceframe.InterpolateFS(fs, from, to, float64(sample)/12)
			test.That(t, err, test.ShouldBeNil)
			steps = append(steps, q)
		}
	}
	test.That(t, len(steps), test.ShouldBeGreaterThan, 15)
	test.That(t, psc.CheckPath(ctx, steps[0], steps[len(steps)-1], true, nil), test.ShouldNotBeNil)
	budget := smoothProbeBudget
	got := partitionShortcuts(ctx, psc, steps, &budget)
	test.That(t, len(got), test.ShouldBeLessThan, len(steps))
	test.That(t, len(got), test.ShouldBeGreaterThanOrEqualTo, 3)
	test.That(t, got[0], test.ShouldResemble, steps[0])
	test.That(t, got[len(got)-1], test.ShouldResemble, steps[len(steps)-1])
	test.That(t, finishingPathCost(psc, got), test.ShouldBeLessThanOrEqualTo, finishingPathCost(psc, steps)+1e-8)
	for i := 1; i < len(got); i++ {
		test.That(t, psc.CheckPath(ctx, got[i-1], got[i], true, nil), test.ShouldBeNil)
	}
	budget = 0
	test.That(t, partitionShortcuts(ctx, psc, steps, &budget), test.ShouldResemble, steps)
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	budget = smoothProbeBudget
	test.That(t, partitionShortcuts(canceled, psc, steps, &budget), test.ShouldResemble, steps)
	test.That(t, budget, test.ShouldEqual, smoothProbeBudget)
}

func finishingPathCost(psc *PlanSegmentContext, steps []*referenceframe.LinearInputs) float64 {
	cost := 0.
	for i := 1; i < len(steps); i++ {
		cost += psc.pc.ConfigurationDistanceFunc(&motionplan.SegmentFS{
			StartConfiguration: steps[i-1], EndConfiguration: steps[i],
		})
	}
	return cost
}
