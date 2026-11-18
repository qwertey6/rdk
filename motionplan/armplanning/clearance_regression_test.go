package armplanning

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
)

// This returned trajectory failed the independent audit of the extracted FK PR.
// Its final chord crosses a blocker between the legacy validation samples.
func TestCheckPathRejectsRecordedClearanceShortcut(t *testing.T) {
	ctx := context.Background()
	req, err := ReadRequestFromFile("data/native-clearance-regression.json")
	test.That(t, err, test.ShouldBeNil)
	var fixture struct {
		Trajectory []referenceframe.FrameSystemInputs `json:"trajectory"`
	}
	data, err := os.ReadFile("data/native-clearance-regression-path.json")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, json.Unmarshal(data, &fixture), test.ShouldBeNil)
	path := make([]*referenceframe.LinearInputs, len(fixture.Trajectory))
	for i, q := range fixture.Trajectory {
		path[i] = q.ToLinearInputs()
	}
	pc, err := NewPlanContext(ctx, logging.NewTestLogger(t), req, &PlanMeta{})
	test.That(t, err, test.ShouldBeNil)
	goals, err := req.Goals[0].ComputePoses(ctx, req.FrameSystem)
	test.That(t, err, test.ShouldBeNil)
	psc, err := NewPlanSegmentContext(ctx, pc, path[0], goals)
	test.That(t, err, test.ShouldBeNil)
	test.That(t, validateDenseEdges(ctx, psc, path), test.ShouldNotBeNil)
	test.That(t, psc.CheckPath(ctx, path[len(path)-2], path[len(path)-1], true, nil), test.ShouldNotBeNil)
}
