package base

import "testing"

func TestValueSet(t *testing.T) {
	values := []int{2, 3}
	vs := NewValueSet(values)
	t.Logf("vs = %03x", vs)
	expectValue := func (vs ValueSet, v int, has bool) {
		if h := vs.HasValue(v); h != has {
			t.Errorf("HasValue(%d): want %v, got %v", v, has, h)
		}
	}
	expectValue(vs, 2, true)
	expectValue(vs, 3, true)
	expectValue(vs, 1, false)
	vs = vs.SetHasValue(3, false)
	t.Logf("vs = %03x", vs)
	expectValue(vs, 3, false)
	vs = vs.SetHasValue(1, true)
	t.Logf("vs = %03x", vs)
	expectValue(vs, 1, true)
	if got := vs.Len(); got != 2 {
		t.Errorf("Bad count,  got %d, want %d", got, 2)
	}
	if vs1 := NewValueSet([]int{4}); vs1 != 8 {
		t.Errorf("NewValueSet returned empty")
	}

	vs = NewValueSet(values)
	valueIndex := 0
	vs.DoValues(func(v int) {
		t.Logf("DoValues function got %d from %03x", v, vs)
		if valueIndex < len(values) {
			if want := values[valueIndex]; want != v {
				t.Errorf("Wrong value in DoValues function, got %d, want %d", v, want)
			}
		}
		valueIndex += 1
	})
	if valueIndex != len(values) {
		t.Errorf("DoValues function called %d times, expected %d times", valueIndex, len(values))
	}
}
