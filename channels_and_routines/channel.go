package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
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
		go checkLink1(link, c)
	}
	// fmt.Println(<-c);

	// for{
	// 	go checkLink(<-c,c)
	// }

	for link := range c {
		//
		// go checkLink(link, c, 15*time.Second)

		//
		go func(l string) {
			time.Sleep(5 * time.Second)
			checkLink1(l, c)
		}(link)
	}
}

func checkLink(link string, c chan string, t time.Duration) {
	time.Sleep(t)
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "error")
		c <- link
		return
	}
	fmt.Println(link, " up!")
	c <- link

}

func checkLink1(link string, c chan string) {
	_, err := http.Get(link)
	if err != nil {
		fmt.Println(link, "error")
		c <- link
		return
	}
	fmt.Println(link, " up!")
	c <- link

}
