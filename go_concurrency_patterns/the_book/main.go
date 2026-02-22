package main

import "fmt"

func main() {
	fmt.Println("starting a goroutine function")
	<-fixedGoroutineWriteLeak()
}
