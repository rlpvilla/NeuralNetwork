package ann

import (
//	"fmt"
//	"math"
)

type Activation struct {
	Function func(float64)float64
	Derivative func(float64)float64
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

func NewNeuron (insignals int, outsignals int, learnrate float64, activation Activation, peripherals Peripherals, cancelchan chan struct{}) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}

	go Nucleus (learnrate, activation, internals, cancelchan)
	go Axon (outsignals, internals.Output, peripherals.Output, cancelchan)
	go Dendrite (insignals, peripherals.Input, internals.Input, cancelchan)
	go Axon (insignals, internals.Downfeed, peripherals.Downfeed, cancelchan)
	go Dendrite (outsignals, internals.Upfeed, peripherals.Upfeed, cancelchan)
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
				select {
				case signalchan <- input:
				case <- cancelchan:
					return
				}
			}
		case <- cancelchan:
			return
		}
	}
}

func Synapse (startweight float64, inputchan chan float64, outputchan chan float64, downfeed chan float64, upfeed chan float64, cancelchan chan struct{}) {
	var weight float64 = startweight
	var output float64
	var inchan chan float64 = inputchan
	for {
		select {
		case input := <- inchan:
			output = input * weight
			select {
			case outputchan <- output:
				inchan = nil
			case <- cancelchan:
				return
			}
		case topmargin := <- upfeed:
			select {
			case downfeed <- topmargin * weight:
				select {
				case adjustment := <- upfeed:
					weight = weight + (adjustment * output)
					inchan = inputchan
				case <- cancelchan:
					return
				}
			case <- cancelchan:
				return
			}
		case <- cancelchan:
			return
		}
	}
}