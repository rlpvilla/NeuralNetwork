package ann

import (
	"math"
	"rand"
)

type TrainingSet struct {
	Inputs []float64
	Expect float64
}

type Regimen struct {
	TrainingSets []TrainingSet
}

func Initialize () {
	cancelneurons := make(chan struct{}); cancelsynapses := make(chan struct{}); endtraining := make(chan struct{})
	activation := Activation {Function: func (x float64) float64 {return 1/(1*math.Exp(-x))}, Derivative: func (x float64) float64 {return 1}}
	neuron := Peripherals {Input: make(chan float64), Output: make(chan float64), Upfeed: make(chan float64), Downfeed: make(chan float64)}
	synapse_1 := Peripherals {Input: make(chan float64), Output: neuron.Input, Upfeed: neuron.Downfeed, Downfeed: make(chan float64)}
	synapse_2 := Peripherals {Input: make(chan float64), Output: neuron.Input, Upfeed: neuron.Downfeed, Downfeed: make(chan float64)}
	synapses := Synapses {Ingoing: 2, Outgoing: 1}
	learningrate := 0.1
	trainingchan := make(chan float64)

	NewNeuron(learningrate, synapses, activation, neuron, cancelneurons)
	go Synapse (0.75, synapse_1, cancelsynapses); go Synapse (0.25, synapse_2, cancelsynapses)

	go StaticInput (100, regimen, neuron, trainingchan, endtraining); go ErrorCatch (neuron, trainingchan, endtraining)

	select {
	case <- trainingchan:
		return
	}
}

func ErrorCatch (peripherals Peripherals, expectchan chan float64, cancelchan chan struct{}) {
	for {
		select {
		case result := <- peripherals.Output:
			expected := <- expectchan
			peripherals.Upfeed <- expected - result
		case <- cancelchan:
			return
		}
	}
}

func StaticInput (cycles int, regimen Regimen, peripherals Peripherals, expectchan chan float64, cancelchan chan struct{}) {
	rand.Seed(time.Now().UTC().UnixNano())
	sets := regimen.TrainingSets
	for {
		set := rand.Intn(len(sets))
		inputs := sets[set].Inputs
		input := rand.Intn(len(inputs))
		select {
		case peripherals.Input <- inputs[input]:
			expectchan <- sets[set].Expect
		case <- peripherals.Downfeed:
			continue
		case <- cancelchan:
			return
		}
		cycles--
		if cycles < 0 {
			close(cancelchan)
			return
		}
	}
}

func NewLayer (synapses Synapses) {

}