package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	quit := make(chan string)
	quitFn(boringQuit("joe"), quit)
	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-quit)
	}
	quit <- "now stop"
	fmt.Printf("quitFn says: %s\n", <-quit)
}

func boringQuit(msg string) <-chan string {
	ch := make(chan string)
	go func() {
		count := 0
		for {
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			count++
			ch <- fmt.Sprintf("%s wrote : %dth time", msg, count)
		}
	}()
	return ch
}

func quitFn(input <-chan string, quit chan string) {
	go func() {
		for {
			select {
			case s := <-input:
				quit <- s
			case <-quit:
				fmt.Println("received quit")
				// do cleanup
				quit <- "cleanup done"
			}
		}
	}()
}
