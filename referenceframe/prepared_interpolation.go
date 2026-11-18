package referenceframe

// PreparedInterpolation interpolates native scalar joint coordinates into private
// reusable storage. Its result is borrowed until the next At call. It is not safe
// for concurrent use. Unknown frame implementations retain Frame.Interpolate.
type PreparedInterpolation struct {
	from, to []Input
	result   *LinearInputs
}

// PrepareInterpolation recognizes only frames whose interpolation is exactly
// from+(to-from)*fraction. Failure to recognize a frame is not a planning error:
// callers must use its original interpolation implementation instead.
func PrepareInterpolation(fs *FrameSystem, from, to *LinearInputs) (*PreparedInterpolation, bool) {
	p := &PreparedInterpolation{result: NewLinearInputs()}
	for name, start := range from.Items() {
		if len(start) == 0 {
			continue
		}
		end, frame := to.Get(name), fs.Frame(name)
		if frame == nil || len(start) != len(end) || len(start) != len(frame.DoF()) || !scalarInterpolation(frame) {
			return nil, false
		}
		p.from = append(p.from, start...)
		p.to = append(p.to, end...)
		p.result.Put(name, start)
	}
	return p, true
}

func scalarInterpolation(frame Frame) bool {
	switch f := frame.(type) {
	case *namedFrame:
		return scalarInterpolation(f.Frame)
	case *staticFrame, *tailGeometryStaticFrame, *rotationalFrame, *translationalFrame:
		return true
	case *SimpleModel:
		for _, child := range f.framesInOrder() {
			if !scalarInterpolation(child) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// At returns borrowed inputs at fraction. Fractions and arithmetic match
// Frame.Interpolate, including its floating point treatment of both endpoints.
func (p *PreparedInterpolation) At(fraction float64) *LinearInputs {
	for i, a := range p.from {
		p.result.inputs[i] = a + ((p.to[i] - a) * fraction)
	}
	return p.result
}
