package armplanning

import (
	"context"
	"slices"

	"go.viam.com/rdk/referenceframe"
)

// Keep speculation small. Rejected trials still leave the ordinary finisher
// its full allowance; this cap limits the extra work paid before that fallback.
const guardedFinishingProbeBudget = 32

// tryGuardedFinishing speculates on a fast compact path. Failed trials leave the
// raw configurations and original fallback allowance intact; all work is counted.
func tryGuardedFinishing(ctx, cleanupCtx context.Context, psc *PlanSegmentContext,
	steps []*referenceframe.LinearInputs,
) ([]*referenceframe.LinearInputs, bool) {
	if len(steps) <= 15 || ctx.Err() != nil {
		return nil, false
	}
	quality := newShortcutQuality(psc, []*referenceframe.LinearInputs{steps[0], steps[len(steps)-1]})
	if !quality.enabled {
		return nil, false
	}
	trialBudget := guardedFinishingProbeBudget
	trialCtx := ctx
	var trialSamples *finishingBudget
	baseSamples, _ := ctx.Value(finishingBudgetKey{}).(*finishingBudget)
	if baseSamples != nil {
		trialSamples = &finishingBudget{remaining: baseSamples.remaining}
		trialCtx = context.WithValue(ctx, finishingBudgetKey{}, trialSamples)
	}
	// Partitioning owns its output slice. All later vertex/idle updates
	// replace configurations with fresh copies, preserving the raw fallback.
	trial := partitionShortcuts(cleanupCtx, psc, slices.Clone(steps), &trialBudget)
	// Pulls reserve two units at a time. Save an odd remainder for the final
	// single-edge removal pass so the trial cannot overdraw its allowance.
	reserved := trialBudget % 2
	pullBudget := trialBudget - reserved
	trial = vertexPullStep(trialCtx, psc, trial, &pullBudget)
	trialBudget = pullBudget + reserved
	trial = simpleSmoothStep(trialCtx, psc, trial, 1, &trialBudget)
	accepted := ctx.Err() == nil && quality.satisfied(trial)
	if baseSamples != nil {
		// Preserve the original fallback allowance for rejected trials.
		// Accepted trials retain their remaining allowance.
		if accepted {
			baseSamples.remaining = trialSamples.remaining
		}
	}

	if accepted {
		return trial, true
	}
	return nil, false
}
