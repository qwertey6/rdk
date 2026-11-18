package motionplan

import (
	"context"
	"math"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func TestSharedGeometryFKCollisionParity(t *testing.T) {
	c, edge := adaptiveArms(t)
	for i := 0; i < 200; i++ {
		q := edge.StartConfiguration.Copy()
		q.Get("left")[0] = .6 + float64(i)/250
		state := &StateFS{Configuration: q, FS: edge.FS}
		got, gotErr := c.CheckStateFSConstraints(context.Background(), state)
		test.That(t, state.geometryEvaluation, test.ShouldBeNil)
		test.That(t, state.shareGeometryFK, test.ShouldBeFalse)
		c.nativeCallbacks = false // Original separate memos, same collision constraints.
		want, wantErr := c.CheckStateFSConstraints(context.Background(), &StateFS{Configuration: q, FS: edge.FS})
		c.nativeCallbacks = true
		test.That(t, gotErr == nil, test.ShouldEqual, wantErr == nil)
		if gotErr != nil && wantErr != nil {
			test.That(t, gotErr.Error(), test.ShouldEqual, wantErr.Error())
		}
		test.That(t, math.Float64bits(got), test.ShouldEqual, math.Float64bits(want))
	}
}

type customFKFrame struct{ referenceframe.Frame }

func TestSharedGeometryFKCustomFallback(t *testing.T) {
	_, edge := adaptiveArms(t)
	fs := edge.FS
	// Even a zero-DoF custom frame outside the moving arm disables sharing.
	test.That(t, fs.AddFrame(customFKFrame{referenceframe.NewZeroStaticFrame("custom")}, fs.World()), test.ShouldBeNil)
	geometries, err := referenceframe.FrameSystemGeometriesLinearInputs(fs, edge.StartConfiguration)
	test.That(t, err, test.ShouldBeNil)
	var robot []spatialmath.Geometry
	moving := map[string]bool{}
	for name, gif := range geometries {
		moving[name] = true
		robot = append(robot, gif.Geometries()...)
	}
	c, err := NewConstraintChecker(1e-8, nil, nil, nil, fs, robot, nil, moving, edge.StartConfiguration, nil,
		logging.NewTestLogger(t), NewCollisionCache())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, c.motionBounds, test.ShouldBeNil)
	_, err = c.CheckStateFSConstraints(context.Background(), &StateFS{Configuration: edge.StartConfiguration, FS: fs})
	test.That(t, err, test.ShouldBeNil)

	c, edge = adaptiveArms(t)
	called := 0
	c.SetCollisionConstraints(CollisionConstraints{Obstacle: func(state *StateFS) (float64, error) {
		called++
		test.That(t, state.shareGeometryFK, test.ShouldBeFalse)
		return 100, nil
	}})
	_, err = c.CheckStateFSConstraints(context.Background(), &StateFS{Configuration: edge.StartConfiguration, FS: edge.FS})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, called, test.ShouldEqual, 1)
}
