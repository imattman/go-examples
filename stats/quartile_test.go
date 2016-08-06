package stats

import (
	"reflect"
	"testing"
)

func TestExceedMax(t *testing.T) {
	max := 10
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

	ds := NewBoundedDataset(max)

	for _, test := range tests {
		err := ds.Add(test.val)
		if (err == nil) != test.ok {
			t.Errorf("Unexpected err of [max:%d] Add(%d) => %s",
				max, test.val, err)
		}
	}
}

func TestAdd(t *testing.T) {
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
		ds := NewBoundedDataset(1000)
		for _, v := range test.vals {
			err := ds.Add(v)
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
		}
		testVerifyDataset(t, ds, test.vals)

		// test batch add
		ds2 := NewBoundedDataset(1000)
		err := ds2.Add(test.vals...)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
		}
		testVerifyDataset(t, ds2, test.vals)
	}
}

func TestAllData(t *testing.T) {
	tests := []struct {
		vals []int
	}{
		{[]int{}},
		{[]int{2}},
		{[]int{0, 1, 2, 3, 4}},
		{[]int{0, 0, 0, 10, 10, 10, 999, 999}},
	}

	for _, test := range tests {
		ds := NewBoundedDataset(1000)
		for _, v := range test.vals {
			err := ds.Add(v)
			if err != nil {
				t.Errorf("Unexpected error %s", err)
			}
		}

		if ds.TotalCount() != len(test.vals) {
			t.Errorf("Wrong number of data points; expected %d, actual %d",
				len(test.vals), ds.TotalCount())
		}
	}
}

func TestSum(t *testing.T) {
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
		ds := NewBoundedDataset(100)
		ds.Add(test.vals...)
		if test.expected != ds.Sum() {
			t.Errorf("Incorrect sum: expected %.4f, actual %.4f",
				test.expected, ds.Sum())
		}
	}
}

func TestMean(t *testing.T) {
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
		ds := NewBoundedDataset(100)
		ds.Add(test.vals...)
		if test.expected != ds.Mean() { // TODO: account for epsilon?
			t.Errorf("Incorrect mean: expected %.4f, actual %.4f",
				test.expected, ds.Mean())
		}
	}
}

func testVerifyDataset(t *testing.T, ds *BoundedDataset, expData []int) {
	if ds.TotalCount() != len(expData) {
		t.Errorf("Wrong number of data points; expected %d, actual %d",
			len(expData), ds.TotalCount())
	}

	if !reflect.DeepEqual(expData, ds.AllData()) {
		t.Errorf("Data points do not match expected: %s, actual %s",
			expData, ds.AllData())
	}
}
