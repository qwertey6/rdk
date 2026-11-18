package motionplan

import (
	"context"
	"encoding/binary"
	"errors"
	"math"
	"strings"
	"sync"
	"testing"

	"go.viam.com/test"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/spatialmath"
)

func TestExactStateKeyOwnershipAndIdentity(t *testing.T) {
	q := referenceframe.NewLinearInputs()
	q.Put("a", []float64{0, .7})
	q.Put("fixed", []float64{.2})
	var buf [exactStateMaxKeyBytes]byte
	k, h, ok := exactStateKey(1, q, buf[:])
	test.That(t, ok, test.ShouldBeTrue)
	m := new(exactStateMemo)
	m.store(k, h, 12.5)
	owned := string(k)
	clear(buf[:])
	value, hit := m.lookup([]byte(owned), h)
	test.That(t, hit, test.ShouldBeTrue)
	test.That(t, value, test.ShouldEqual, 12.5)
	for _, mutation := range []func() *referenceframe.LinearInputs{
		func() *referenceframe.LinearInputs { x := q.Copy(); x.Get("a")[0] = math.Copysign(0, -1); return x },
		func() *referenceframe.LinearInputs { x := q.Copy(); x.Get("a")[1] = math.Nextafter(.7, 1); return x },
		func() *referenceframe.LinearInputs {
			x := referenceframe.NewLinearInputs()
			x.Put("a", []float64{0, .7})
			return x
		},
		func() *referenceframe.LinearInputs {
			x := referenceframe.NewLinearInputs()
			x.Put("fixed", []float64{.2})
			x.Put("a", []float64{0, .7})
			return x
		},
		func() *referenceframe.LinearInputs {
			x := referenceframe.NewLinearInputs()
			x.Put("a", []float64{0})
			x.Put("fixed", []float64{.7})
			return x
		},
	} {
		key, hash, ok := exactStateKey(1, mutation(), buf[:])
		test.That(t, ok, test.ShouldBeTrue)
		test.That(t, string(key) == owned, test.ShouldBeFalse)
		_, hit := m.lookup(key, hash)
		test.That(t, hit, test.ShouldBeFalse)
	}
	k, _, ok = exactStateKey(2, q, buf[:])
	test.That(t, ok, test.ShouldBeTrue)
	test.That(t, string(k) == owned, test.ShouldBeFalse)
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		q.Get("a")[0] = bad
		_, _, ok := exactStateKey(1, q, buf[:])
		test.That(t, ok, test.ShouldBeFalse)
	}
	q = referenceframe.NewLinearInputs()
	q.Put(strings.Repeat("x", exactStateMaxKeyBytes), []float64{0})
	_, _, ok = exactStateKey(1, q, buf[:])
	test.That(t, ok, test.ShouldBeFalse)
}

func TestExactStateMemoBoundedAndConcurrent(t *testing.T) {
	for _, size := range []int{8, exactStateMaxKeyBytes} {
		m := new(exactStateMemo)
		var wg sync.WaitGroup
		for worker := range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				key := make([]byte, size)
				for i := range 1500 {
					id := uint64(worker*1500 + i)
					binary.LittleEndian.PutUint64(key, id)
					// Deliberately collide every key in the shard selector.
					m.store(key, 0, float64(id))
					if got, hit := m.lookup(key, 0); hit && got != float64(id) {
						t.Errorf("hash collision reused another state's clearance")
					}
				}
			}()
		}
		wg.Wait()
		s := &m.shards[0]
		test.That(t, len(s.values) <= exactStateEntriesPerShard, test.ShouldBeTrue)
		test.That(t, s.bytes <= exactStateBytesPerShard, test.ShouldBeTrue)
		bytes, entries := 0, 0
		for _, key := range s.fifo {
			if key == "" {
				continue
			}
			_, present := s.values[key]
			test.That(t, present, test.ShouldBeTrue)
			bytes += len(key)
			entries++
		}
		test.That(t, bytes, test.ShouldEqual, s.bytes)
		test.That(t, entries, test.ShouldEqual, len(s.values))
	}
}

func TestExactStateNativeChecksAndInvalidation(t *testing.T) {
	c, edge := adaptiveArms(t)
	test.That(t, c.exactStateScope != 0, test.ShouldBeTrue)
	for range 2 {
		_, err := c.CheckStateFSConstraints(context.Background(), &StateFS{FS: edge.FS, Configuration: edge.StartConfiguration})
		test.That(t, err, test.ShouldBeNil)
	}
	test.That(t, memoEntries(&c.cache.exactStates), test.ShouldEqual, 0)
	state := func(q *referenceframe.LinearInputs) *StateFS { return &StateFS{FS: edge.FS, Configuration: q} }
	collisionCalls := 0
	original := c.collisionConstraints.SelfCollision
	c.collisionConstraints.SelfCollision = func(s *StateFS) (float64, error) {
		collisionCalls++
		return original(s)
	}
	want, err := c.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), state(edge.StartConfiguration))
	test.That(t, err, test.ShouldBeNil)
	for range 2 {
		got, err := c.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), state(edge.StartConfiguration.Copy()))
		test.That(t, err, test.ShouldBeNil)
		test.That(t, math.Float64bits(got), test.ShouldEqual, math.Float64bits(want))
	}
	test.That(t, memoEntries(&c.cache.exactStates), test.ShouldEqual, 1)
	test.That(t, collisionCalls, test.ShouldEqual, 1)
	// Topology still executes on an already-cached clear state.
	topologyError := errors.New("topology changed")
	c.topoConstraint = func(*StateFS) error { return topologyError }
	_, err = c.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), state(edge.StartConfiguration))
	test.That(t, err, test.ShouldEqual, topologyError)
	c.topoConstraint = nil
	// Both arms moving to the same pose must still collide, not reuse an endpoint.
	q := edge.StartConfiguration.Copy()
	q.Get("left")[0], q.Get("right")[0] = 1, 1
	for range 2 {
		_, err = c.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), state(q))
		test.That(t, err, test.ShouldNotBeNil)
	}
	test.That(t, memoEntries(&c.cache.exactStates), test.ShouldEqual, 1)
	ctx, cancel := context.WithCancel(WithFinishingCollisionReuse(context.Background()))
	cancel()
	_, err = c.CheckStateConstraintsAcrossSegmentFS(ctx, edge, 1, true)
	test.That(t, err, test.ShouldEqual, context.Canceled)
	// Independent validation must use original materialization again.
	disableGeometryReuse(c)
	test.That(t, c.exactStateScope, test.ShouldEqual, 0)
	_, err = c.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), state(edge.StartConfiguration))
	test.That(t, err, test.ShouldBeNil)
	calls := 0
	c.SetCollisionConstraints(CollisionConstraints{Obstacle: func(*StateFS) (float64, error) {
		calls++
		return 2, nil
	}})
	for range 2 {
		_, err = c.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), state(edge.StartConfiguration))
		test.That(t, err, test.ShouldBeNil)
	}
	test.That(t, calls, test.ShouldEqual, 2)
}

type exactStateCustomPose struct{ spatialmath.Pose }

func TestExactStateCustomPoseFallback(t *testing.T) {
	fs := referenceframe.NewEmptyFrameSystem("custom-pose")
	f, err := referenceframe.NewStaticFrame("fixed", exactStateCustomPose{spatialmath.NewZeroPose()})
	test.That(t, err, test.ShouldBeNil)
	test.That(t, fs.AddFrame(f, fs.World()), test.ShouldBeNil)
	c, err := NewConstraintChecker(1e-8, nil, nil, nil, fs, nil, nil, map[string]bool{"fixed": true},
		referenceframe.NewLinearInputs(), nil, logging.NewTestLogger(t), NewCollisionCache())
	test.That(t, err, test.ShouldBeNil)
	test.That(t, c.motionBounds, test.ShouldNotBeNil)
	test.That(t, c.exactStateScope, test.ShouldEqual, 0)
}

func TestExactStateDenseParity(t *testing.T) {
	c, edge := adaptiveArms(t)
	audit, _ := adaptiveArms(t)
	disableGeometryReuse(audit)
	for pass := range 2 {
		for i := range 201 {
			q := edge.StartConfiguration.Copy()
			q.Get("left")[0], q.Get("right")[0] = .7+float64(i)*.003, 1.3-float64(i)*.003
			_, got := c.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), &StateFS{FS: edge.FS, Configuration: q})
			_, want := audit.CheckStateFSConstraints(WithFinishingCollisionReuse(context.Background()), &StateFS{FS: edge.FS, Configuration: q})
			if (got == nil) != (want == nil) {
				t.Fatalf("pass %d sample %d collision mismatch", pass, i)
			}
		}
	}
}

// Test-only dense oracle retains the original geometry materialization.
func disableGeometryReuse(c *ConstraintChecker) {
	c.exactStateScope = 0
	c.geometry = nil
	c.motionBounds = nil
	c.streaming = false
}

func memoEntries(m *exactStateMemo) int {
	n := 0
	for i := range m.shards {
		s := &m.shards[i]
		s.mu.Lock()
		n += len(s.values)
		s.mu.Unlock()
	}
	return n
}
