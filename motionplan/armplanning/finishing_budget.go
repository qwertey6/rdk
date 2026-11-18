package armplanning

import (
	"context"
	"errors"

	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

type (
	finishingBudgetKey struct{}
	finishingBudget    struct{ remaining int }
)

var errFinishingBudget = errors.New("optional finishing validation budget exhausted")

// Budget units are the original dense sample count of candidate edges. Charging
// before validation bounds work even when an edge is expensive. This is private
// to optional finishing; the final validation sweep never consumes this budget.
func newFinishingBudget(psc *PlanSegmentContext, steps []*referenceframe.LinearInputs) *finishingBudget {
	work := 0
	for i := 1; i < len(steps); i++ {
		n, err := motionplan.SegmentSampleCount(&motionplan.SegmentFS{
			FS:                 psc.pc.fs,
			StartConfiguration: steps[i-1], EndConfiguration: steps[i],
		}, psc.pc.planOpts.Resolution)
		if err != nil {
			return &finishingBudget{}
		}
		work += n
	}
	return &finishingBudget{remaining: min(12000, max(1000, 2*work))}
}

func checkSmoothingPath(ctx context.Context, psc *PlanSegmentContext,
	start, end *referenceframe.LinearInputs, checkFinal bool,
) error {
	if budget, _ := ctx.Value(finishingBudgetKey{}).(*finishingBudget); budget != nil {
		count, err := motionplan.SegmentSampleCount(&motionplan.SegmentFS{
			FS:                 psc.pc.fs,
			StartConfiguration: start, EndConfiguration: end,
		}, psc.pc.planOpts.Resolution)
		if err != nil {
			return err
		}
		if !checkFinal {
			count--
		}
		if count > budget.remaining {
			return errFinishingBudget
		}
		budget.remaining -= count
	}
	return psc.checkPath(ctx, start, end, checkFinal, nil, true)
}
