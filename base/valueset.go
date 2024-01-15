package base

import "fmt"

// ValueSet represents a set of values that appear in a sudoku Cell.
type ValueSet uint16

const MaxValue = 9

// NewValueSet returns a new ValueSet that contains every value in universe.
func NewValueSet(universe []int) ValueSet {
	var vs ValueSet = 0
	for _, v := range universe {
		vs = vs.SetHasValue(v, true)
	}
	return vs
}

// IsEmpty returns true if the ValueSet contains no values.
func (vs ValueSet) IsEmpty() bool {
	return vs == 0
}

// Len returns the number of values in the ValueSet.
func (vs ValueSet) Len() int {
	count := 0
	for i := 1; i <= MaxValue; i++ {
		if vs.HasValue(i) {
			count += 1
		}
	}
	return count
}

func bitmask(v int) ValueSet {
	if v < 1 {
		return 0
	}
	return 1 << uint16(v-1)
}

// HasValue returns true of vs contains v.
func (vs ValueSet) HasValue(v int) bool {
	m := bitmask(v)
	return vs&m == m
}

// SetHasValue returns a new ValueSet with v added to vs.
func (vs ValueSet) SetHasValue(v int, has bool) ValueSet {
	if has {
		return vs | bitmask(v)
	} else {
		return vs & ^bitmask(v)
	}
}

// Union eturns the union of two value sets.
func (vs1 ValueSet) Union(vs2 ValueSet) ValueSet {
	return vs1 | vs2
}

// Intersection eturns the intersection of two value sets.
func (vs1 ValueSet) Intersection(vs2 ValueSet) ValueSet {
	return vs1 & vs2
}

// SetDifference returns a ValueSet containing those values in vs1 which are not in vs2.
func (vs1 ValueSet) SetDifference(vs2 ValueSet) ValueSet {
	return vs1 & ^vs2
}

// DoValues calls f on each value in the ValueSet.  If f returns false
// then DoValues doen't call it on further values but returns immediately.
func (vs ValueSet) DoValues(f func(int) bool) {
	for v := 1; v <= MaxValue; v++ {
		if vs.HasValue(v) {
			if !f(v) {
				return
			}
		}
	}
}

func (vs ValueSet) String(separator string) string {
	s := ""
	vs.DoValues(func(v int) bool {
		if s != "" {
			s += separator
		}
		s += fmt.Sprintf("%d", v)
		return true
	})
	return s
}

// Get returns the indexth value from the ValueSet.
func (vs ValueSet) Get(index int) (int, error) {
	count := 0
	value := -1
	if index < 0 {
		panic(fmt.Sprintf("Index %d out of range", index))
	}
	if index >= vs.Len() {
		return 0, fmt.Errorf("Index %d out of range, 0x%x", index, vs)
	}
	f := func(v int) bool {
		if index == count {
			value = v
			return false
		}
		count += 1
		return true
	}
	vs.DoValues(f)
	return value, nil
}

// MustGet returns the indexth value from the ValueSet.
// MustGet panics if index is out of range.
func (vs ValueSet) MustGet(index int) int {
	v, err := vs.Get(index)
	if err != nil {
		panic(err)
	}
    return v
}
