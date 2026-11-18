package referenceframe

import "math"

// JointInterpolationSteps bounds interpolation increments in each joint's own
// units. It covers built-in rotational/translational frames and models, including
// joints outside the primary output chain. Unknown frames return false so their
// caller can retain its general frame interpolation policy. This is a sampling
// bound, not a continuous collision certificate or a bound on link displacement.
func JointInterpolationSteps(frame Frame, from, to []Input, linearStep, angularStep float64) (int, bool) {
	if len(from) != len(to) || len(from) != len(frame.DoF()) || linearStep <= 0 || angularStep <= 0 {
		return 0, false
	}
	if len(from) == 0 {
		return 0, true
	}
	stepsFor := func(delta, step float64) int { return int(math.Ceil(math.Abs(delta) / step)) }
	switch f := frame.(type) {
	case *namedFrame:
		return JointInterpolationSteps(f.Frame, from, to, linearStep, angularStep)
	case *rotationalFrame:
		return stepsFor(to[0]-from[0], angularStep), true
	case *translationalFrame:
		return stepsFor(to[0]-from[0], linearStep), true
	case *SimpleModel:
		steps := 0
		for _, meta := range f.inputSchema.metas {
			if meta.dof == 0 {
				continue
			}
			n, ok := JointInterpolationSteps(f.internalFS.Frame(meta.frameName),
				from[meta.offset:meta.offset+meta.dof], to[meta.offset:meta.offset+meta.dof], linearStep, angularStep)
			if !ok {
				return 0, false
			}
			steps = max(steps, n)
		}
		for name, mimic := range f.mimicMappings {
			a := from[mimic.sourceInputIdx]*mimic.valueMultiplier + mimic.valueOffset
			b := to[mimic.sourceInputIdx]*mimic.valueMultiplier + mimic.valueOffset
			n, ok := JointInterpolationSteps(f.internalFS.Frame(name), []Input{a}, []Input{b}, linearStep, angularStep)
			if !ok {
				return 0, false
			}
			steps = max(steps, n)
		}
		return steps, true
	default:
		return 0, false
	}
}
