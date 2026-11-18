package armplanning

import (
	"context"
	"fmt"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func TestFinishingBudgetLeavesMandatoryValidation(t *testing.T) {
	fs := referenceframe.NewEmptyFrameSystem("budget")
	joint, err := referenceframe.NewRotationalFrame("joint", spatialmath.R4AA{RZ: 1}, referenceframe.Limit{Min: -2, Max: 2})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(joint, fs.World()), test.ShouldBeNil)
	q := func(v float64) *referenceframe.LinearInputs {
		return referenceframe.FrameSystemInputs{"joint": {v}}.ToLinearInputs()
	}
	checker := motionplan.NewEmptyConstraintChecker(logging.NewTestLogger(t))
	checks := 0
	checker.SetCollisionConstraints(motionplan.CollisionConstraints{Obstacle: func(s *motionplan.StateFS) (float64, error) {
		checks++
		if s.Configuration.Get("joint")[0] > .1 {
			return 0, fmt.Errorf("blocked")
		}
		return 100, nil
	}})
	psc := &PlanSegmentContext{pc: &PlanContext{fs: fs, planOpts: NewBasicPlannerOptions()}, Checker: checker}
	budget := &finishingBudget{remaining: 1}
	ctx := context.WithValue(context.Background(), finishingBudgetKey{}, budget)
	test.That(t, checkSmoothingPath(ctx, psc, q(0), q(1), true), test.ShouldEqual, errFinishingBudget)
	test.That(t, checks, test.ShouldEqual, 0)
	_, blocked, err := findCloseObstacleWaypoints(ctx, psc, q(0), q(1))
	test.That(t, err, test.ShouldBeNil)
	test.That(t, blocked, test.ShouldBeTrue)
	test.That(t, checks > 0, test.ShouldBeTrue)
}
