package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
	NOTE:
	If a goroutine is responsible for creating a goroutine, it is also
	responsible for ensuring that it can stop the goroutine.
*/

// goroutineReadLeak showcases an example of a goroutine leak and mentions the
// reason why.
func goroutineReadLeak() {
	doWork := func(strs <-chan string) <-chan any {
		completed := make(chan any)
		go func() {
			defer fmt.Println("goroutine completed")
			defer close(completed)
			for str := range strs {
				// do something worthwhile
				fmt.Println(str)
			}
		}()
		return completed
	}
	// since we've passed nil to the strs channel, the inner goroutine never
	// actually gets any work done and gets hung up in the for-range loop
	// forever, locking the goroutine in place - in memory.
	doWork(nil)
	// perhaps more work is done here
	fmt.Println("function complete")
}

// fixedGoroutineReadLeak showcases an example of a proper exiting pattern and mentions
// the reason why.
func fixedGoroutineReadLeak() {
	// as a convention done should be the first parameter.
	doWork := func(done <-chan struct{}, strs <-chan string) <-chan any {
		completed := make(chan any)
		go func() {
			defer fmt.Println("doWork goroutine exited")
			defer close(completed)
		outer:
			for {
				select {
				case s := <-strs:
					fmt.Println(s)
					// done case ensures that this goroutine exits when
					// this channel is closed, indicating a exit signal.
				case <-done:
					break outer
				}
			}
		}()
		return completed
	}
	// we create a done channel that will be used to exit out of the for-select
	// loop.
	done := make(chan struct{})

	// we receive the work the goroutine does while passing it the done channel
	completed := doWork(done, nil)

	go func() {
		// we close the done channel here so that the for-select loop in the
		// above goroutine selects the done chan as available and exists out
		// of the loop.
		time.AfterFunc(2*time.Second, func() {
			fmt.Println("cancelling do work goroutine")
			close(done)
		})
	}()
	<-completed
	fmt.Println("function completed; exiting...")
	return
}

// goroutineWriteLeak showcases a goroutine leak that occurs when a goroutine
// is created without a done channel and there are no reads from the channel,
// and you attempt to write.
func goroutineWriteLeak() {
	makeRandStream := func() <-chan float64 {
		randStream := make(chan float64)
		go func() {
			defer fmt.Println("randStream goroutine exited")
			defer close(randStream)
			for {
				// if no one reads from randStream, this
				// operation blocks.
				randStream <- rand.Float64()
			}
		}()
		return randStream
	}
	randStream := makeRandStream()

	fmt.Println("3 random numbers:")
	for i := 0; i < 4; i++ {
		// reads from randStream, but stops at 4 => no reads => writing
		// to it blocks in the makeRandStream goroutine.
		fmt.Println(<-randStream)
	}
}

// goroutineWriteLeak showcases a correct pattern for creating a writer goroutine.
// Where even if there are no reads from the channel, and you attempt to write -
// which would normally end up the goroutine getting hung up - can now exit safely
// using a done channel.
//
// See: fixedGoroutineReadLeak for inner workings.
func fixedGoroutineWriteLeak() <-chan struct{} {
	makeRandStream := func(done <-chan struct{}) <-chan float64 {
		randStream := make(chan float64)
		go func() {
			defer fmt.Println("randStream goroutine exited")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Float64():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}
	exit := make(chan struct{})
	done := make(chan struct{})
	randStream := makeRandStream(done)

	go func() {
		time.AfterFunc(2*time.Second, func() {
			fmt.Println("closing done channel")
			close(done)
			close(exit)
		})
	}()
	for i := 0; i < 4; i++ {
		fmt.Println(<-randStream)
	}
	return exit
}
