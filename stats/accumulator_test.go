package stats

import (
	"reflect"
	"testing"
)

type factoryFunc func() Accumulator

func testClear(t *testing.T, accFn factoryFunc) {
	ac := accFn()

	ac.Clear()
	if 0 != ac.TotalCount() {
		t.Errorf("Unexpected count %d after Clear()", ac.TotalCount())
	}

	ac.Add(1)
	ac.Add(2)
	ac.Add(3)
	ac.Clear()
	if 0 != ac.TotalCount() {
		t.Errorf("Unexpected count %d after Clear()", ac.TotalCount())
	}
}

func testExceedMax(t *testing.T, accFn factoryFunc, max int) {
	tests := []struct {
		val int
		ok  bool
	}{
		{val: 0, ok: true},
		{val: 1, ok: true},
		{val: 10, ok: true},
		{val: 11, ok: false},
		{val: 99, ok: false},
	}

	for _, test := range tests {
		ac := accFn()
		err := ac.Add(test.val)
		if (err == nil) != test.ok {
			t.Errorf("Unexpected err of [max:%d] Add(%d) => %s",
				max, test.val, err)
		}
	}
}

func testAdd(t *testing.T, accFn factoryFunc) {
	tests := []struct {
		vals []int
	}{
		{[]int{}},
		{[]int{2}},
		{[]int{0, 1, 2, 3, 4}},
		{[]int{0, 0, 0, 10, 10, 10, 999, 999}},
	}

	for _, test := range tests {
		// test add one at a time
		ac := accFn()
		for _, v := range test.vals {
			err := ac.Add(v)
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
		}
		testVerifyDataset(t, ac, test.vals)

		// test batch add
		ac = accFn()
		err := ac.Add(test.vals...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		testVerifyDataset(t, ac, test.vals)
	}
}

func testAllData(t *testing.T, accFn factoryFunc) {
	tests := []struct {
		vals []int
	}{
		{[]int{}},
		{[]int{2}},
		{[]int{0, 1, 2, 3, 4}},
		{[]int{0, 0, 0, 10, 10, 10, 999, 999}},
	}

	for _, test := range tests {
		ac := accFn()
		for _, v := range test.vals {
			err := ac.Add(v)
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
		}

		if ac.TotalCount() != len(test.vals) {
			t.Errorf("Wrong number of data points; expected %d, actual %d",
				len(test.vals), ac.TotalCount())
		}
	}
}

func testSum(t *testing.T, accFn factoryFunc) {
	tests := []struct {
		vals     []int
		expected int
	}{
		{[]int{}, 0},
		{[]int{2}, 2},
		{[]int{1, 2, 3}, 6},
		{[]int{1, 9, 10, 5, 7, 8}, 40},
		{[]int{100, 100, 100, 100, 100, 100}, 600},
	}

	for _, test := range tests {
		ac := accFn()
		ac.Add(test.vals...)
		if test.expected != ac.Sum() {
			t.Errorf("Incorrect sum: expected %d, actual %d",
				test.expected, ac.Sum())
		}
	}
}

func testMean(t *testing.T, accFn factoryFunc) {
	tests := []struct {
		vals     []int
		expected float64
	}{
		{[]int{}, 0.0},
		{[]int{2}, 2.0},
		{[]int{5, 5, 5}, 5.0},
		{[]int{1, 5, 1, 5, 1, 5}, 3.0},
	}

	for _, test := range tests {
		ac := accFn()
		ac.Add(test.vals...)
		if test.expected != ac.Mean() { // TODO: account for epsilon?
			t.Errorf("Incorrect mean: expected %.4f, actual %.4f",
				test.expected, ac.Mean())
		}
	}
}

func testVerifyDataset(t *testing.T, ac Accumulator, expected []int) {
	if ac.TotalCount() != len(expected) {
		t.Errorf("Wrong number of data points; expected %d, actual %d",
			len(expected), ac.TotalCount())
	}

	if !reflect.DeepEqual(expected, ac.AllData()) {
		t.Errorf("Data points do not match expected: %v, actual %v",
			expected, ac.AllData())
	}
}
