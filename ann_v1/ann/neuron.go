package ann

import (
//	"fmt"
//	"math"
)

func NewNeuron (insignals int, outsignals int, learnrate float64, activationfunc func(float64)float64, derivativefunc func(float64)float64, inchan chan float64, outchan chan float64, downchan chan float64, upchan chan float64, cancelchan chan struct{}) {
	neuroninput := make(chan float64)
	neuronoutput := make(chan float64)
	neurondown := make(chan float64)
	neuronup := make(chan float64)

	go Neuron (learnrate, activationfunc, derivativefunc, neuroninput, neuronoutput, neurondown, neuronup, cancelchan)
	go Axon (outsignals, neuronoutput, outchan, cancelchan)
	go Dendrite (insignals, inchan, neuroninput, cancelchan)
	go Axon (insignals, neurondown, downchan, cancelchan)
	go Dendrite (outsignals, neuronup, upchan, cancelchan)
}

func Neuron (learnrate float64, activationfunc func(float64)float64, derivativefunc func(float64)float64, inputchan chan float64, outputchan chan float64, downfeed chan float64, upfeed chan float64, cancelchan chan struct{}) {
	var excitement float64
	var inchan = inputchan
	for {
		select {
		case input := <- inchan:
			excitement = input
			select {
			case outputchan <- activationfunc(excitement):
				inchan = nil
			case <- cancelchan:
				return
			}
		case difference := <- upfeed:
			select {
			case downfeed <- difference:
				select {
				case downfeed <- learnrate * difference * derivativefunc(excitement):
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

func Dendrite (signals int, inputchan chan float64, outputchan chan float64, cancelchan chan struct{}) {
	var signalcount int
	var sum float64
	for {
		select {
		case input := <- inputchan:
			signalcount++
			sum = sum + input
			if signalcount == signals {
				select {
				case outputchan <- sum:
					sum = 0
				case <- cancelchan:
					return
				}
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
		case difference := <- upfeed:
			select {
			case downfeed <- difference * weight:
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