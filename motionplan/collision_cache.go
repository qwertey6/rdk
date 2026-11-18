package motionplan

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/singleflight"

	"go.viam.com/rdk/spatialmath"
)

// CollisionCache holds planner-level temporal-coherence state for collision
// queries within a single planning request. Owned by planContext and threaded
// down through the constraint checker. Holds collision hints:
//
//  1. Geometry-pair "last-violated" hints — one slot per constraint type. When
//     a constraint check finds (geomA, geomB) in collision, it stores that pair
//     so the next call to the same constraint tries that pair first.
//
// Native path validity is memoized separately by PlanSegmentContext with the
// checker, direction, resolution and endpoint policy in its key. The legacy
// endpoint-only methods below must not be used for that purpose.
//
// Per-mesh witness caches (the inner-loop temporal-coherence short-circuit)
// live on spatialmath.Mesh.state, not here. Threading a cache through the
// motionplan call chain was measurably slower than direct field access on
// the mesh.
type CollisionCache struct {
	exactStates      exactStateMemo
	exactStateScopes atomic.Uint64
	// obstaclePairHint, selfPairHint, robotPairHint each cache the
	// most-recently-violated geometry-label pair for one constraint type.
	// Lock-free reads via atomic; stale reads are harmless.
	obstaclePairHint atomic.Pointer[[2]string]
	selfPairHint     atomic.Pointer[[2]string]
	robotPairHint    atomic.Pointer[[2]string]

	// edgeResults memoizes the outcome of CheckStateConstraintsAcrossSegmentFS
	// for an interpolated edge. Key is the canonical {hashA, hashB} pair —
	// uint64 fits inside sync.Map's interface{} slot without allocation.
	edgeResults sync.Map // edgeResultKey -> edgeResultValue
}

// NewCollisionCache constructs an empty cache. Safe for concurrent use.
func NewCollisionCache() *CollisionCache {
	return &CollisionCache{}
}

// edgeResultKey identifies an interpolated edge by hashed-config endpoints.
// Symmetric (edges are bidirectional) — sorting the hashes canonicalizes the key.
type edgeResultKey struct {
	a, b uint64
}

// edgeResultValue is retained for the legacy endpoint-only cache API.
type edgeResultValue struct {
	isClear bool
}

// LookupEdgeResult returns whether the edge between two configurations has been
// previously verified collision-free. Returns (false, false) for "no cached result".
// Caller hashes must be deterministic for the same inputs across calls.
//
// Deprecated: this API omits the checker, direction, resolution and endpoint policy.
// Native planning uses a segment-scoped cache instead.
func (c *CollisionCache) LookupEdgeResult(hashA, hashB uint64) (isClear, ok bool) {
	if c == nil {
		return false, false
	}
	if hashA > hashB {
		hashA, hashB = hashB, hashA
	}
	v, ok := c.edgeResults.Load(edgeResultKey{a: hashA, b: hashB})
	if !ok {
		return false, false
	}
	return v.(edgeResultValue).isClear, true
}

// StoreEdgeResult records that an edge was found collision-free.
//
// Deprecated: use a cache scoped to the exact validation context instead.
func (c *CollisionCache) StoreEdgeResult(hashA, hashB uint64, isClear bool) {
	if c == nil {
		return
	}
	if hashA > hashB {
		hashA, hashB = hashB, hashA
	}
	c.edgeResults.Store(edgeResultKey{a: hashA, b: hashB}, edgeResultValue{isClear: isClear})
}

// sdfRegistry caches voxel distance fields process-wide. A field depends only
// on the static scene (keyed by the shape-aware spatialmath.GeometrySetHash),
// never on the plan, and costs ~35-70ms plus a multi-MB grid allocation to
// build. Bound retained storage by bytes as well as entry count, so a sequence
// of small recurring scenes need not rebuild fields just because it uses more
// than eight keys. Changed poses/shapes still invalidate the key.
type sdfRegistry struct {
	mu              sync.Mutex
	tick            uint64
	entries         map[uint64]*sdfRegistryEntry
	bytes, maxBytes int64
	builds          singleflight.Group
}

type sdfRegistryEntry struct {
	sdf      *spatialmath.VoxelSDF // nil records a scene the SDF cannot represent
	lastUsed uint64
	bytes    int64
}

const (
	sdfRegistryCap   = 128
	sdfRegistryBytes = 512 << 20
)

var globalSDFRegistry = sdfRegistry{entries: map[uint64]*sdfRegistryEntry{}}

func (r *sdfRegistry) lookup(key uint64) (*spatialmath.VoxelSDF, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[key]
	if !ok {
		return nil, false
	}
	r.tick++
	e.lastUsed = r.tick
	return e.sdf, true
}

func (r *sdfRegistry) insert(key uint64, sdf *spatialmath.VoxelSDF) *spatialmath.VoxelSDF {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.entries[key]; ok {
		// Another plan built the same field concurrently; keep the incumbent.
		return e.sdf
	}
	r.tick++
	bytes := sdf.MemoryBytes() + 128
	budget := r.maxBytes
	if budget == 0 {
		budget = sdfRegistryBytes
	}
	if bytes > budget {
		return sdf
	}
	r.entries[key] = &sdfRegistryEntry{sdf: sdf, lastUsed: r.tick, bytes: bytes}
	r.bytes += bytes
	for len(r.entries) > sdfRegistryCap || r.bytes > budget {
		var oldestKey uint64
		oldest := uint64(math.MaxUint64)
		for k, e := range r.entries {
			if e.lastUsed < oldest {
				oldest, oldestKey = e.lastUsed, k
			}
		}
		r.bytes -= r.entries[oldestKey].bytes
		delete(r.entries, oldestKey)
	}
	return sdf
}

// SDFFor returns the voxel distance field for the given static geometry set,
// building it on first request anywhere in the process and sharing it across
// plans of the same scene. sdfResolutionMM balances build time (~70ms for a
// dual-arm scene at 10mm) against the conservative margin subtracted from
// every query (half the voxel diagonal, ~8.7mm at 10mm).
func (c *CollisionCache) SDFFor(static []spatialmath.Geometry) *spatialmath.VoxelSDF {
	if c == nil || len(static) == 0 {
		return nil
	}
	key := spatialmath.GeometrySetHash(static)
	if sdf, ok := globalSDFRegistry.lookup(key); ok {
		return sdf
	}
	//nolint:errcheck // The build closure below always returns a nil error.
	result, _, _ := globalSDFRegistry.builds.Do(fmt.Sprintf("%x", key), func() (interface{}, error) {
		// Recheck after joining the flight: another caller may have finished
		// between the initial miss and acquiring the build slot.
		if sdf, ok := globalSDFRegistry.lookup(key); ok {
			return sdf, nil
		}
		built := spatialmath.NewVoxelSDF(static, sdfResolutionMM)
		return globalSDFRegistry.insert(key, built), nil
	})
	return result.(*spatialmath.VoxelSDF)
}

const sdfResolutionMM = 10.0
