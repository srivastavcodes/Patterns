package main

import (
	"fmt"
	"math/rand"
	"testing"
)

// repeatAndTake showcases how repeat and take pipelines can be used together
// to create something useful.
func TestRepeatAndTake(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	// generate 10 wanted numbers
	for res := range take(done, repeat(done, 5), 10) {
		fmt.Printf("%v", res)
	}
	println()
	// generate 10 random numbers.
	randInt := func() any { return rand.Int() }
	for res := range take(done, repeatFn(done, randInt), 10) {
		fmt.Printf("%v\n", res)
	}
	for str := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		fmt.Printf("%s", str)
	}
	println()
}

// repeat repeats the values you pass to you indefinitely unless you tell it
// to stop.
func repeat(done <-chan struct{}, values ...any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream)
		for {
			for _, val := range values {
				select {
				case valueStream <- val:
					// keeps sending the values unless the reader stops
					// reader and then the channel will block and wait
					// for either done being called or for reading to
					// start again.
				case <-done:
					return
				}
			}
		}
	}()
	return valueStream
}

// take takes the first count items off of the incoming stream and then exit.
func take(done <-chan struct{}, inputStream <-chan any, count int) <-chan any {
	takeStream := make(chan any)
	go func() {
		defer close(takeStream)
		for i := 0; i < count; i++ {
			select {
			case takeStream <- <-inputStream:
			// this is what controls the repeat from generating values
			// infinitely, because after the loop exists, there is
			// nothing that is there to read from the inputStream which
			// blocks the channel.
			case <-done:
				return
			}
		}
	}()
	return takeStream
}

// repeatFn calls the given function indefinitely until stopped and sends the
// result of the execution to the stream.
func repeatFn(done <-chan struct{}, fn func() any) <-chan any {
	valueStream := make(chan any)
	go func() {
		defer close(valueStream)
		for {
			select {
			case valueStream <- fn():
				// the same principle as the above repeat, but
				// sends the result of the function call to
				// the stream.
			case <-done:
				return
			}
		}
	}()
	return valueStream
}

// toString asserts that the values from the incoming stream are strings and
// sends them to the outgoing stream.
func toString(done <-chan struct{}, inputStream <-chan any) <-chan string {
	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range inputStream {
			select {
			case stringStream <- v.(string):
			case <-done:
				return
			}
		}
	}()
	return stringStream
}
