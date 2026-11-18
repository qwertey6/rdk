package armplanning

import (
	"math"

	"go.viam.com/rdk/motionplan"
	"go.viam.com/rdk/referenceframe"
)

// shortcutQuality accepts a speculative finisher only when both the planner
// metric and whole-configuration Euclidean length are within 10% of their direct
// lower bounds. It never accepts an edge or bypasses mandatory final validation.
// The bound applies to these endpoints; IK choices and later multi-arm cleanup
// still require separate final path-quality measurement.
type shortcutQuality struct {
	psc                      *PlanSegmentContext
	length, metric           float64
	lengthLimit, metricLimit float64
	enabled                  bool
}

func jointCoordinateDistance(a, b *referenceframe.LinearInputs) float64 {
	if len(a.GetLinearizedInputs()) != len(b.GetLinearizedInputs()) {
		return math.Inf(1)
	}
	total := 0.
	for name, x := range a.Items() {
		y := b.Get(name)
		if len(x) != len(y) {
			return math.Inf(1)
		}
		for i, v := range x {
			delta := v - y[i]
			total += delta * delta
		}
	}
	return math.Sqrt(total)
}

func newShortcutQuality(psc *PlanSegmentContext, steps []*referenceframe.LinearInputs) *shortcutQuality {
	q := &shortcutQuality{psc: psc}
	if len(steps) < 2 || psc.pc.planOpts.ConfigurationDistanceMetric != motionplan.FSConfigurationL2DistanceMetric {
		return q
	}
	a, b := steps[0], steps[len(steps)-1]
	q.lengthLimit = 1.10 * jointCoordinateDistance(a, b)
	q.metricLimit = 1.10 * q.metricDistance(a, b)
	q.enabled = q.lengthLimit > 0 && q.metricLimit > 0 && !math.IsInf(q.lengthLimit, 0) && !math.IsInf(q.metricLimit, 0)
	q.refresh(steps)
	return q
}

func (q *shortcutQuality) metricDistance(a, b *referenceframe.LinearInputs) float64 {
	return q.psc.pc.ConfigurationDistanceFunc(&motionplan.SegmentFS{StartConfiguration: a, EndConfiguration: b})
}

func (q *shortcutQuality) refresh(steps []*referenceframe.LinearInputs) {
	q.length, q.metric = 0, 0
	for i := 1; i < len(steps); i++ {
		q.length += jointCoordinateDistance(steps[i-1], steps[i])
		q.metric += q.metricDistance(steps[i-1], steps[i])
	}
}

func (q *shortcutQuality) satisfied(steps []*referenceframe.LinearInputs) bool {
	if !q.enabled {
		return false
	}
	q.refresh(steps)
	return q.length <= q.lengthLimit && q.metric <= q.metricLimit
}
