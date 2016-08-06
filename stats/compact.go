package stats

import "fmt"

type CompactAccumulator struct {
	data  []int
	count int
}

func NewCompactAccumulator(maxVal int) *CompactAccumulator {
	return &CompactAccumulator{data: make([]int, maxVal+1)}

}

func (ac *CompactAccumulator) Clear() {
	ac.data = make([]int, len(ac.data))
	ac.count = 0
}

func (ac *CompactAccumulator) Add(datapoints ...int) error {
	maxAllowed := ac.maxAllowed()
	for _, v := range datapoints {
		if v < 0 {
			return fmt.Errorf("illegal negative value %d", v)
		}
		if v > maxAllowed {
			return fmt.Errorf("%d exceeds max allowed value %d", v, maxAllowed)
		}

		ac.data[v]++
		ac.count++
	}

	return nil
}

func (ac *CompactAccumulator) maxAllowed() int {
	return len(ac.data) - 1
}

func (ac *CompactAccumulator) TotalCount() int {
	return ac.count
}

func (ac *CompactAccumulator) AllData() []int {
	all := make([]int, 0, ac.count)

	// recreate dataset by using counts
	for i, cnt := range ac.data {
		for j := 0; j < cnt; j++ {
			all = append(all, i)
		}
	}

	return all
}

func (ac *CompactAccumulator) Sum() int {
	sum := 0
	for i, cnt := range ac.data {
		sum += i * cnt
	}

	return sum
}

func (ac *CompactAccumulator) Mean() float64 {
	if ac.count == 0 {
		return 0.0
	}
	return float64(ac.Sum()) / float64(ac.count)
}

func (ac *CompactAccumulator) String() string {
	return fmt.Sprintf("acc(<%d) count: %d  sum: %d  mean: %0.3f",
		ac.maxAllowed(), ac.TotalCount(), ac.Sum(), ac.Mean())
}
