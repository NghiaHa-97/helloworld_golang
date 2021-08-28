package main

import "testing"

func TestNewDeck(t *testing.T) {
	card := newCard()
	if len(card) != 24 {
		t.Errorf("Expected deck length of 10 %v", len(card))
	}
}
