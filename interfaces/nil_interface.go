package main

import (
	"fmt"
)

type Processor interface {
	Process()
}

type NoOp int

func (n *NoOp) Process() {}

func buildProcessor() Processor {
	var noop *NoOp = nil

	// appears to return nil but interface has a type and does not test as 'nil'
	return noop
}

func main() {
	proc := buildProcessor()
	if proc == nil {
		fmt.Printf("Did not expect nil: %v [%T]\n", proc, proc)
	}

	// An interface must have both a nil type and a nil value to be considered nil
	fmt.Printf("interface == nil  %v\n", proc == nil)
	fmt.Printf("interface type    %T\n", proc)
	fmt.Printf("interface value   %v\n", proc)

	fmt.Println()
	fmt.Printf("Switch type check: proc.(type) => ")
	switch t := proc.(type) {
	case *NoOp:
		fmt.Printf("*NoOp\n")
	default:
		fmt.Printf("unknown: %T %V\n", t, t)
	}
}
