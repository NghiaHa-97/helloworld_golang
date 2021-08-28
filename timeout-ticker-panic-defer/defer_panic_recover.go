package main

import (
	"fmt"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			println("Recover")
		}
	}()

	name := "ABCDEF"

	for _, val := range []rune(name) {
		defer fmt.Println(string(val))
	}

	panic("ERROR1")

}
