package main

import (
	"fmt"
	"time"
)

/*
	func main() {
		ch := fanIn(boring("joe"), boring("ann"))
		for i := 0; i < 10; i++ {
			fmt.Println(<-ch)
		}
		fmt.Println("using fanInSelect")
		ch1 := fanInSelectTimeout(boring("joe"), boring("ann"))
		for i := 0; i < 10; i++ {
			if _, ok := <-ch1; ok {
				fmt.Println(<-ch1)
			} else {
				fmt.Println("channel closed")
			}
		}
		fmt.Println("both are boring; leaving main")
	}
*/
func fanInSelectTimeout(input1, input2 <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				ch <- s
			case s := <-input2:
				ch <- s
			case <-time.After(1 * time.Second):
				fmt.Println("timeout")
				close(ch)
				return
			}
		}
	}()
	return ch
}

func fanInSelect(input1, input2 <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				ch <- s
			case s := <-input2:
				ch <- s
			}
		}
	}()
	return ch
}

func fanIn(input1, input2 <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		for {
			ch <- <-input1
		}
	}()
	go func() {
		for {
			ch <- <-input2
		}
	}()
	return ch
}
