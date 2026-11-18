package spatialmath

import "sync"

// GeometryTransformBuffer owns a reusable transformed geometry. The source must
// keep its shape and label fixed, but may itself be a buffer whose pose changes.
// Returned geometry is borrowed until the next Transform. Neither the buffer nor
// its borrowed output may be used concurrently. Unknown geometry types retain
// their ordinary Transform implementation.
type GeometryTransformBuffer struct {
	source Geometry
	output Geometry
	pose   DualQuaternion
	child  *GeometryTransformBuffer
}

// NewGeometryTransformBuffer prepares private storage without changing source.
func NewGeometryTransformBuffer(source Geometry) *GeometryTransformBuffer {
	b := &GeometryTransformBuffer{source: source}
	switch g := source.(type) {
	case *box:
		b.output = &box{center: &b.pose, halfSize: g.halfSize, boundingSphereR: g.boundingSphereR, label: g.label}
	case *sphere:
		b.output = &sphere{pose: &b.pose, radius: g.radius, label: g.label}
	case *capsule:
		b.output = &capsule{pose: &b.pose, radius: g.radius, length: g.length, label: g.label}
	case *Mesh:
		out := &Mesh{
			pose: &b.pose, triangles: g.triangles, label: g.label, fileType: g.fileType,
			rawBytes: g.rawBytes, originalFilePath: g.originalFilePath, state: g.state,
		}
		out.bvh.Store(g.ensureBVH())
		vertices := g.ensureUniqueVertices()
		out.uniqueVerts.Store(&vertices)
		b.output = out
	case *Cylinder:
		b.child = NewGeometryTransformBuffer(g.mesh)
		b.output = &Cylinder{pose: &b.pose, radius: g.radius, height: g.height, capped: g.capped, label: g.label}
	case *point:
		b.output = &point{label: g.label}
	}
	return b
}

// Transform uses the same composition and collision implementations as Geometry.Transform.
func (b *GeometryTransformBuffer) Transform(parent Pose) Geometry {
	if b.output == nil {
		return b.source.Transform(parent)
	}
	if out, ok := b.output.(*point); ok {
		out.position = TransformPointByPose(parent, b.source.(*point).position)
		return out
	}
	b.pose.Number = DualQuaternionFromPose(parent).Transformation(DualQuaternionFromPose(b.source.Pose()).Number)
	switch out := b.output.(type) {
	case *box:
		out.centerPt = b.pose.Point()
		out.mesh, out.rotMatrix = nil, nil
		out.meshOnce = sync.Once{}
		out.once, out.aabbOnce = sync.Once{}, sync.Once{}
	case *capsule:
		src := b.source.(*capsule)
		out.segA = TransformPointByPose(parent, src.segA)
		out.segB = TransformPointByPose(parent, src.segB)
		out.center = b.pose.Point()
		out.capVec = out.segB.Sub(out.center)
		out.rotMatrix, out.once = nil, sync.Once{}
	case *Mesh:
		out.poseCacheOnce = sync.Once{}
	case *Cylinder:
		out.mesh = b.child.Transform(parent).(*Mesh)
	}
	return b.output
}
