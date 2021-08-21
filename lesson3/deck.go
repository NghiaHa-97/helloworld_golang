package main

import "fmt"

type deck []string

func newCard() deck {
	card := deck{}

	cardNumbers := deck{"1", "2", "3", "4", "5", "6"}
	cardABCs := deck{"A", "B", "C", "D"}

	for _, cardN := range cardNumbers {
		for _, cardA := range cardABCs {
			card = append(card, cardN+" of "+cardA)
		}
	}

	return card

}

func (d deck) print() {
	for i, cards := range d {
		fmt.Println(i, cards)
	}
}
