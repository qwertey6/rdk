package motionplan

import (
	"context"
	"encoding/binary"
	"math"
	"math/bits"
	"sync"

	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func exactStateNativeGeometry(sets ...[]spatialmath.Geometry) bool {
	for _, set := range sets {
		for _, g := range set {
			if _, _, ok := spatialmath.BoundingSphere(g); !ok {
				return false
			}
			if _, ok := g.Pose().(*spatialmath.DualQuaternion); !ok {
				return false
			}
		}
	}
	return true
}

const (
	exactStateShards          = 16
	exactStateEntriesPerShard = 512
	exactStateBytesPerShard   = 256 << 10
	exactStateMaxKeyBytes     = 4096
)

// exactStateMemo retains successful collision clearances within one planning
// request. Each immutable native checker has a distinct scope. Keys own every
// frame name, width and floating-point bit, including zero-DoF entries; schema
// pointers and hashes alone are never evidence of equality. FIFO eviction bounds
// retained key payload to 4 MiB and entries to 8192, excluding Go map overhead.
// Topological constraints always execute before a lookup.
type exactStateMemo struct {
	shards [exactStateShards]exactStateShard
}

type exactStateShard struct {
	mu          sync.Mutex
	values      map[string]float64
	fifo        []string
	head, bytes int
}

func exactStateKey(scope uint64, q *referenceframe.LinearInputs, buffer []byte) ([]byte, uint64, bool) {
	if scope == 0 || q == nil || len(buffer) < 8 {
		return nil, 0, false
	}
	binary.LittleEndian.PutUint64(buffer, scope)
	n, hash := 8, scope
	for name, values := range q.Items() {
		if len(name) > len(buffer)-n-8 || len(values) > (len(buffer)-n-8-len(name))/8 {
			return nil, 0, false
		}
		binary.LittleEndian.PutUint32(buffer[n:], uint32(len(name)))
		binary.LittleEndian.PutUint32(buffer[n+4:], uint32(len(values)))
		n += 8
		copy(buffer[n:], name)
		n += len(name)
		for _, value := range values {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return nil, 0, false
			}
			v := math.Float64bits(value)
			binary.LittleEndian.PutUint64(buffer[n:], v)
			n += 8
			hash = bits.RotateLeft64(hash, 13) ^ v
		}
	}
	// Hash chooses a shard only. The map compares the complete owned key.
	hash ^= hash >> 32
	hash *= 0x9e3779b97f4a7c15
	hash ^= hash >> 32
	return buffer[:n], hash, true
}

func (m *exactStateMemo) lookup(key []byte, hash uint64) (float64, bool) {
	s := &m.shards[hash%exactStateShards]
	s.mu.Lock()
	v, ok := s.values[string(key)] // Lookup conversion does not retain borrowed bytes.
	s.mu.Unlock()
	return v, ok
}

func (m *exactStateMemo) store(key []byte, hash uint64, clearance float64) {
	s := &m.shards[hash%exactStateShards]
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.values[string(key)]; ok {
		return
	}
	if s.values == nil {
		s.values = make(map[string]float64)
		s.fifo = make([]string, exactStateEntriesPerShard)
	}
	for len(s.values) >= exactStateEntriesPerShard || s.bytes+len(key) > exactStateBytesPerShard {
		old := s.fifo[s.head]
		delete(s.values, old)
		s.bytes -= len(old)
		s.fifo[s.head] = ""
		s.head = (s.head + 1) % len(s.fifo)
	}
	owned := string(key)
	s.fifo[(s.head+len(s.values))%len(s.fifo)] = owned
	s.values[owned] = clearance
	s.bytes += len(owned)
}

// WithFinishingCollisionReuse permits exact native collision-state reuse in this
// finishing call tree. Search contexts remain untagged. Native eligibility,
// immutable checker scope, topology and the exact input key are still required.
func WithFinishingCollisionReuse(ctx context.Context) context.Context {
	return context.WithValue(ctx, finishingCollisionReuseKey{}, true)
}

type finishingCollisionReuseKey struct{}

func (c *ConstraintChecker) checkExactStateCollisions(ctx context.Context, state *StateFS) (float64, error) {
	var keyBuffer [exactStateMaxKeyBytes]byte
	memoKey, memoHash, eligible := exactStateKey(c.exactStateScope, state.Configuration, keyBuffer[:])
	if eligible {
		if distance, ok := c.cache.exactStates.lookup(memoKey, memoHash); ok {
			return distance, nil
		}
	}
	closest, err := c.checkStateCollisions(state)
	if err == nil && eligible && ctx.Err() == nil && closest >= 0 && !math.IsNaN(closest) {
		c.cache.exactStates.store(memoKey, memoHash, closest)
	}
	return closest, err
}
