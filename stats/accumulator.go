package stats

type Accumulator interface {
	Clear()
	Add(data ...int) error
	TotalCount() int
	AllData() []int
	Sum() int
	Mean() float64
}
