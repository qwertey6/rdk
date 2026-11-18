package armplanning

import (
	"context"
	"math"
	"sync"

	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/motionplan/ik"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

type preparedGoal struct {
	pose      spatialmath.Pose
	transform *referenceframe.PreparedTransform
}

type goalCostScratch struct {
	pose spatialmath.DualQuaternion
}

func (pc *PlanContext) goalCost(goals referenceframe.FrameSystemPoses) (ik.CostFunc, ik.BatchCostFunc) {
	fallback := pc.LinearizeFSMetric(pc.planOpts.GetGoalMetric(goals))
	var prepared []preparedGoal
	for name, goal := range goals {
		// Goal clouds have their own region membership semantics. Keep the full
		// existing metric until a separately validated accelerated evaluator exists.
		if goal.GoalCloud != nil {
			return fallback, nil
		}
		transform, err := pc.fs.PrepareTransform(pc.lis, name, goal.Parent())
		if err != nil {
			return fallback, nil
		}
		g := preparedGoal{pose: goal.Pose(), transform: transform}
		prepared = append(prepared, g)
	}

	dimension := len(pc.lis.GetLimits())
	orientationScale := 100.
	if pc.planOpts.GoalMetricType == motionplan.PositionOnly {
		orientationScale = 0
	}
	// Scratch belongs to one invocation, never shared between concurrent IK workers.
	pool := sync.Pool{New: func() any { return new(goalCostScratch) }}
	batch := func(ctx context.Context, inputs []float64, count int, scores []float64) {
		// Cancellation belongs to the IK solve loop, as for the original metric.
		// Returning infinities here gives NLopt NaN gradients and can keep losing
		// workers iterating until MaxEval while the planner waits for shutdown.
		clear(scores[:count])
		storage := pool.Get().(*goalCostScratch)
		defer pool.Put(storage)
		for _, g := range prepared {
			for lane := range count {
				// The metric takes a Pose interface. A loop-local value escapes
				// once per lane; keep it in this invocation's private scratch.
				pose := &storage.pose
				var err error
				*pose, err = g.transform.Evaluate(inputs[lane*dimension : (lane+1)*dimension])
				if err != nil {
					scores[lane] = math.NaN()
					continue
				}
				scores[lane] += motionplan.WeightedSquaredNormDistanceWithTolerance(g.pose, pose, 1, orientationScale, 0)
			}
		}
		for lane := range count {
			if math.IsNaN(scores[lane]) {
				scores[lane] = fallback(ctx, inputs[lane*dimension:(lane+1)*dimension])
			}
		}
	}
	scalar := func(ctx context.Context, inputs []float64) float64 {
		if len(inputs) != dimension {
			return fallback(ctx, inputs)
		}
		var result [1]float64
		batch(ctx, inputs, 1, result[:])
		return result[0]
	}
	return scalar, batch
}
