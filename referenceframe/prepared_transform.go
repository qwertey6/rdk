package referenceframe

import (
	"fmt"

	"github.com/golang/geo/r3"
	"gonum.org/v1/gonum/num/dualquat"
	"gonum.org/v1/gonum/num/quat"

	"go.viam.com/rdk/spatialmath"
)

// PreparedTransform resolves a frame chain and its flat input offsets once.
// It is immutable and safe for concurrent evaluations. Construct it for each
// planning snapshot; it must not outlive changes to the model or frame system.
// Unsupported custom frames return an error during preparation so callers can
// retain the normal FrameSystem.TransformToDQ implementation.
type PreparedTransform struct {
	from, to  []preparedFrame
	dimension int
}

type preparedFrame struct {
	kind             byte // fixed, rotational, translational, or model
	fixed            spatialmath.DualQuaternion
	axis             r3.Vector
	input            int
	multiplier, bias float64
	limit            *Limit
	model            []preparedFrame
}

func preparedIdentity() dualquat.Number { return dualquat.Number{Real: quat.Number{Real: 1}} }

// PrepareTransform constructs an evaluator that consumes inputs in schema order.
func (fs *FrameSystem) PrepareTransform(schema *LinearInputsSchema, from, to string) (*PreparedTransform, error) {
	if schema == nil {
		return nil, fmt.Errorf("missing input schema")
	}
	p := &PreparedTransform{}
	for _, meta := range schema.metas {
		p.dimension += meta.dof
	}
	offsetFor := func(name string) (int, error) {
		for _, meta := range schema.metas {
			if meta.frameName == name && meta.dof > 0 {
				return meta.offset, nil
			}
		}
		for component, bundle := range fs.flattened {
			for _, outer := range schema.metas {
				if outer.frameName != component {
					continue
				}
				for _, inner := range bundle.schema.metas {
					if inner.frameName == name && inner.dof > 0 {
						return outer.offset + inner.offset, nil
					}
				}
			}
		}
		return 0, fmt.Errorf("no flat input offset for %q", name)
	}
	chain := func(name string) ([]preparedFrame, error) {
		frame := fs.lookupFrame(name)
		if frame == nil {
			return nil, NewFrameMissingError(name)
		}
		var result []preparedFrame
		for fs.parents[frame.Name()] != "" {
			offset, multiplier, bias := 0, 1., 0.
			var err error
			if mimic := fs.mimicFrames[frame.Name()]; mimic != nil {
				offset, err = offsetFor(mimic.sourceFrameName)
				multiplier, bias = mimic.multiplier, mimic.offset
			} else if len(frame.DoF()) > 0 {
				offset, err = offsetFor(frame.Name())
			}
			if err != nil {
				return nil, err
			}
			op, err := prepareFrame(frame, offset, multiplier, bias, true)
			if err != nil {
				return nil, err
			}
			result = append(result, op)
			frame = fs.lookupFrame(fs.parents[frame.Name()])
			if frame == nil {
				return nil, fmt.Errorf("missing parent in prepared transform")
			}
		}
		return result, nil
	}
	var err error
	p.from, err = chain(from)
	if err != nil {
		return nil, err
	}
	p.to, err = chain(to)
	return p, err
}

func prepareFrame(frame Frame, offset int, multiplier, bias float64, checkLimits bool) (preparedFrame, error) {
	for {
		named, ok := frame.(*namedFrame)
		if !ok {
			break
		}
		frame = named.Frame
	}
	op := preparedFrame{input: offset, multiplier: multiplier, bias: bias}
	switch f := frame.(type) {
	case *tailGeometryStaticFrame:
		return prepareFrame(f.staticFrame, offset, multiplier, bias, checkLimits)
	case *staticFrame:
		dq, ok := f.transform.(*spatialmath.DualQuaternion)
		if !ok {
			return op, fmt.Errorf("unsupported static pose %T", f.transform)
		}
		op.kind, op.fixed = 'f', *dq
	case *rotationalFrame:
		op.kind, op.axis = 'r', f.rotAxis
		if checkLimits {
			limit := f.DoF()[0]
			op.limit = &limit
		}
	case *translationalFrame:
		op.kind, op.axis = 't', f.transAxis
		if checkLimits {
			limit := f.DoF()[0]
			op.limit = &limit
		}
	case *SimpleModel:
		op.kind = 'm'
		for i, inner := range f.transformChain {
			index, scale, shift := f.transformChainInputOffsets[i], 1., 0.
			if index == -1 {
				mimic := f.mimicMappings[inner.Name()]
				if mimic == nil {
					return op, fmt.Errorf("missing mimic mapping")
				}
				index, scale, shift = mimic.sourceInputIdx, mimic.valueMultiplier, mimic.valueOffset
			}
			// SimpleModel.Transform handles revolute joints without a limits check.
			_, rotation := inner.(*rotationalFrame)
			child, err := prepareFrame(inner, offset+index, scale, shift, !rotation)
			if err != nil {
				return op, err
			}
			op.model = append(op.model, child)
		}
	default:
		return op, fmt.Errorf("unsupported prepared frame %T", frame)
	}
	return op, nil
}

func (op *preparedFrame) evaluate(inputs []float64) (spatialmath.DualQuaternion, error) {
	if op.kind == 'f' {
		return op.fixed, nil
	}
	if op.kind == 'm' {
		result := spatialmath.DualQuaternion{Number: preparedIdentity()}
		for i := range op.model {
			local, err := op.model[i].evaluate(inputs)
			if err != nil {
				return result, err
			}
			result.Number = result.Transformation(local.Number)
		}
		return result, nil
	}
	value := inputs[op.input]*op.multiplier + op.bias
	if op.limit != nil && !op.limit.IsValid(value) {
		return spatialmath.DualQuaternion{}, fmt.Errorf("prepared joint outside limits")
	}
	if op.kind == 'r' {
		aa := spatialmath.R4AA{Theta: value, RX: op.axis.X, RY: op.axis.Y, RZ: op.axis.Z}
		return spatialmath.DualQuaternion{Number: dualquat.Number{Real: aa.Quaternion()}}, nil
	}
	result := spatialmath.DualQuaternion{Number: preparedIdentity()}
	result.SetTranslation(op.axis.Mul(value))
	return result, nil
}

// Evaluate computes the relative pose without constructing LinearInputs or
// traversing frame-name maps. It preserves native dual-quaternion composition order.
func (p *PreparedTransform) Evaluate(inputs []float64) (spatialmath.DualQuaternion, error) {
	if len(inputs) != p.dimension {
		return spatialmath.DualQuaternion{}, NewIncorrectDoFError(len(inputs), p.dimension)
	}
	eval := func(chain []preparedFrame) (spatialmath.DualQuaternion, error) {
		result := spatialmath.DualQuaternion{Number: preparedIdentity()}
		for i := range chain {
			local, err := chain[i].evaluate(inputs)
			if err != nil {
				return result, err
			}
			result.Number = local.Transformation(result.Number)
		}
		return result, nil
	}
	from, err := eval(p.from)
	if err != nil {
		return from, err
	}
	if len(p.to) > 0 {
		to, err := eval(p.to)
		if err != nil {
			return to, err
		}
		inverse := spatialmath.DualQuaternion{Number: dualquat.ConjQuat(to.Number)}
		from.Number = inverse.Transformation(from.Number)
	}
	from.Number = from.Transformation(preparedIdentity())
	return from, nil
}
