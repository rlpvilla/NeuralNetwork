package ann_v1

import (
	"fmt"
	"math"
)

type SynapseDef struct {
	Input chan float
	Output chan float
	Feedput chan float
	StartWeight float
}

type Neuron struct {
	ID string
	Network Network
	Errors chan error
}

type Network struct {
	ID string

}

func Synapse (startweight float, input chan float, output chan float, feedput chan float, cancelchan chan struct{}) {
	weight := startweight
	for {
		select {
		case value := <- input:
			output <- value * weight
			select {
			case feedback := <- feedput:

			case <- cancelchan:
				close(cancelchan)
				return
			}
		case <- cancelchan:
			close(cancelchan)
			return
		}
	}
}

func Neuron (inputs []chan float, outputs []chan float, infeeds []chan float, outfeeds []chan float, cancelchan chan struct {}) {
	input := make(chan float)
	output := make(chan float)

	for {
		select {
			
		}
	}
}

func Dendrite (inputs []chan float, output chan float, cancelchan chan struct{})


func Axon (outputs []chan float, cancelchan chan struct{})