package referenceframe

import (
	"testing"

	"github.com/golang/geo/r3"
	"go.viam.com/test"

	"go.viam.com/rdk/spatialmath"
)

func TestJointInterpolationStepsUnits(t *testing.T) {
	rotation, err := NewRotationalFrame("rotate", spatialmath.R4AA{RZ: 1}, Limit{Min: -100, Max: 100})
	test.That(t, err, test.ShouldBeNil)
	translation, err := NewTranslationalFrame("slide", r3.Vector{X: 1}, Limit{Min: -100000, Max: 100000})
	test.That(t, err, test.ShouldBeNil)
	model, err := NewSerialModel("mixed", []Frame{rotation, translation})
	test.That(t, err, test.ShouldBeNil)
	for _, tc := range []struct {
		frame    Frame
		from, to []Input
		want     int
	}{
		{rotation, []Input{0}, []Input{.01}, 5},
		{translation, []Input{0}, []Input{10}, 5},
		{model, []Input{0, 0}, []Input{.02, 10}, 10},
		{NewNamedFrame(model, "renamed"), []Input{0, 0}, []Input{.02, 10}, 10},
	} {
		n, ok := JointInterpolationSteps(tc.frame, tc.from, tc.to, 2, .002)
		test.That(t, ok, test.ShouldBeTrue)
		test.That(t, n, test.ShouldEqual, tc.want)
	}
}
