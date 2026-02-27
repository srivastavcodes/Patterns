package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// ReplicatedRequest is a pattern where you can spin up multiple processes for
// concurrent processing and wait for the first process that returns and
// discard the rest. This process needs to ensure that whatever data the handlers
// are using needs to be replicated as well.
// If your handlers are operating on something that is too uniform as in the
// requests are too much alike, there's a very small possibility that you'll see
// any difference in response time, as everything will happen the same, however,
// if you are accessing a different path, processes, or access to a different data
// store altogether; using this pattern can improve performance significantly.
func TestReplicatedRequest(t *testing.T) {
	doWork := func(done <-chan struct{}, id int, wg *sync.WaitGroup, result chan<- int) {
		started := time.Now()
		defer wg.Done()

		// simulate random load
		simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
		select {
		case <-time.After(simulatedLoadTime):
		case <-done:
		}
		select {
		case result <- id:
		case <-done:
		}
		took := time.Since(started)
		// to showcase how long it would have taken had done not been closed.
		if took < simulatedLoadTime {
			took = simulatedLoadTime
		}
		fmt.Printf("%v took %v\n", id, took)
	}
	done := make(chan struct{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)
	// here we schedule 10 handlers to handle our request.
	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}
	// we received the first result returned from one of the handlers.
	returnId := <-result
	// we close to stop ongoing unnecessary computation and processes.
	close(done)
	wg.Wait()

	fmt.Printf("received result: %d\n", returnId)
}
