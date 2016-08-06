package stats

import "fmt"

type BoundedDataset struct {
	data  []int
	count int
}

func NewBoundedDataset(maxVal int) *BoundedDataset {
	return &BoundedDataset{data: make([]int, maxVal+1)}

}

func (ds *BoundedDataset) Add(datapoints ...int) error {
	maxAllowed := ds.maxAllowed()
	for _, v := range datapoints {
		if v < 0 {
			return fmt.Errorf("illegal negative value %d", v)
		}
		if v > maxAllowed {
			return fmt.Errorf("%d exceeds max allowed value %d", v, maxAllowed)
		}

		ds.data[v] += 1
		ds.count++
	}

	return nil
}

func (ds *BoundedDataset) maxAllowed() int {
	return len(ds.data) - 1
}

func (ds *BoundedDataset) TotalCount() int {
	return ds.count
}

func (ds *BoundedDataset) AllData() []int {
	all := make([]int, 0, ds.count)

	// recreate dataset by using counts
	for i, cnt := range ds.data {
		for j := 0; j < cnt; j++ {
			all = append(all, i)
		}
	}

	return all
}

func (ds *BoundedDataset) Sum() int {
	sum := 0
	for i, cnt := range ds.data {
		sum += i * cnt
	}

	return sum
}

func (ds *BoundedDataset) Mean() float64 {
	if ds.count == 0 {
		return 0.0
	}
	return float64(ds.Sum()) / float64(ds.count)
}

func (ds *BoundedDataset) String() string {
	return fmt.Sprintf("dataset(max: %d) total: %d, sum: %d, mean: %0.3f",
		ds.maxAllowed(), ds.TotalCount(), ds.Sum(), ds.Mean())
}
