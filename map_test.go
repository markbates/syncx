package syncx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Map_Len(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	m := &Map[int, int]{}

	exp := 0
	act := m.Len()

	r.Equal(exp, act)

	r.NoError(m.Set(1, 1))

	exp = 1
	act = m.Len()

	r.Equal(exp, act)

	ok := m.Delete(1)
	r.True(ok)

	exp = 0
	act = m.Len()

	r.Equal(exp, act)
}

func Test_Map_Keys(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	m := &Map[int, int]{}

	r.NoError(m.Set(2, 200))

	r.NoError(m.Set(1, 100))

	r.NoError(m.Set(3, 300))

	exp := []int{1, 2, 3}

	act := m.Keys()

	r.Equal(len(exp), len(act))
	for i := range exp {
		r.Equal(exp[i], act[i])
	}
}

func Test_NewMap(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	m := NewMap(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})

	r.Equal(3, m.Len())
}

func Test_Map_Get_Set_Delete(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	m := &Map[int, int]{}

	exp := 42

	_, ok := m.Get(1)
	r.False(ok)

	r.NoError(m.Set(1, exp))

	act, ok := m.Get(1)
	r.True(ok)

	r.Equal(exp, act)

	ok = m.Delete(1)
	r.True(ok)

	_, ok = m.Get(1)
	r.False(ok)

	ok = m.Delete(1)
	r.False(ok)
}

func Test_Map_Range(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	exp := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	m := NewMap(exp)

	m.Range(func(k int, v string) bool {
		r.Equal(exp[k], v)

		return true
	})

}

func Test_Map_Range_Break(t *testing.T) {
	t.Parallel()

	exp := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	m := NewMap(exp)

	var count int
	m.Range(func(k int, v string) bool {
		count++
		return false
	})

	if count != 1 {
		t.Fatalf("expected 1, got %v", count)
	}
}
func Test_Map_Clear(t *testing.T) {
	t.Parallel()

	r := require.New(t)

	m := NewMap(map[int]string{
		1: "one",
	})

	v, ok := m.Get(1)
	r.True(ok)

	r.Equal("one", v)

	m.Clear()

	_, ok = m.Get(1)
	r.False(ok)

}

func Test_Map_BulkSet(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	m := &Map[int, string]{}

	bm := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	err := m.BulkSet(bm)

	r.NoError(err)

	r.Equal(len(bm), m.Len())

	for k, e := range bm {
		a, ok := m.Get(k)
		r.True(ok)
		r.Equal(e, a)
	}

}
