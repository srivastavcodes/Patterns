package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*func main() {
	joe := boring("joe")
	ann := boring("ann")

	for i := 0; i < 5; i++ {
		fmt.Printf("%s :: read at -> %d\n", <-joe, time.Now().Unix())
		fmt.Printf("%s :: read at -> %d\n", <-ann, time.Now().Unix())
	}
	fmt.Println("the boring function exited")
}*/

// it'll return a receive channel of string
func boring(msg string) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			ch <- fmt.Sprintf("%s wrote : %dth time :: at %d", msg, i, time.Now().Unix())
		}
	}()
	return ch
}
