package syncx

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// NewMap returns a new Map from the given map.
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

func (m *Map[K, V]) init() {
	if m == nil {
		return
	}

	m.once.Do(func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.m == nil {
			m.m = make(map[K]V)
		}
	})
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
func (m *Map[K, V]) Set(k K, v V) {
	if m == nil {
		return
	}

	m.init()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.m[k] = v
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
func (m *Map[K, V]) Range(f func(k K, v V) bool) {
	if m == nil || f == nil {
		return
	}

	m.init()

	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
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

	m.m = make(map[K]V)
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