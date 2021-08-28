package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {

	var waitGr = sync.WaitGroup{}

	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}

	c := make(chan string)

	for _, link := range links {
		// go checkLink(link, c, 0*time.Second)
		waitGr.Add(1)
		go checkLink(link, c, &waitGr)
	}

	waitGr.Wait()

}

func checkLink(link string, c chan string, wait *sync.WaitGroup) {

	defer wait.Done()

	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "error")
		c <- link
		return
	}
	fmt.Println(link, " up!")
	c <- link

}
