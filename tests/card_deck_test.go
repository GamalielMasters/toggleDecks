package tests

import (
	"sort"
	"strings"
	"testing"
	"toggleDecks"
)

const STANDARD_DECK = "AS 2S 3S 4S 5S 6S 7S 8S 9S 10S JS QS KS AD 2D 3D 4D 5D 6D 7D 8D 9D 10D JD QD KD AC 2C 3C 4C 5C 6C 7C 8C 9C 10C JC QC KC AH 2H 3H 4H 5H 6H 7H 8H 9H 10H JH QH KH"

// Helper Functions

func DeckToSSortedString(d toggleDecks.Deck) string {
	return SortDeckString(d.String())
}

func SortDeckString(expected string) string {
	sortedStrings := strings.Split(expected, " ")
	sort.Strings(sortedStrings)
	return strings.Join(sortedStrings, " ")
}

func DeckContainsCards(d toggleDecks.Deck, expected string) (bool, string, string) {
	the_deck := DeckToSSortedString(d)
	expected_deck := SortDeckString(expected)
	return the_deck == expected_deck, expected_deck, the_deck
}

// Tests

func TestDefaultDeckIs52Cards(t *testing.T) {
	deck := toggleDecks.CreateDeck()
	if deck.Len() != 52 {
		t.Errorf("Wanted %v cards in deck, but got %v instead.", 52, deck.Len())
	}
}

func TestDefaultDeckIsAFullDeck(t *testing.T) {
	deck := toggleDecks.CreateDeck()
	equal, sorted_deck, expected := DeckContainsCards(deck, STANDARD_DECK)
	if !equal {
		t.Errorf("Deck does not contain the proper cards.\n\tExpected: '%v'\n\tGot:      '%v'", expected, sorted_deck)
	}
}

func TestDefaultDeckForCorrectOrder(t *testing.T) {
	deck := toggleDecks.CreateDeck()
	if deck.String() != STANDARD_DECK {
		t.Errorf("Deck is not in the proper order.\n\tExpected: '%v'\n\tGot:      '%v'", STANDARD_DECK, deck.String())
	}
}

func TestShuffledDeckIsNotInInitialOrder(t *testing.T) {
	deck := toggleDecks.CreateDeck()
	deck.Shuffle()
	if deck.String() == STANDARD_DECK {
		t.Errorf("Deck is not shuffled.")
	}
}

func TestDeckKnowsIfItHasBeenShuffled(t *testing.T) {
	unshuffled_deck := toggleDecks.CreateDeck()
	if unshuffled_deck.Shuffled {
		t.Error("Unshuffled Deck thinks it's shuffdled.")
	}

	shuffled_deck := toggleDecks.CreateDeck()
	shuffled_deck.Shuffle()
	if !shuffled_deck.Shuffled {
		t.Error("Shuffled Deck thinks it's unshuffled.")
	}
}
