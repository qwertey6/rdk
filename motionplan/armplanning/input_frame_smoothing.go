package armplanning

import "go.viam.com/rdk/referenceframe"

// smoothingInputFrames separates independently interpolated inputs from geometry
// that moves with an ancestor. Only call this in complete native scope; opaque
// frame implementations retain the original optional smoothing behavior.
func smoothingInputFrames(fs *referenceframe.FrameSystem, moving []string) ([]string, bool) {
	if fs == nil {
		return nil, false
	}
	owners := make([]string, 0, len(moving))
	for _, name := range moving {
		frame := fs.Frame(name)
		if frame == nil {
			return nil, false
		}
		if len(frame.DoF()) > 0 {
			owners = append(owners, name)
		}
	}
	return owners, true
}
