package ann

import (
	"math"
	"math/rand"
	"time"
	"fmt"
)

const devnetwork bool = true

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

func Init () {
	cancelneuron := make(chan struct{})
	activation := Activation {Function: func (x float64) float64 {return 1/(1*math.Exp(-x))}, Derivative: func (x float64) float64 {return 1/(1*math.Exp(-x)) *(1 - 1/(1*math.Exp(-x)))}}
	neuron := Peripherals {Input: make(chan float64), Output: make(chan float64), Upfeed: make(chan float64), Downfeed: make(chan float64)}
	synapses := Synapses {Ingoing: 1, Outgoing: 1}; learningrate := 0.1

	NewNeuron(learningrate, synapses, activation, neuron, cancelneuron)

	select {
	case neuron.Input <- 1:
		result := <- neuron.Output
		if devnetwork {fmt.Printf("\n%v: Network output [%f]...\n", time.Now(), result)}
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
			if devnetwork {fmt.Printf("\n%v: Closing training session...\n")}
			close(cancelchan)
			return
		}
	}
}

func NewLayer (synapses Synapses) {

}