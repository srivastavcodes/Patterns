package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Site struct {
	url string
}

type Result struct {
	url        string
	statusCode int
}

func crawl(wid int, jobs <-chan Site, results chan<- Result) {
	for site := range jobs {
		fmt.Println("worker id: ", wid)

		resp, err := http.Get(site.url)
		if err != nil {
			log.Println(err.Error())
		}
		if resp == nil {
			return
		}
		resp.Body.Close()
		results <- Result{
			url:        site.url,
			statusCode: resp.StatusCode,
		}
	}
}

func main() {
	jobs := make(chan Site, 3)
	result := make(chan Result, 3)

	for w := 1; w <= 3; w++ {
		go crawl(w, jobs, result)
	}
	fmt.Println("worker pool initialized")

	time.Sleep(3 * time.Second)

	urls := []string{
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.instagram.com",
		"https://tutorialedge.net/pricing/",
		"https://www.reddit.com",
	}
	for _, url := range urls {
		jobs <- Site{url: url}
	}
	close(jobs)

	time.Sleep(3 * time.Second)
	for i := 0; i < 5; i++ {
		r := <-result
		fmt.Printf("%+v\n", r)
	}
}
