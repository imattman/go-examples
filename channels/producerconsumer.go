package main

import (
	"context"
	"flag"
	"fmt"
	"time"
)

type WorkUnit struct {
	producer int
	pseq     int
	prodTime time.Duration
	consumer int
	cseq     int
	consTime time.Duration
}

func (w WorkUnit) String() string {
	return fmt.Sprintf("p(%d-%d) c(%d-%d)", w.producer, w.pseq, w.consumer, w.cseq)
}

func main() {
	var pcnt int
	var pAvg time.Duration
	var pVar time.Duration
	var ccnt int
	var cAvg time.Duration
	var cVar time.Duration

	flag.IntVar(&pcnt, "p", 1, "producer count")
	flag.DurationVar(&pAvg, "pavg", (100 * time.Millisecond), "average produce duration")
	flag.DurationVar(&pAvg, "pvar", (50 * time.Millisecond), "produce duration variance")
	flag.IntVar(&ccnt, "c", 1, "consumer count")
	flag.DurationVar(&cAvg, "cavg", (100 * time.Millisecond), "average consume duration")
	flag.DurationVar(&cAvg, "cvar", (50 * time.Millisecond), "consume duration variance")
	flag.Parse()

	prodcons := make(chan WorkUnit)
	finished := make(chan WorkUnit, ccnt)
	ctx, cancelFunc := context.WithCancel(context.Background())

	for i := 1; i <= ccnt; i++ {
		go consume(ctx, i, prodcons, finished, cAvg, cVar)
	}

	for i := 1; i <= pcnt; i++ {
		go produce(ctx, i, prodcons, pAvg, pVar)
	}

	maxQuiet := 2 * time.Second
	for i := 0; i < 20; i++ {
		select {
		case work := <-finished:
			fmt.Printf("complete: %s\n", work)
		case <-time.After(1 * time.Second):
			fmt.Printf("Exceeded %s max quiet period of no data\n", maxQuiet)
			cancelFunc()
			return
		}
	}
	cancelFunc()
}

func produce(ctx context.Context, id int, data chan<- WorkUnit, pAvg, pVar time.Duration) {
	for i := 0; true; i++ {
		prodTime := pAvg

		time.Sleep(prodTime)
		w := WorkUnit{producer: id, pseq: i, prodTime: prodTime}
		select {
		case data <- w:
		case <-ctx.Done():
			fmt.Printf("cancelling producer %d\n", id)
			return
		}
	}
}

func consume(ctx context.Context, id int, data <-chan WorkUnit, collect chan<- WorkUnit, cAvg, cVar time.Duration) {
	for i := 0; true; i++ {
		consTime := cAvg
		select {
		case w := <-data:
			w.consumer = id
			w.cseq = i
			w.consTime = consTime
			time.Sleep(consTime)
			collect <- w
		case <-ctx.Done():
			fmt.Printf("cancelling consumer %d\n", id)
			return
		}
	}
}
