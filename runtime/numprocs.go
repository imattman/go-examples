package main

import (
	"fmt"
	"runtime"
)

func main() {
	// passing value of < 1 to GOMAXPROCS does not change the setting
	fmt.Println("logical procs:", runtime.GOMAXPROCS(0))
}
