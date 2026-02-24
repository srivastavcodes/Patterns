package main

import (
	"bytes"
	"fmt"
	"sync"
)

/*
NOTE:
	Confinement is the idea of ensuring information is only ever available
	from one concurrent process. When this is achieved, a concurrent program
	is implicitly safe and no synchronization is needed.
*/

// adhocConfinement means achieving confinement through the means of convention
// and best practices rather than anything enforcing it.
func adhocConfinement() {
	// data as a variable is available to the entire function to use, but as
	// convention would dictate (presumably), we only use it inside the loop
	// function and access its values through the handler channel rather than
	// using it directly assuming it had values in a real program.
	data := make([]int, 4)

	// loop uses data as part of the convention.
	loop := func(handler chan<- int) {
		defer close(handler)
		for i := range data {
			handler <- data[i]
		}
	}
	handler := make(chan int)
	go loop(handler)

	// we access the values of data not by looping over data var itself but
	// over what the loop function produced using data's elements.
	for v := range handler {
		fmt.Println(v)
	}
	// problem is: conventions are easy to break, by simple negligence under
	// a deadline or because of the lack of context, which might end up
	// corrupting other parts of the program;
	// in this case, if someone accesses data directly and changes it, loop
	// starts behaving unexpectedly.
}

// lexicalConfinement means using the languages scoping constraints as a tool to
// enforce confinement.
func lexicalConfinement() {
	// chanOwner function here confines the very scope of the results variable,
	// ensuring no one outside this scope will be able to write to it and only
	// exposing the read aspect of this channel.
	chanOwner := func() <-chan int {
		// here results var is equivalent to data, from which data will be read
		// later on, but the writing aspect of this variable is confined to the
		// closure goroutine defined below it.
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i < 5; i++ {
				results <- i
			}
		}()
		// the chanOwner function here returns only a read view of the channel.
		return results
	}
	// consumer expects a read-only channel and doesn't concern itself with the
	// original data's location.
	consumer := func(results <-chan int) {
		for r := range results {
			fmt.Printf("received %d\n", r)
		}
		fmt.Println("done receiving")
	}
	results := chanOwner()
	consumer(results)
	// but channels are inherently concurrent safe, so next is an example that is
	// safe for concurrent usage.
}

// lexicalConfinementBuf is to showcase how lexical confinement would work a data type
// that is not inherently concurrent safe.
func lexicalConfinementBuf() {
	// because print data does not close the data variable being created below
	// holding the originally produced data, we cannot do the wrong thing and
	// mutate something we weren't supposed to. We can only operate on the data
	// provided to us when calling printData because of lexical confinement
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()
		buf := new(bytes.Buffer)
		for _, b := range data {
			_, _ = fmt.Fprintf(buf, "%c", b)
		}
		fmt.Println(buf.String())
	}
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("goland")
	// printData is only ever aware of the part of "goland" that is passed to it
	// and nothing more, ensuring confinement without the use of mutexes.
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])
	wg.Wait()
}
