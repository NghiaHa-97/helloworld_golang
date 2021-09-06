package main

import (
	"fmt"
	"time"
)

func main()  {
	n:=0
	go func() {
		fmt.Println("GO")
	}()
	for {
		fmt.Println(n)
		n++
		time.Sleep(5*time.Second)
	}

}
