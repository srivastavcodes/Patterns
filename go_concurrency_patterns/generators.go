package main

import "fmt"

func main() {
	ch := boring("boring")
	for i := 0; i < 5; i++ {
		fmt.Printf("%s\n", <-ch)
	}
	fmt.Println("the boring function exited")
}

func boring(msg string) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- fmt.Sprintf("ch wrote %s : %dth time", msg, i)
		}
	}()
	return ch
}
