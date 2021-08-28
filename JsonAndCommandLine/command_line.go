package main

import (
	"fmt"
	"os"
)

func main() {

	argsWithProg := os.Args
	argsWithoutProg := os.Args[0:]

	// arg := os.Args[3]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	// 	fmt.Println(arg)
}
