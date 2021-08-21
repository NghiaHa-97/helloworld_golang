package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// var card string = "11111"

	var i int
	rand.Seed(time.Now().UnixNano())
	for i = 0; i < 10; i++ {
		random := rand.Intn(100)
		fmt.Println(random)
	}

	card := "Reusable"
	card = newCardL()

	fmt.Println(card)

	// slices := []string{}
	slices := deck{}
	slices = append(slices, "Diamonds")

	// for index, current := range slices {
	// 	fmt.Println(index, current)
	// }

	// cards := newCard()
	newCard := newCard()
	newCard.print()

	slices.print()

}

func newCardL() string {
	return "LLLLL"
}
