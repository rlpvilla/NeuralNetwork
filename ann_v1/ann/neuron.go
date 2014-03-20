package ann

import (
//	"math"
	"time"
	"fmt"
)

const devneuron bool = true
const devsynapse bool = true

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

func NewNeuron (learnrate float64, synapses Synapses, activation Activation, peripherals Peripherals, cancelchan chan struct{}) {
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
	if devneuron {fmt.Printf("\n%v: Neuron initialized...\n", time.Now())}
}

func Nucleus (learnrate float64, activation Activation, peripherals Peripherals, cancelchan chan struct{}) {
	var excitement float64
	for {
		select {
		case input := <- peripherals.Input:
			if devneuron {fmt.Printf("\n%v: Nucleus received [%f]...\n", time.Now(), input)}
			excitement = input
			peripherals.Output <- activation.Function(excitement)
			if devneuron {fmt.Printf("\n%v: Nucleus output [%f]...\n", time.Now(), activation.Function(excitement))}
		case errormargin := <- peripherals.Upfeed:
			if devneuron {fmt.Printf("\n%v: Nucleus upfed [%f]...\n", time.Now(), errormargin)}
			peripherals.Downfeed <- errormargin
			if devneuron {fmt.Printf("\n%v: Nucleus downfed [%f]...\n", time.Now(), errormargin)}
			peripherals.Downfeed <- learnrate * errormargin * activation.Derivative(excitement)
			if devneuron {fmt.Printf("\n%v: Nucleus received [%f]...\n", time.Now(), learnrate * errormargin * activation.Derivative(excitement))}
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
			if devneuron {fmt.Printf("\n%v: Dendrite received signal [%i]...\n", time.Now(), input)}
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
			if devneuron {fmt.Printf("\n%v: Axon relayed signal [%i]...\n", time.Now(), input)}
			for i := 0; i < signals; i++ {
				signalchan <- input
				if devneuron {fmt.Printf("\n%v: Axon relayed signal [%i]...\n", time.Now(), input)}
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
			if devsynapse {fmt.Printf("\n%v: Synapse received input [%f]...\n", time.Now(), input)}
			output = input * weight
			peripherals.Output <- output
			if devsynapse {fmt.Printf("\n%v: Synapse relayed signal [%f]...\n", time.Now(), output)}
			inputchan = nil
		case errormargin := <- peripherals.Upfeed:
			if devsynapse {fmt.Printf("\n%v: Synapse received error margin [%f]...\n", time.Now(), errormargin)}
			peripherals.Downfeed <- errormargin * weight
			if devsynapse {fmt.Printf("\n%v: Synapse relayed margin [%f]...\n", time.Now(), errormargin * weight)}
			adjustment := <- peripherals.Upfeed
			if devsynapse {fmt.Printf("\n%v: Synapse received adjustment [%f]...\n", time.Now(), adjustment)}
			weight = weight + (adjustment * output)
			inputchan = peripherals.Input
		case <- cancelchan:
			return
		}
	}
}