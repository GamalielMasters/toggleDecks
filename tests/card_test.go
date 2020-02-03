package tests

import (
	"testing"
	"toggleDecks"
)

// Card Tests

func TestCardKnowsItsRank(t *testing.T) {
	card := toggleDecks.Card("AC")
	card2 := toggleDecks.Card("10D")
	if card.Rank() != "ACE" {
		t.Errorf("Card reports wrong rank.  Expected 'ACE', got '%v'", card.Rank())
	}

	if card2.Rank() != "10" {
		t.Errorf("Card reports wrong rank.  Expected '10', got '%v'", card2.Rank())
	}
}

func TestCardKnowsItsSuite(t *testing.T) {
	card := toggleDecks.Card("AC")
	if card.Suite() != "CLUBS" {
		t.Errorf("Card reports wrong suite. Expected 'CLUBS', got '%v'", card.Suite())
	}
}
