package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("starting a goroutine function")
	withoutErrorHandling()
	time.Sleep(1 * time.Second)
	withErrorHandling()
	time.Sleep(1 * time.Second)
	withErrorHandlingLimit()
	time.Sleep(1 * time.Second)
}
