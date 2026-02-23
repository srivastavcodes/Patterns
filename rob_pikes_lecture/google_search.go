package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Search func(query string) string

func fakeSearch(kind string) Search {
	return func(query string) string {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return fmt.Sprintf("%s result for %q\n", kind, query)
	}
}

// using replication pattern

func Google2o(query string) []string {
	ch := make(chan string)
	var result []string

	go func() { ch <- First(query, Web, Image) }()
	go func() { ch <- First(query, Image, Video) }()
	go func() { ch <- First(query, Web, Image, Video) }()

	timeout := time.After(80 * time.Millisecond)

	for i := 0; i < 3; i++ {
		select {
		case r := <-ch:
			result = append(result, r)
		case <-timeout:
			fmt.Println("timeout")
		}
	}
	return result
}

// First returns the first value returned out of all the searches used
func First(query string, replicas ...Search) string {
	ch := make(chan string)
	replica := func(i int) { ch <- replicas[i](query) }
	for i := range replicas {
		go replica(i)
	}
	return <-ch
}

// using timeout pattern

func Google1o(query string) []string {
	ch := make(chan string)
	var res []string

	go func() { ch <- Web(query) }()
	go func() { ch <- Image(query) }()
	go func() { ch <- Video(query) }()

	timeout := time.After(80 * time.Millisecond)

	for i := 0; i < 3; i++ {
		select {
		case r := <-ch:
			res = append(res, r)
		case <-timeout:
			fmt.Println("timed out")
		}
	}
	return res
}

func main() {
	start := time.Now()
	fmt.Println(Google2o("golang"))
	fmt.Println(time.Since(start))
}
