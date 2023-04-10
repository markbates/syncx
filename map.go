package syncx

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// NewMap returns a new Map from the given map.
// If the map is nil, a new map is created.
func NewMap[K constraints.Ordered, V any](m map[K]V) *Map[K, V] {
	if m == nil {
		m = map[K]V{}
	}

	return &Map[K, V]{
		m: m,
	}
}

// Map is a synchronized map.
type Map[K constraints.Ordered, V any] struct {
	m    map[K]V
	mu   sync.RWMutex
	once sync.Once
}

// Get returns the value for the given key.
// Returns false if the key does not exist in the map.
func (m *Map[K, V]) Get(k K) (V, bool) {
	if m == nil {
		return *new(V), false
	}

	m.init()

	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.m[k]
	return v, ok
}

// Set sets the value for the given key.
func (m *Map[K, V]) Set(k K, v V) error {
	if m == nil {
		return ErrNilMap
	}

	m.init()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[k] = v

	return nil
}

func (m *Map[K, V]) BulkSet(kvs map[K]V) error {
	if m == nil {
		return ErrNilMap
	}

	m.init()

	m.mu.Lock()

	for k, v := range kvs {
		m.m[k] = v
	}

	m.mu.Unlock()

	return nil
}

// Delete removes the item for the given key.
// Returns true if an item was in the map and deleted.
func (m *Map[K, V]) Delete(k K) bool {
	if m == nil {
		return false
	}

	m.init()

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.m[k]; !ok {
		return false
	}

	delete(m.m, k)

	return true
}

// Range calls f sequentially for each key and value present in the map.
// Returning true from f will terminate the iteration.
func (m *Map[K, V]) Range(f func(k K, v V) bool) bool {
	if m == nil || f == nil {
		return false
	}

	m.init()

	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.m {
		if !f(k, v) {
			return false
		}
	}

	return true
}

// Len returns the number of items in the map.
func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}

	m.init()

	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.m)
}

// Clear removes all items from the map.
func (m *Map[K, V]) Clear() {
	if m == nil {
		return
	}

	m.init()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.m = map[K]V{}
}

// Keys returns a sorted slice of the keys in the map.
func (m *Map[K, V]) Keys() []K {
	if m == nil {
		return nil
	}

	m.init()

	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < (keys[j])
	})

	return keys
}

// Clone returns a new Map with a copy of the underlying map.
func (m *Map[K, V]) Clone() (*Map[K, V], error) {
	if m == nil {
		return nil, ErrNilMap
	}

	nm := &Map[K, V]{
		m: map[K]V{},
	}

	var err error
	m.Range(func(k K, v V) bool {
		if err = nm.Set(k, v); err != nil {
			return false
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return nm, nil
}

func (m *Map[K, V]) Map() map[K]V {
	mm := map[K]V{}
	if m == nil {
		return mm
	}

	m.Range(func(k K, v V) bool {
		mm[k] = v
		return true
	})

	return mm
}

func (m *Map[K, V]) init() {
	if m == nil {
		return
	}

	m.once.Do(func() {
		m.mu.Lock()
		defer m.mu.Unlock()

		if m.m == nil {
			m.m = map[K]V{}
		}
	})
}
