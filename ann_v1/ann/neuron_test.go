package ann

import (
	"testing"
	"time"
	"math"
)

func Test_Synapse_Input_Workflow (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 1
	var timeout time.Duration = 2

	go Synapse (startweight, internals, cancelchan)
	go func() {
		internals.Input <- 1
		<- internals.Output
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Synapse workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Synapse workflow is clear")
		return
	}
}

func Test_Synapse_Input_Weighting (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 0.3
	var timeout time.Duration = 1

	go Synapse (startweight, internals, cancelchan)
	go func() {
		internals.Input <- 1246
		result := <- internals.Output
		if result != 373.8 {
			t.Log("Failure - Synapse weighting returned wrong value")
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Synapse workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Synapse weighting is accurate")
		return
	}
}

func Test_Synapse_Input_Blocking (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 1
	var timeout time.Duration = 1000

	go Synapse (startweight, internals, cancelchan)
	go func() {
		internals.Input <- 1
		<- internals.Output
		internals.Input <- 1
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Microsecond):
		t.Log("Success - Synapse input is blocked")
		return
	case <- resultchan:
		t.Log("Failure - Synapse input not blocking")
		t.Fail()
		return
	}
}

func Test_Synapse_Feedback_Workflow (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 1
	var timeout time.Duration = 10

	go Synapse (startweight, internals, cancelchan)
	go func() {
		internals.Upfeed <- 1
		<- internals.Downfeed
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Synapse feedback workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Synapse feedback workflow is clear")
		return
	}
}

func Test_Synapse_Feedback_ErrorMargin (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 0.5
	var timeout time.Duration = 10

	go Synapse (startweight, internals, cancelchan)
	go func() {
		internals.Upfeed <- 0.25
		result := <- internals.Downfeed
		if result != .125 {
			t.Log("Failure - Synapse error margin inaccurate")
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Synapse error margin timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Synapse error margin is accurate")
		return
	}
}

func Test_Synapse_Loop_Workflow (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}
	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	var startweight float64 = 1
	var timeout time.Duration = 10

	go Synapse (startweight, internals, cancelchan)
	go func() {
		internals.Input <- 0; <- internals.Output
		internals.Upfeed <- 0; <- internals.Downfeed; internals.Upfeed <- 0
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Synapse loop workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Synapse loop workflow is clear")
		return
	}
}

func Test_Nucleus_Input_Workflow (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}

	activation := Activation {
		Function: func (x float64) float64 {return math.Pow(x, 2)},
		Derivative: func (x float64) float64 {return 2*x},
	}

	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})

	var learnrate float64 = 1
	var timeout time.Duration = 2

	go Nucleus (learnrate, activation, internals, cancelchan)
	go func() {
		internals.Input <- 1
		<- internals.Output
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Nucleus input workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Nucleus input workflow is clear")
		return
	}
}

func Test_Nucleus_Input_ComputeFunction (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}

	activation := Activation {
		Function: func (x float64) float64 {return math.Pow(x, 2)},
		Derivative: func (x float64) float64 {return 2*x},
	}

	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	
	var learnrate float64 = 1
	var timeout time.Duration = 2

	go Nucleus (learnrate, activation, internals, cancelchan)
	go func() {
		internals.Input <- 4.8
		result := <- internals.Output
		if result != 23.04 {
			t.Log("Failure -  Nucleus compute function is inaccurate")
			t.Log(result)
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Nucleus compute function timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Nucleus compute function is accurate")
		return
	}
}

func Test_Nucleus_Feedback_Workflow (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}

	activation := Activation {
		Function: func (x float64) float64 {return math.Pow(x, 2)},
		Derivative: func (x float64) float64 {return 2*x},
	}

	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	
	var learnrate float64 = 1
	var timeout time.Duration = 2

	go Nucleus (learnrate, activation, internals, cancelchan)
	go func() {
		internals.Upfeed <- 1
		<- internals.Downfeed
		<- internals.Downfeed
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Nucleus feedback workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Nucleus feedback workflow is accurate")
		return
	}
}

func Test_Nucleus_Feedback_ErrorMargin (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}

	activation := Activation {
		Function: func (x float64) float64 {return math.Pow(x, 2)},
		Derivative: func (x float64) float64 {return 2*x},
	}

	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	
	var learnrate float64 = 1
	var timeout time.Duration = 2

	go Nucleus (learnrate, activation, internals, cancelchan)
	go func() {
		internals.Upfeed <- 0.5
		result := <- internals.Downfeed
		if result != 0.5 {
			t.Log("Failure -  Nucleus margin error is inaccurate")
			t.Log(result)
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Nucleus margin error timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Nucleus margin error is accurate")
		return
	}
}

func Test_Nucleus_Feedback_ComputeDerivative (t *testing.T) {
	internals := Peripherals {
		Input: make(chan float64),
		Output: make(chan float64),
		Upfeed: make(chan float64),
		Downfeed: make(chan float64),
	}

	activation := Activation {
		Function: func (x float64) float64 {return 2*x},
		Derivative: func (x float64) float64 {return 2},
	}

	cancelchan := make(chan struct{})
	resultchan := make(chan struct{})
	
	var learnrate float64 = 0.5
	var timeout time.Duration = 2

	go Nucleus (learnrate, activation, internals, cancelchan)
	go func() {
		internals.Upfeed <- 0.5
		<- internals.Downfeed
		result := <- internals.Downfeed
		if result != 0.5 {
			t.Log("Failure -  Nucleus compute derivative is inaccurate")
			t.Log(result)
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Millisecond):
		t.Log("Failure - Nucleus compute derivative timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("Success - Nucleus compute derivative is accurate")
		return
	}
}