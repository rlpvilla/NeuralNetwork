package ann

import (

)

func Initialize () {
	cancelneurons := make(chan struct{}); cancelsynapses := make(chan struct{})
	neural := Peripherals {Input: make(chan float64), Output: make(chan float64), Upfeed: make(chan float64), Downfeed: make(chan float64)}
	activation := Activation {Function: func (x float64) float64 {return x}, Derivative: func (x float64) float64 {return 1}}
	synapse_1 := Peripherals {Input: make(chan float64), Output: neural.Input, Upfeed: neural.Downfeed, Downfeed: make(chan float64)}
	synapse_2 := Peripherals {Input: make(chan float64), Output: neural.Input, Upfeed: neural.Downfeed, Downfeed: make(chan float64)}
	synapses := Synapses {Ingoing: 2, Outgoing: 1}

	NewNeuron(1, synapses, activation, neural, cancelneurons)
	go Synapse (0.75, synapse_1, cancelsynapses); go Synapse (0.25, synapse_2, cancelsynapses)
}