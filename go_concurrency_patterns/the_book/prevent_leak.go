package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("starting a goroutine with done channel")
	fixingGoroutineLeak()
}

// goroutineLeak showcases an example of a goroutine leak and mentions the
// reason why.
func goroutineLeak() {
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

// fixingGoroutineLeak showcases an example of a proper exiting pattern and mentions
// the reason why.
func fixingGoroutineLeak() {
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
		time.Sleep(2 * time.Second)
		// we close the done channel here so that the for-select loop in the
		// above goroutine selects the done chan as available and exists out
		// of the loop.
		close(done)
		fmt.Println("closed done channel")
	}()
	<-completed
	fmt.Println("function completed; exiting...")
	return
}
