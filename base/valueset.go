package base

type ValueSet uint16

const MaxValue = 9

func NewValueSet(universe []int) ValueSet {
	var vs ValueSet = 0
	for _, v := range universe {
		vs = vs.SetHasValue(v, true)
	}
	return vs
}

func (vs ValueSet) IsEmpty() bool {
	return vs == 0
}

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
	return 1 << uint16(v - 1)
}

func (vs ValueSet) HasValue(v int) bool {
	m := bitmask(v)
	return vs & m == m
}

func (vs ValueSet) SetHasValue(v int, has bool) ValueSet {
	if has {
		return vs | bitmask(v)
	} else {
		return vs & ^bitmask(v)
	}
}

func (vs ValueSet) DoValues(f func(int)) {
	for v := 1; v <= MaxValue; v++ {
		if vs.HasValue(v) {
			f(v)
		}
	}
}
