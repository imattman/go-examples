package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
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
	return fmt.Sprintf("p(%d-%d) c(%d-%d) produce: %12s  consume: %12s",
		w.producer, w.pseq,
		w.consumer, w.cseq,
		w.prodTime, w.consTime)
}

func main() {
	var pcnt int
	var pAvg time.Duration
	var pVar time.Duration
	var ccnt int
	var cAvg time.Duration
	var cVar time.Duration

	var maxQuiet time.Duration

	flag.IntVar(&pcnt, "p", 1, "producer count")
	flag.DurationVar(&pAvg, "pavg", (100 * time.Millisecond), "average produce duration")
	flag.DurationVar(&pVar, "pvar", (50 * time.Millisecond), "produce duration variance")
	flag.IntVar(&ccnt, "c", 1, "consumer count")
	flag.DurationVar(&cAvg, "cavg", (100 * time.Millisecond), "average consume duration")
	flag.DurationVar(&cVar, "cvar", (50 * time.Millisecond), "consume duration variance")
	flag.DurationVar(&maxQuiet, "max", (2 * time.Second), "max quiet threshold of no data")
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

	for i := 0; i < 20; i++ {
		select {
		case work := <-finished:
			fmt.Printf("complete: %s\n", work)
		case <-time.After(maxQuiet):
			fmt.Printf("Exceeded %s max quiet period of no data\n", maxQuiet)
			cancelFunc()
			return
		}
	}
	cancelFunc()
}

func produce(ctx context.Context, id int, data chan<- WorkUnit, pAvg, pVar time.Duration) {
	for i := 0; true; i++ {
		prodTime := duration(pAvg, pVar)
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
		consTime := duration(cAvg, cVar)
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

func duration(avg, variance time.Duration) time.Duration {
	offset := avg.Nanoseconds() - variance.Nanoseconds()
	ubound := 2 * variance.Nanoseconds()
	if ubound <= 0 {
		ubound = 1
	}
	nsecs := rand.Int63n(ubound) + offset

	return time.Duration(nsecs)
}
