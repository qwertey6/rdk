package armplanning

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

func TestGuardedFinishingPreservesFallback(t *testing.T) {
	for _, radius := range []float64{20, 80} {
		t.Run(fmt.Sprint(radius), func(t *testing.T) {
			x, err := referenceframe.NewTranslationalFrame("x", r3.Vector{X: 1}, referenceframe.Limit{Min: -300, Max: 300})
			test.That(t, err, test.ShouldBeNil)
			y, err := referenceframe.NewTranslationalFrame("y", r3.Vector{Y: 1}, referenceframe.Limit{Min: -300, Max: 300})
			test.That(t, err, test.ShouldBeNil)
			model, err := referenceframe.NewSerialModel("xy", []referenceframe.Frame{x, y})
			test.That(t, err, test.ShouldBeNil)
			fs := referenceframe.NewEmptyFrameSystem("guarded")
			test.That(t, fs.AddFrame(model, fs.World()), test.ShouldBeNil)
			q := func(x, y float64) *referenceframe.LinearInputs {
				return referenceframe.FrameSystemInputs{"xy": {x, y}}.ToLinearInputs()
			}
			checker := motionplan.NewEmptyConstraintChecker(logging.NewTestLogger(t))
			checker.SetCollisionConstraints(motionplan.CollisionConstraints{Obstacle: func(s *motionplan.StateFS) (float64, error) {
				v := s.Configuration.Get("xy")
				d := math.Hypot(v[0], v[1]) - radius
				if d <= 0 {
					return d, fmt.Errorf("inside obstacle")
				}
				return d, nil
			}})
			meta := &PlanMeta{}
			psc := &PlanSegmentContext{pc: &PlanContext{
				fs: fs, planOpts: NewBasicPlannerOptions(), planMeta: meta,
				ConfigurationDistanceFunc: motionplan.FSConfigurationL2Distance,
			}, Checker: checker}
			height := 35.
			if radius == 80 {
				height = 180
			}
			corners := []*referenceframe.LinearInputs{q(-100, 0), q(0, height), q(100, 0)}
			steps := []*referenceframe.LinearInputs{corners[0]}
			for i := 1; i < len(corners); i++ {
				for j := 1; j <= 20; j++ {
					v, err := referenceframe.InterpolateFS(fs, corners[i-1], corners[i], float64(j)/20)
					test.That(t, err, test.ShouldBeNil)
					steps = append(steps, v)
				}
			}
			saved := make([]*referenceframe.LinearInputs, len(steps))
			for i, q := range steps {
				saved[i] = q.Copy()
			}
			budget := &finishingBudget{remaining: 12000}
			ctx := context.WithValue(context.Background(), finishingBudgetKey{}, budget)
			cleanupCtx := context.WithValue(ctx, finishingBudgetKey{}, (*finishingBudget)(nil))
			trial, accepted := tryGuardedFinishing(ctx, cleanupCtx, psc, steps)
			test.That(t, accepted, test.ShouldEqual, radius == 20)
			for i := range steps {
				test.That(t, steps[i].Get("xy"), test.ShouldResemble, saved[i].Get("xy"))
			}

			if accepted {
				test.That(t, budget.remaining, test.ShouldBeLessThan, 12000)
				for i := 1; i < len(trial); i++ {
					test.That(t, psc.CheckPath(ctx, trial[i-1], trial[i], true, nil), test.ShouldBeNil)
				}
			} else {
				test.That(t, budget.remaining, test.ShouldEqual, 12000)
			}
		})
	}
}

func TestGuardedFinishingProbeLimitAllRejected(t *testing.T) {
	x, err := referenceframe.NewTranslationalFrame("x", r3.Vector{X: 1}, referenceframe.Limit{Min: -100, Max: 100})
	test.That(t, err, test.ShouldBeNil)
	fs := referenceframe.NewEmptyFrameSystem("bounded-trial")
	test.That(t, fs.AddFrame(x, fs.World()), test.ShouldBeNil)
	calls := 0
	checker := motionplan.NewEmptyConstraintChecker(logging.NewTestLogger(t))
	checker.SetCollisionConstraints(motionplan.CollisionConstraints{Obstacle: func(*motionplan.StateFS) (float64, error) {
		calls++
		return -1, fmt.Errorf("reject every speculative chord")
	}})
	meta := &PlanMeta{}
	psc := &PlanSegmentContext{pc: &PlanContext{
		fs: fs, planOpts: NewBasicPlannerOptions(), planMeta: meta,
		ConfigurationDistanceFunc: motionplan.FSConfigurationL2Distance,
	}, Checker: checker}
	var steps []*referenceframe.LinearInputs
	for i := 0; i <= 100; i++ {
		v := float64(i / 2)
		if i%2 != 0 {
			v += 20
		}
		steps = append(steps, referenceframe.FrameSystemInputs{"x": {v}}.ToLinearInputs())
	}
	budget := &finishingBudget{remaining: 1234}
	ctx := context.WithValue(context.Background(), finishingBudgetKey{}, budget)
	cleanupCtx := context.WithValue(ctx, finishingBudgetKey{}, (*finishingBudget)(nil))
	trial, accepted := tryGuardedFinishing(ctx, cleanupCtx, psc, steps)
	test.That(t, accepted, test.ShouldBeFalse)
	test.That(t, trial, test.ShouldBeNil)
	test.That(t, calls, test.ShouldEqual, guardedFinishingProbeBudget)
	test.That(t, budget.remaining, test.ShouldEqual, 1234)
}
