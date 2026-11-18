package armplanning

import (
	"context"

	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

// partitionShortcuts simplifies long raw paths by trying their outer chord,
// then splitting blocked chords at the largest configuration-space detour.
// Every accepted shortcut receives the same full-resolution path check as the
// ordinary removal pass. Exhausted work or cancellation retains the raw span.
func partitionShortcuts(ctx context.Context, psc *PlanSegmentContext,
	steps []*referenceframe.LinearInputs, budget *int,
) []*referenceframe.LinearInputs {
	if len(steps) < 3 {
		return steps
	}
	type span struct{ first, last int }
	pending := []span{{0, len(steps) - 1}}
	out := make([]*referenceframe.LinearInputs, 0, len(steps))
	out = append(out, steps[0])
	cost := func(a, b int) float64 {
		return psc.pc.ConfigurationDistanceFunc(&motionplan.SegmentFS{
			StartConfiguration: steps[a], EndConfiguration: steps[b],
		})
	}
	for len(pending) > 0 {
		part := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if part.last-part.first < 2 || *budget <= 0 || ctx.Err() != nil {
			out = append(out, steps[part.first+1:part.last+1]...)
			continue
		}
		*budget--
		if checkSmoothingPath(ctx, psc, steps[part.first], steps[part.last], false) == nil {
			out = append(out, steps[part.last])
			continue
		}
		if ctx.Err() != nil {
			out = append(out, steps[part.first+1:part.last+1]...)
			continue
		}
		// Start at the midpoint so equal-cost straight spans split evenly.
		// The direct chord cost is constant and need not be subtracted.
		split := (part.first + part.last) / 2
		largest := cost(part.first, split) + cost(split, part.last)
		for i := part.first + 1; i < part.last; i++ {
			if detour := cost(part.first, i) + cost(i, part.last); detour > largest+1e-12 {
				split, largest = i, detour
			}
		}
		// LIFO processes the left span first, preserving trajectory order.
		pending = append(pending, span{split, part.last}, span{part.first, split})
	}
	return out
}
