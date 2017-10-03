package base

import "testing"

func TestValueSet(t *testing.T) {
	vs := NewValueSet([]int{2, 3})
	expectValue := func (vs ValueSet, v int, has bool) {
		if h := vs.HasValue(v); h != has {
			t.Errorf("HasValue(%d): want %v, got %v", v, has, h)
		}
	}
	expectValue(vs, 2, true)
	expectValue(vs, 3, true)
	expectValue(vs, 1, false)
	vs.SetHasValue(3, false)
	expectValue(vs, 3, false)
	vs.SetHasValue(1, false)
	expectValue(vs, 1, true)
	if got := vs.Len(); got != 2 {
		t.Errorf("Bad count,  got %d, want %d", got, 2)
	}
	if vs1 := NewValueSet([]int{4}); vs1 != 8 {
		t.Errorf("NewValueSet returned empty")
	}
}
