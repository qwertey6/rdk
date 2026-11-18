package armplanning

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/golang/geo/r3"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func TestInputFrameSmoothingRetainsMovingGeometryAndValidation(t *testing.T) {
	fs := referenceframe.NewEmptyFrameSystem("input-frame-smoothing")
	x, err := referenceframe.NewTranslationalFrame("x", r3.Vector{X: 1}, referenceframe.Limit{Min: -20, Max: 20})
	if err != nil {
		t.Fatal(err)
	}
	if err := fs.AddFrame(x, fs.World()); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"tool", "sensor"} {
		f, err := referenceframe.NewStaticFrame(name, spatialmath.NewZeroPose())
		if name == "tool" {
			g, gerr := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "tool-sphere")
			if gerr != nil {
				t.Fatal(gerr)
			}
			f, err = referenceframe.NewStaticFrameWithGeometry(name, spatialmath.NewZeroPose(), g)
		}
		if err != nil {
			t.Fatal(err)
		}
		if err := fs.AddFrame(f, x); err != nil {
			t.Fatal(err)
		}
		if !frameCanMove(fs, f) {
			t.Fatal("attached geometry must still count as moving for collision/fatal-error handling")
		}
	}
	var steps []*referenceframe.LinearInputs
	var saved []referenceframe.FrameSystemInputs
	for _, value := range []float64{-10, 10, -5} {
		q := referenceframe.FrameSystemInputs{"x": {value}, "tool": {}, "sensor": {}}.ToLinearInputs()
		steps = append(steps, q)
		saved = append(saved, q.Copy().ToFrameSystemInputs())
	}
	logger := logging.NewTestLogger(t)
	geometries, err := referenceframe.FrameSystemGeometriesLinearInputs(fs, steps[0])
	if err != nil {
		t.Fatal(err)
	}
	var robot []spatialmath.Geometry
	for _, gif := range geometries {
		robot = append(robot, gif.Geometries()...)
	}
	obstacle, err := spatialmath.NewSphere(spatialmath.NewZeroPose(), 1, "obstacle")
	if err != nil {
		t.Fatal(err)
	}
	cache := motionplan.NewCollisionCache()
	checker, err := motionplan.NewConstraintChecker(1e-8, nil, nil, nil, fs, robot, nil,
		map[string]bool{"x": true, "tool": true, "sensor": true}, steps[0],
		referenceframe.NewGeometriesInFrame(referenceframe.World, []spatialmath.Geometry{obstacle}), logger, cache)
	if err != nil || !checker.HasNativeMotionBounds() {
		t.Fatalf("native fixture unavailable: %v", err)
	}
	meta := &PlanMeta{}
	psc := &PlanSegmentContext{pc: &PlanContext{
		fs: fs, movableFrames: []string{"tool", "x", "sensor"},
		planOpts: NewBasicPlannerOptions(), planMeta: meta, logger: logger,
	}, Checker: checker}
	smoothMultiArms(context.Background(), psc, steps)

	for i, step := range steps {
		if !reflect.DeepEqual(step.ToFrameSystemInputs(), saved[i]) {
			t.Fatal("optional no-op changed configuration or frame keys")
		}
	}
	if _, err := addCloseObstacleWaypoints(context.Background(), psc, steps); err == nil {
		t.Fatal("mandatory final validation must still reject the unchanged obstructed path")
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := addCloseObstacleWaypoints(canceled, psc, steps); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
	owners, ok := smoothingInputFrames(fs, psc.pc.movableFrames)
	if !ok || !reflect.DeepEqual(owners, []string{"x"}) {
		t.Fatalf("input owners: %v, valid=%v", owners, ok)
	}
	if _, ok := smoothingInputFrames(fs, []string{"x", "missing"}); ok {
		t.Fatal("unknown frame must retain original behavior")
	}
	if !reflect.DeepEqual(psc.pc.movableFrames, []string{"tool", "x", "sensor"}) {
		t.Fatal("collision moving-frame list was mutated")
	}
	// A valid same-side trajectory still passes the mandatory final sweep.
	valid := []*referenceframe.LinearInputs{
		referenceframe.FrameSystemInputs{"x": {-10}, "tool": {}, "sensor": {}}.ToLinearInputs(),
		referenceframe.FrameSystemInputs{"x": {-5}, "tool": {}, "sensor": {}}.ToLinearInputs(),
	}
	if _, err := addCloseObstacleWaypoints(context.Background(), psc, valid); err != nil {
		t.Fatalf("valid attached geometry path rejected: %v", err)
	}
	// Opaque callbacks retain the old optional trial order and call behavior.
	calls := 0
	checker.SetCollisionConstraints(motionplan.CollisionConstraints{Obstacle: func(*motionplan.StateFS) (float64, error) {
		calls++
		return math.Inf(1), nil
	}})
	if checker.HasNativeMotionBounds() {
		t.Fatal("custom callbacks entered native scope")
	}
	smoothMultiArms(context.Background(), psc, steps)
	if calls == 0 {
		t.Fatal("custom optional trial behavior changed")
	}
}

type inputFrameSmoothingOpaqueFrame struct{ referenceframe.Frame }

func TestInputFrameSmoothingMultipleOwners(t *testing.T) {
	fs := referenceframe.NewEmptyFrameSystem("two-input-owners")
	for _, name := range []string{"x", "y"} {
		f, err := referenceframe.NewTranslationalFrame(name, r3.Vector{X: 1}, referenceframe.Limit{Min: -20, Max: 20})
		if err != nil {
			t.Fatal(err)
		}
		if err := fs.AddFrame(f, fs.World()); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{"tool", "sensor"} {
		f, err := referenceframe.NewStaticFrame(name, spatialmath.NewZeroPose())
		if err != nil {
			t.Fatal(err)
		}
		if err := fs.AddFrame(f, fs.Frame("x")); err != nil {
			t.Fatal(err)
		}
	}
	var steps []*referenceframe.LinearInputs
	for _, xy := range [][2]float64{{-10, 0}, {-8, 5}, {-5, 0}} {
		steps = append(steps, referenceframe.FrameSystemInputs{
			"x": {xy[0]}, "y": {xy[1]}, "tool": {}, "sensor": {},
		}.ToLinearInputs())
	}
	moving := []string{"tool", "y", "sensor", "x"}
	owners, ok := smoothingInputFrames(fs, moving)
	if !ok || !reflect.DeepEqual(owners, []string{"y", "x"}) {
		t.Fatalf("input-owner order changed: %v", owners)
	}
	logger := logging.NewTestLogger(t)
	cache := motionplan.NewCollisionCache()
	checker, err := motionplan.NewConstraintChecker(1e-8, nil, nil, nil, fs, nil, nil,
		map[string]bool{"x": true, "y": true, "tool": true, "sensor": true}, steps[0], nil, logger, cache)
	if err != nil || !checker.HasNativeMotionBounds() {
		t.Fatalf("native fixture unavailable: %v", err)
	}
	meta := &PlanMeta{}
	psc := &PlanSegmentContext{pc: &PlanContext{
		fs: fs, movableFrames: moving,
		planOpts: NewBasicPlannerOptions(), planMeta: meta, logger: logger,
	}, Checker: checker}
	smoothMultiArms(context.Background(), psc, steps)
	for i, step := range steps {
		if step.Get("y")[0] != 0 || step.Get("x")[0] != -10+float64(i)*2.5 {
			t.Fatalf("independent input trial lost at step %d: %v", i, step)
		}
	}
}

func TestInputFrameSmoothingOpaqueFrameScope(t *testing.T) {
	fs := referenceframe.NewEmptyFrameSystem("opaque")
	f, err := referenceframe.NewStaticFrame("opaque", spatialmath.NewZeroPose())
	if err != nil {
		t.Fatal(err)
	}
	if err := fs.AddFrame(&inputFrameSmoothingOpaqueFrame{f}, fs.World()); err != nil {
		t.Fatal(err)
	}
	start := referenceframe.FrameSystemInputs{"opaque": {}}.ToLinearInputs()
	c, err := motionplan.NewConstraintChecker(1e-8, nil, nil, nil, fs, nil, nil,
		map[string]bool{"opaque": true}, start, nil, logging.NewTestLogger(t), motionplan.NewCollisionCache())
	if err != nil {
		t.Fatal(err)
	}
	if c.HasNativeMotionBounds() {
		t.Fatal("opaque zero-DoF frame must keep ordinary optional trials")
	}
}
