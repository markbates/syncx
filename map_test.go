package syncx

import "testing"

func Test_Map_Len(t *testing.T) {
	t.Parallel()

	m := &Map[int, int]{}

	exp := 0
	act := m.Len()

	if exp != act {
		t.Fatalf("expected %v, got %v", exp, act)
	}

	m.Set(1, 1)

	exp = 1
	act = m.Len()

	if exp != act {
		t.Fatalf("expected %v, got %v", exp, act)
	}

	ok := m.Delete(1)
	if !ok {
		t.Fatal("expected true")
	}

	exp = 0
	act = m.Len()

	if exp != act {
		t.Fatalf("expected %v, got %v", exp, act)
	}

}

func Test_Map_Keys(t *testing.T) {
	t.Parallel()

	m := &Map[int, int]{}

	m.Set(2, 200)
	m.Set(1, 100)
	m.Set(3, 300)

	exp := []int{1, 2, 3}

	act := m.Keys()

	if len(act) != len(exp) {
		t.Fatalf("expected %v, got %v", exp, act)
	}

	for i := range exp {
		if act[i] != exp[i] {
			t.Fatalf("expected %v, got %v", exp, act)
		}
	}
}

func Test_NewMap(t *testing.T) {
	t.Parallel()

	m := NewMap(map[int]string{
		1: "one",
		2: "two",
		3: "three",
	})

	if m.Len() != 3 {
		t.Fatalf("expected 3, got %v", m.Len())
	}
}

func Test_Map_Get_Set_Delete(t *testing.T) {
	t.Parallel()

	m := &Map[int, int]{}

	exp := 42

	_, ok := m.Get(1)
	if ok {
		t.Fatalf("expected false, got %v", ok)
	}

	m.Set(1, exp)

	act, ok := m.Get(1)
	if !ok {
		t.Fatalf("expected true, got %v", ok)
	}

	if act != exp {
		t.Fatalf("expected %v, got %v", exp, act)
	}

	ok = m.Delete(1)
	if !ok {
		t.Fatal("expected true")
	}

	_, ok = m.Get(1)
	if ok {
		t.Fatalf("expected false, got %v", ok)
	}

	ok = m.Delete(1)
	if ok {
		t.Fatalf("expected false, got %v", ok)
	}
}

func Test_Map_Range(t *testing.T) {
	t.Parallel()

	exp := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	m := NewMap(exp)

	m.Range(func(k int, v string) bool {
		if exp[k] != v {
			t.Fatalf("expected %v, got %v", exp[k], v)
		}

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

	m := NewMap(map[int]string{
		1: "one",
	})

	v, ok := m.Get(1)
	if !ok {
		t.Fatalf("expected true, got %v", ok)
	}

	if v != "one" {
		t.Fatalf("expected one, got %v", v)
	}

	m.Clear()

	_, ok = m.Get(1)
	if ok {
		t.Fatalf("expected false, got %v", ok)
	}

}
