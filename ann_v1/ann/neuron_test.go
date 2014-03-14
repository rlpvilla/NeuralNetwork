package ann

import (
	"testing"
	"time"
)

func Test_SynapseWorkflow (t *testing.T) {
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
		t.Log("\nFailure - Synapse Workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("\nSuccess - Synapse Workflow is clear")
		return
	}
}

func Test_SynapseWeighting (t *testing.T) {
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
		compare := <- outputchan
		if compare != 373.8 {
			t.Log("\nFailure - Synapse Weighting returned wrong value")
			return
		}
		resultchan <- struct{}{}
		}()

	select {
	case <- time.After(timeout * time.Second):
		t.Log("\nFailure - Synapse Workflow timed out")
		t.Fail()
		return
	case <- resultchan:
		t.Log("\nSuccess - Synapse Weighting is accurate")
		return
	}
}