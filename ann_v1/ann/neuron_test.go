package ann

import (
	"testing"
	"time"
)

func Test_Synapse_InputWorkflow (t *testing.T) {
	inputchan := make(chan float64)
	outputchan := make(chan float64)
	downfeed := make(chan float64)
	upfeed := make(chan float64)
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 1
	var timeout time.Duration = 2

	go Synapse (startweight, inputchan, outputchan, downfeed, upfeed, cancelchan)
	go func() {
		inputchan <- 1
		<- outputchan
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Second):
		t.Log("\nFailure - Synapse workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("\nSuccess - Synapse workflow is clear")
		return
	}
}

func Test_Synapse_InputWeighting (t *testing.T) {
	inputchan := make(chan float64)
	outputchan := make(chan float64)
	downfeed := make(chan float64)
	upfeed := make(chan float64)
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 0.3
	var timeout time.Duration = 1

	go Synapse (startweight, inputchan, outputchan, downfeed, upfeed, cancelchan)
	go func() {
		inputchan <- 1246
		result := <- outputchan
		if result != 373.8 {
			t.Log("\nFailure - Synapse weighting returned wrong value")
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Second):
		t.Log("\nFailure - Synapse workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("\nSuccess - Synapse weighting is accurate")
		return
	}
}

func Test_Synapse_InputBlocking (t *testing.T) {
	inputchan := make(chan float64)
	outputchan := make(chan float64)
	downfeed := make(chan float64)
	upfeed := make(chan float64)
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 1
	var timeout time.Duration = 1000

	go Synapse (startweight, inputchan, outputchan, downfeed, upfeed, cancelchan)
	go func() {
		inputchan <- 1
		<- outputchan
		inputchan <- 1
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Microsecond):
		t.Log("\nSuccess - Synapse input is blocked")
		return
	case <- resultchan:
		t.Log("\nFailure - Synapse input not blocking")
		t.Fail()
		return
	}
}

func Test_Synapse_FeedbackWorkflow (t *testing.T) {
	inputchan := make(chan float64)
	outputchan := make(chan float64)
	downfeed := make(chan float64)
	upfeed := make(chan float64)
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 1
	var timeout time.Duration = 10

	go Synapse (startweight, inputchan, outputchan, downfeed, upfeed, cancelchan)
	go func() {
		upfeed <- 1
		<- downfeed
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("\nFailure - Synapse feedback workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("\nSuccess - Synapse feedback workflow is clear")
		return
	}
}

func Test_Synapse_FeedbackMargin (t *testing.T) {
	inputchan := make(chan float64)
	outputchan := make(chan float64)
	downfeed := make(chan float64)
	upfeed := make(chan float64)
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 0.5
	var timeout time.Duration = 10

	go Synapse (startweight, inputchan, outputchan, downfeed, upfeed, cancelchan)
	go func() {
		upfeed <- 0.25
		result := <- downfeed
		if result != .125 {
			t.Log("\nFailure - Synapse weighting returned wrong value")
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("\nFailure - Synapse feedback margin timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("\nSuccess - Synapse feedback margin is accurate")
		return
	}
}

func Test_Nucleus_InputWorkflow (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}
	activation := Activation {
		Function: func (x float64) {return 2*x},
		Derivative: func (x float64) {return 2},
	}
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var learnrate float64 = 1
	var timeout time.Duration = 2

	go Nucleus (learnrate, activation, internals, cancelchan)
	go func() {
		inputchan <- 1
		<- outputchan
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Second):
		t.Log("\nFailure - Synapse workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("\nSuccess - Synapse workflow is clear")
		return
	}
}