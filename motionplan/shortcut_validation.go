package motionplan

import "context"

// CheckShortcutFS checks the same samples as CheckStateConstraintsAcrossSegmentFS,
// but can reject an invalid shortcut before checking its complete valid prefix.
// It is for callers needing only a verdict. Custom callbacks/interpolation retain
// their original evaluation order and ownership. Accepted edges still visit every
// required sample, with the existing conservative clearance certificates.
func (c *ConstraintChecker) CheckShortcutFS(ctx context.Context, ci *SegmentFS, resolution float64, checkFinal bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	validation, err := c.SegmentValidation(ci, resolution)
	if err != nil {
		return err
	}
	samples := validation.Samples
	// A native interpolation can still traverse custom zero-DoF frames. The
	// prepared motion-bound tree recognizes every frame and geometry, including
	// those without joint inputs. Its presence establishes the native-only scope;
	// the actual edge need not have a numerically stable clearance certificate.
	if !c.nativeCallbacks || c.motionBounds == nil || samples.prepared == nil || samples.steps < 32 {
		_, err := checkSegmentInOrder(ctx, ci, validation, checkFinal)
		return err
	}

	check := func(i int) error {
		_, _, err := validation.Check(ctx, i, 0)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}
	if err := check(0); err != nil {
		return err
	}
	// Breadth-first subdivision of the original grid: midpoint, quarters,
	// eighths, sixteenths. These are original samples, never a coarser acceptance
	// test. Integer rounding and endpoint inclusion match the original grid.
	for stride := 8; stride > 0; stride /= 2 {
		for k := stride; k < 16; k += 2 * stride {
			if err := check(samples.steps * k / 16); err != nil {
				return err
			}
		}
	}
	end := samples.Len()
	if !checkFinal {
		end--
	}
	// Skip the already checked samples exactly once, in ascending grid order.
	nextK := 1
	for i := 1; i < end; i++ {
		if nextK < 16 && i == samples.steps*nextK/16 {
			nextK++
			continue
		}
		if err := check(i); err != nil {
			return err
		}
	}
	return nil
}
