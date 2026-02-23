package main

import (
	"fmt"
	"net/http"
)

/*
NOTE:
	Errors should be considered first-class citizens when constructing values to
	return from goroutines. If your goroutine can produce errors, those errors
	should be tightly coupled with your result type and passed along through the
	same lines of communication – just like regular synchronous functions.
*/

// withoutErrorHandling showcases a general problem of error handling in concurrent
// programs; the error can't be passed to any other process, and there is no limit
// as to how many errors are too many.
// The best case scenario would be that someone is paying attention to the logs,
// else the error gets ignored, and that could be costly.
func withoutErrorHandling() {
	checkStatus := func(done <-chan struct{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					// can't return the error so it logs it. no fixed
					// limit as to when to break, and can't pass it
					// to anything else as well.
					fmt.Printf("couldn't fetch url: %v\n", err)
					continue
				}
				select {
				case responses <- resp:
				case <-done:
					return
				}
			}
		}()
		return responses
	}
	done := make(chan struct{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	responses := checkStatus(done, urls...)

	for resp := range responses {
		fmt.Println(resp.Status)
	}
}

// Result is a data type that is supposed to be aware of the application's state,
// say (backend). The error will be passed to the Err field and Result decides
// how to move forward.
type Result struct {
	Err  error
	Resp *http.Response
}

// withErrorHandling demonstrates how you could handle error in a goroutine by
// passing it to another part of your program that is aware of the application
// state and can make decisions accordingly.
func withErrorHandling() {
	// The error handling gets completely separated from the producer goroutine,
	// which is desirable because the process that spawned this producer
	// goroutine is more aware of the application's context and can handle the
	// error better.
	checkStatus := func(done <-chan struct{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				// passes the entire result of the operation to the
				// responsible part of the program where, decision
				// regarding errors/results can be made.
				result := Result{
					Resp: resp, Err: err,
				}
				select {
				case results <- result:
				case <-done:
					return
				}
			}
		}()
		return results
	}
	done := make(chan struct{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}

	// This simulates the caller of the goroutine - the parent process that
	// spawns the checkStatus (producer) goroutine.
	for result := range checkStatus(done, urls...) {
		if err := result.Err; err != nil {
			fmt.Printf("couldn't fetch url: %v\n", err)
			continue
		}
		fmt.Printf("response: %v\n", result.Resp.Status)
	}
}

// withErrorHandlingLimit is the same as the above function but only differs in
// the number of errors it will tolerate after which it breaks.
func withErrorHandlingLimit() {
	checkStatus := func(done <-chan struct{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{
					Resp: resp, Err: err,
				}
				select {
				case results <- result:
				case <-done:
					return
				}
			}
		}()
		return results
	}
	done := make(chan struct{})
	defer close(done)
	var errCount uint8

	urls := []string{"https://www.google.com", "a", "b", "c", "d"}

	// This simulates the caller of the goroutine - the parent process that
	// spawns the checkStatus (producer) goroutine.
	for result := range checkStatus(done, urls...) {
		if err := result.Err; err != nil {
			fmt.Printf("couldn't fetch url: %v\n", err)
			// here we immediately break as soon as the error limit of 3
			// is hit, which is a limiter the parent process can put to
			// signal that too many errors have occurred and the process
			// should be stopped.
			if errCount++; errCount >= 3 {
				fmt.Println("Too many errors, exiting!")
				break
			}
			continue
		}
		fmt.Printf("response: %v\n", result.Resp.Status)
	}
}
