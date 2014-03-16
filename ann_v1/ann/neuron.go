package ann

import (
//	"fmt"
//	"math"
)

type Activation struct {
	Function func(float64)float64
	Derivative func(float64)float64
}

type Synapses struct {
	Ingoing int
	Outgoing int
}

type Peripherals struct {
	Input chan float64
	Output chan float64
	Upfeed chan float64
	Downfeed chan float64
}


func PushOrCancel (value float64, outchan chan float64, cancelchan chan struct{}) bool {
	select {
	case outchan <- value:
		return true
	case <- cancelchan:
		return false
	}
}

func PullOrCancel (inchan chan float64, cancelchan chan struct{}) (bool, float64) {
	select {
	case value := <- inchan:
		return true, value
	case <- cancelchan:
		return false, 0
	}
}

func NewNeuron (synapses Synapses, learnrate float64, activation Activation, peripherals Peripherals, cancelchan chan struct{}) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}

	go Nucleus (learnrate, activation, internals, cancelchan)
	go Axon (synapses.Outgoing, internals.Output, peripherals.Output, cancelchan)
	go Dendrite (synapses.Ingoing, peripherals.Input, internals.Input, cancelchan)
	go Axon (synapses.Ingoing, internals.Downfeed, peripherals.Downfeed, cancelchan)
	go Dendrite (synapses.Outgoing, peripherals.Upfeed, internals.Upfeed, cancelchan)
}

func Nucleus (learnrate float64, activation Activation, peripherals Peripherals, cancelchan chan struct{}) {
	var excitement float64
	var inputchan = peripherals.Input
	for {
		select {
		case input := <- inputchan:
			excitement = input
			peripherals.Output <- activation.Function(excitement)
			inputchan = nil
		case errormargin := <- peripherals.Upfeed:
			peripherals.Downfeed <- errormargin
			peripherals.Downfeed <- learnrate * errormargin * activation.Derivative(excitement)
			inputchan = peripherals.Input
		case <- cancelchan:
			return
		}
	}
}

func Dendrite (signals int, inputchan chan float64, outputchan chan float64, cancelchan chan struct{}) {
	var signalcount int
	var sum float64
	for {
		select {
		case input := <- inputchan:
			signalcount++
			sum = sum + input
			if signalcount == signals {
				outputchan <- sum
				sum = 0
			}
		case <- cancelchan:
			return
		}
	}
}

func Axon (signals int, inputchan chan float64, signalchan chan float64, cancelchan chan struct{}) {
	for {
		select {
		case input := <- inputchan:
			for i := 0; i < signals; i++ {
				signalchan <- input
			}
		case <- cancelchan:
			return
		}
	}
}

func Synapse (weight float64, peripherals Peripherals, cancelchan chan struct{}) {
	var output float64
	var inputchan chan float64 = peripherals.Input
	for {
		select {
		case input := <- inputchan:
			output = input * weight
			peripherals.Output <- output
			inputchan = nil
		case topmargin := <- peripherals.Upfeed:
			peripherals.Downfeed <- topmargin * weight
			adjustment := <- peripherals.Upfeed
			weight = weight + (adjustment * output)
			inputchan = peripherals.Input
		case <- cancelchan:
			return
		}
	}
}