package ann

import (
	"math"
	"math/rand"
	"time"
	"fmt"
)

const devnetwork bool = false

var TestSet = Regimen {
	TrainingSets: []TrainingSet {
		{
			Input: 1,
			Expect: 1,
		},
		{
			Input: 0,
			Expect: 0,
		},
	},
}

type TrainingSet struct {
	Input float64
	Expect float64
}

type Regimen struct {
	TrainingSets []TrainingSet
}

func Initialize (regimen Regimen) {
	if devnetwork {fmt.Printf("\n%v: Network initialized...\n", time.Now())}
	cancelneurons := make(chan struct{}); cancelsynapses := make(chan struct{}); endtraining := make(chan struct{})
	activation := Activation {Function: func (x float64) float64 {return 1/(1*math.Exp(-x))}, Derivative: func (x float64) float64 {return 1/(1*math.Exp(-x)) *(1 - 1/(1*math.Exp(-x)))}}
	neuron := Peripherals {Input: make(chan float64), Output: make(chan float64), Upfeed: make(chan float64), Downfeed: make(chan float64)}
	synapse_1 := Peripherals {Input: make(chan float64), Output: neuron.Input, Upfeed: neuron.Downfeed, Downfeed: make(chan float64)}
	synapse_2 := Peripherals {Input: make(chan float64), Output: neuron.Input, Upfeed: neuron.Downfeed, Downfeed: make(chan float64)}
	synapses := Synapses {Ingoing: 2, Outgoing: 1}
	learningrate := 0.1
	trainingchan := make(chan float64)
	if devnetwork {fmt.Printf("\n%v: Network set up...\n", time.Now())}

	NewNeuron(learningrate, synapses, activation, neuron, cancelneurons); go Synapse (0.75, synapse_1, cancelsynapses); go Synapse (0.25, synapse_2, cancelsynapses)
	go StaticInput (1000, regimen, synapse_1, trainingchan, endtraining); go ErrorCatch (neuron, trainingchan, endtraining)
	go func() {
		if devnetwork {fmt.Printf("\n%v: Feeding bias...\n", time.Now())}
		count := 0
		for {
			count ++
			select {
			case synapse_2.Input <- 1:
				if devnetwork {fmt.Printf("\n%v: Input [%d] out...\n", time.Now(), count)}
				<- synapse_2.Downfeed
				if devnetwork {fmt.Printf("\n%v: Feedback [%d] received...\n", time.Now(), count)}
			case <- trainingchan:
				return
			}
		}
		}()

	select {
	case <- trainingchan:
		close(cancelsynapses); close(cancelneurons)
		return
	}
}

func ErrorCatch (peripherals Peripherals, expectchan chan float64, cancelchan chan struct{}) {
	for {
		select {
		case result := <- peripherals.Output:
			if devnetwork {fmt.Printf("\n%v: Result was [%d]...\n", time.Now(), result)}
			expected := <- expectchan
			if devnetwork {fmt.Printf("\n%v: Expected [%d]...\n", time.Now(), expected)}
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
		select {
		case peripherals.Input <- sets[set].Input:
			if devnetwork {fmt.Printf("\n%v: Input out...\n", time.Now())}
			expectchan <- sets[set].Expect
			if devnetwork {fmt.Printf("\n%v: Expected out...\n", time.Now())}
		case <- peripherals.Downfeed:
			if devnetwork {fmt.Printf("\n%v: BROKEN...\n", time.Now())}
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