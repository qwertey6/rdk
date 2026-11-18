package armplanning

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func TestCheckPathCacheScope(t *testing.T) {
	logger := logging.NewTestLogger(t)
	fs := referenceframe.NewEmptyFrameSystem("cache")
	joint, err := referenceframe.NewRotationalFrame("joint", spatialmath.R4AA{RZ: 1}, referenceframe.Limit{Min: -2, Max: 2})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(joint, fs.World()), test.ShouldBeNil)
	q := func(v float64) *referenceframe.LinearInputs {
		return referenceframe.FrameSystemInputs{"joint": {v}}.ToLinearInputs()
	}
	pc := &PlanContext{fs: fs, planOpts: NewBasicPlannerOptions(), collisionCache: motionplan.NewCollisionCache()}
	checker := motionplan.NewEmptyConstraintChecker(logger)
	checker.SetCollisionConstraints(motionplan.CollisionConstraints{Obstacle: func(s *motionplan.StateFS) (float64, error) {
		if s.Configuration.Get("joint")[0] == 1 {
			return -1, fmt.Errorf("blocked endpoint")
		}
		return 0, nil
	}})
	psc := &PlanSegmentContext{pc: pc, Checker: checker}
	test.That(t, psc.CheckPath(context.Background(), q(0), q(1), false, nil), test.ShouldBeNil)
	// A prefix check must not certify its unchecked endpoint or the reverse edge.
	test.That(t, psc.CheckPath(context.Background(), q(0), q(1), true, nil), test.ShouldNotBeNil)
	test.That(t, psc.CheckPath(context.Background(), q(1), q(0), false, nil), test.ShouldNotBeNil)
	// A new segment can impose a different constraint even in the same plan.
	other := motionplan.NewEmptyConstraintChecker(logger)
	other.SetCollisionConstraints(motionplan.CollisionConstraints{Obstacle: func(*motionplan.StateFS) (float64, error) {
		return -1, fmt.Errorf("new segment constraint")
	}})
	psc2 := &PlanSegmentContext{pc: pc, Checker: other}
	test.That(t, psc2.CheckPath(context.Background(), q(0), q(1), false, nil), test.ShouldNotBeNil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	test.That(t, psc.CheckPath(ctx, q(0), q(1), false, nil), test.ShouldNotBeNil)
}

func nativeFixture(t *testing.T) *PlanRequest {
	t.Helper()
	req, err := ReadRequestFromFile("data/native-dual-xarm6-crossing.json")
	test.That(t, err, test.ShouldBeNil)
	req.PlannerOptions.Timeout = 10
	return req
}

func TestCheckPathDetectsThinCrossing(t *testing.T) {
	req := nativeFixture(t)
	data, err := os.ReadFile("data/native-thin-crossing-edge.json")
	test.That(t, err, test.ShouldBeNil)
	var edge struct {
		From, To referenceframe.FrameSystemInputs
	}
	test.That(t, json.Unmarshal(data, &edge), test.ShouldBeNil)
	pc, err := NewPlanContext(context.Background(), logging.NewTestLogger(t), req, &PlanMeta{})
	test.That(t, err, test.ShouldBeNil)
	goal, err := req.Goals[0].ComputePoses(context.Background(), req.FrameSystem)
	test.That(t, err, test.ShouldBeNil)
	psc, err := NewPlanSegmentContext(context.Background(), pc, req.StartState.LinearConfiguration(), goal)
	test.That(t, err, test.ShouldBeNil)
	path := []*referenceframe.LinearInputs{edge.From.ToLinearInputs(), edge.To.ToLinearInputs()}
	for _, q := range path {
		_, err = psc.Checker.CheckStateFSConstraints(context.Background(), &motionplan.StateFS{FS: req.FrameSystem, Configuration: q})
		test.That(t, err, test.ShouldBeNil)
	}
	test.That(t, validateDenseEdges(context.Background(), psc, path), test.ShouldNotBeNil)
	test.That(t, psc.CheckPath(context.Background(), path[0], path[1], true, nil), test.ShouldNotBeNil)
}

// Dense regression oracle bypasses edge caching and clearance skipping.
func validateDenseEdges(ctx context.Context, psc *PlanSegmentContext, path []*referenceframe.LinearInputs) error {
	for edge, q := range path {
		previous := q
		if edge > 0 {
			previous = path[edge-1]
		}
		steps := 1
		for name, vals := range q.Items() {
			for j, v := range vals {
				steps = max(steps, int(math.Ceil(math.Abs(v-previous.Get(name)[j])/.002)))
			}
		}
		for k := 0; k <= steps; k++ {
			if err := ctx.Err(); err != nil {
				return err
			}
			inputs := referenceframe.NewLinearInputs()
			for name, vals := range q.Items() {
				interp := make([]float64, len(vals))
				for j, v := range vals {
					interp[j] = previous.Get(name)[j] + (v-previous.Get(name)[j])*float64(k)/float64(steps)
				}
				inputs.Put(name, interp)
			}
			if _, err := psc.Checker.CheckStateFSConstraints(ctx, &motionplan.StateFS{FS: psc.pc.fs, Configuration: inputs}); err != nil {
				return fmt.Errorf("RDK rejected edge %d sample %d/%d: %w", edge, k, steps, err)
			}
		}
	}
	return nil
}
