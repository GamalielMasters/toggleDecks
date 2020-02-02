package tests

import (
	"sort"
	"strings"
	"testing"
	"toggleDecks"
)

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
	deck := toggleDecks.CreateFullDeck()
	if deck.Len() != 52 {
		t.Errorf("Wanted %v cards in deck, but got %v instead.", 52, deck.Len())
	}
}

func TestDefaultDeckIsAFullDeck(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	equal, sortedDeck, expected := DeckContainsCards(deck, toggleDecks.STANDARD_DECK)
	if !equal {
		t.Errorf("Deck does not contain the proper cards.\n\tExpected: '%v'\n\tGot:      '%v'", expected, sortedDeck)
	}
}

func TestDefaultDeckForCorrectOrder(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	if deck.String() != toggleDecks.STANDARD_DECK {
		t.Errorf("Deck is not in the proper order.\n\tExpected: '%v'\n\tGot:      '%v'", toggleDecks.STANDARD_DECK, deck.String())
	}
}

func TestShuffledDeckIsNotInInitialOrder(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	deck.Shuffle()
	if deck.String() == toggleDecks.STANDARD_DECK {
		t.Errorf("After calling shuffle, the deck is still in un-shuffled order.")
	}
}

func TestDeckKnowsIfItHasBeenShuffled(t *testing.T) {
	unshuffled_deck := toggleDecks.CreateFullDeck()
	if unshuffled_deck.Shuffled {
		t.Error("Unshuffled Deck thinks it's shuffdled.")
	}

	shuffled_deck := toggleDecks.CreateFullDeck()
	shuffled_deck.Shuffle()
	if !shuffled_deck.Shuffled {
		t.Error("Shuffled Deck thinks it's unshuffled.")
	}
}

func TestCanCreatePartialDeck(t *testing.T) {
	deckDefinition := "AS KD AC 2C KH"
	partialDeck := toggleDecks.CreateDeck(deckDefinition)

	if partialDeck.Len() != 5 {
		t.Errorf("Partial Deck Not Correct Size.  Expected %v cards, but received %v cards instead.", 5, partialDeck.Len())
	}

	if DeckToSSortedString(partialDeck) != SortDeckString(deckDefinition) {
		t.Errorf("Partial deck does not contain the correct cards.\n\tExpected: %v\n\tGot     : %v", DeckToSSortedString(partialDeck), SortDeckString(deckDefinition))
	}
}
