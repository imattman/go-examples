package stats

import "testing"

const defaultMax = 1000

func compactFactFn() Accumulator {
	return NewCompactAccumulator(defaultMax)
}

func TestClear_Compact(t *testing.T) {
	testClear(t, compactFactFn)
}

func TestExceedMax_Compact(t *testing.T) {
	max := 10
	testExceedMax(t,
		func() Accumulator {
			return NewCompactAccumulator(max)
		},
		max)
}

func TestAdd_Compact(t *testing.T) {
	testAdd(t, compactFactFn)
}

func TestAllData_Compact(t *testing.T) {
	testAllData(t, compactFactFn)
}

func TestSum_Compact(t *testing.T) {
	testSum(t, compactFactFn)
}

func TestMean_Compact(t *testing.T) {
	testMean(t, compactFactFn)
}
