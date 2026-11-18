package armplanning

import (
	"context"
	"math"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

func TestShortcutQualityRequiresBothLengthBounds(t *testing.T) {
	psc := &PlanSegmentContext{pc: &PlanContext{
		planOpts:                  NewBasicPlannerOptions(),
		ConfigurationDistanceFunc: motionplan.FSConfigurationL2Distance,
	}}
	q := func(left, right float64) *referenceframe.LinearInputs {
		return referenceframe.FrameSystemInputs{"left": {left}, "right": {right}}.ToLinearInputs()
	}
	// Sequential arm movement has optimal summed per-arm distance, but its
	// whole-configuration path is 41% longer than simultaneous movement.
	steps := []*referenceframe.LinearInputs{q(0, 0), q(1, 0), q(1, 1)}
	quality := newShortcutQuality(psc, steps)
	test.That(t, quality.satisfied(steps), test.ShouldBeFalse)
	steps = []*referenceframe.LinearInputs{steps[0], steps[2]}
	quality = newShortcutQuality(psc, steps)
	test.That(t, quality.satisfied(steps), test.ShouldBeTrue)
	// A small reversal can fit the whole-configuration bound while exceeding
	// the planner's summed per-arm metric bound. Both must pass.
	steps = []*referenceframe.LinearInputs{q(0, 0), q(.5, .1), q(1, 0)}
	quality = newShortcutQuality(psc, steps)
	test.That(t, quality.length, test.ShouldBeLessThan, quality.lengthLimit)
	test.That(t, quality.satisfied(steps), test.ShouldBeFalse)
	for _, endpoint := range []float64{0, math.NaN(), math.Inf(1)} {
		steps = []*referenceframe.LinearInputs{q(0, 0), q(endpoint, 0)}
		test.That(t, newShortcutQuality(psc, steps).satisfied(steps), test.ShouldBeFalse)
	}
	psc.pc.planOpts.ConfigurationDistanceMetric = motionplan.FSConfigurationDistanceMetric
	steps = []*referenceframe.LinearInputs{q(0, 0), q(1, 1)}
	test.That(t, newShortcutQuality(psc, steps).satisfied(steps), test.ShouldBeFalse)
}

func TestShortcutQualityKeepsMandatoryCollisionValidation(t *testing.T) {
	ctx := context.Background()
	fs, startJoints, goalJoints, req := nudgeBlockedScene(t)
	logger := logging.NewTestLogger(t)
	meta := &PlanMeta{}
	pc, err := NewPlanContext(ctx, logger, req, meta)
	test.That(t, err, test.ShouldBeNil)
	goals, err := req.Goals[0].ComputePoses(ctx, fs)
	test.That(t, err, test.ShouldBeNil)
	psc, err := NewPlanSegmentContext(ctx, pc, req.StartState.LinearConfiguration(), goals)
	test.That(t, err, test.ShouldBeNil)
	// Use the known blocked joint-space chord. IK may return a different goal
	// configuration whose direct path is unobstructed.
	steps := []*referenceframe.LinearInputs{
		referenceframe.FrameSystemInputs{"arm": startJoints}.ToLinearInputs(),
		referenceframe.FrameSystemInputs{"arm": goalJoints}.ToLinearInputs(),
	}
	test.That(t, psc.CheckPath(ctx, steps[0], steps[1], true, nil), test.ShouldNotBeNil)
	test.That(t, newShortcutQuality(psc, steps).satisfied(steps), test.ShouldBeTrue)
	// This chord is as short as possible, yet collides. A length target is
	// never evidence of validity; final full-resolution checking must reject.
	_, _, err = smoothPath(ctx, psc, steps)
	test.That(t, err, test.ShouldNotBeNil)
}
