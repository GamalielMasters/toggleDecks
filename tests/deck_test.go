package tests

import (
	"github.com/GamalielMasters/toggleDecks"
	"reflect"
	"testing"
)

// Test Underlying deck functionality.

// By default a deck is 52 playing cards.
func TestDefaultDeckIs52Cards(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	if deck.Len() != 52 {
		t.Errorf("Wanted %v cards in deck, but got %v instead.", 52, deck.Len())
	}
}

// And they are the standard cards you would expect.
func TestDefaultDeckIsAFullDeck(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	equal, sortedDeck, expected := DeckContainsCards(deck, toggleDecks.STANDARD_DECK)
	if !equal {
		t.Errorf("Deck does not contain the proper cards.\n\tExpected: '%v'\n\tGot:      '%v'", expected, sortedDeck)
	}
}

// In the order that a new deck comes in.
func TestDefaultDeckForCorrectOrder(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	if deck.String() != toggleDecks.STANDARD_DECK {
		t.Errorf("Deck is not in the proper order.\n\tExpected: '%v'\n\tGot:      '%v'", toggleDecks.STANDARD_DECK, deck.String())
	}
}

// Unless you call Shuffle.
func TestShuffledDeckIsNotInInitialOrder(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	deck.Shuffle()
	if deck.String() == toggleDecks.STANDARD_DECK {
		t.Errorf("After calling shuffle, the deck is still in un-shuffled order.")
	}
}

// In which case, the deck remembers that it has been shuffled
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

// You can also create a partial deck with only the cards you want.
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

// Once you have a deck, you can draw cards.
func TestCanDrawOneCardFromTheDeck(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	cards := deck.Draw(1)

	if len(cards) != 1 {
		t.Errorf("Got the wrong number of cards.  Drew 1, got %v", len(cards))
	}

	if cards[0].Code() != "AS" {
		t.Errorf("Got Wrong Card.  Expected 'AS' got '%v'", cards[0])
	}
}

// Even more then one at a time.
func TestCanDrawMultipleCardsFromTheDeck(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	cards := deck.Draw(5)

	if len(cards) != 5 {
		t.Errorf("Got the wrong number of cards.  Drew 5, got %v", len(cards))
	}

	expected := []toggleDecks.Card{"AS", "2S", "3S", "4S", "5S"}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Got Wrong Cards.  Expected '%v' got '%v'", expected, cards)
	}
}

// Once you draw cards, the deck doesn't have them anymore.
func TestDrawnCardsAreNoLongerInTheDeck(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	deck.Draw(1)

	if deck.Cards[0] != "2S" {
		t.Error("Drawn card is still present on the deck")
	}
}

// Multiple draws return progressive cards.  No dealing from the bottom here.
func TestMultipleDrawsTakeConsecutiveCards(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	deck.Draw(15)
	var cards [3]toggleDecks.Card
	var expected = [3]toggleDecks.Card{"3D", "4D", "5D"}

	for i := 0; i < 3; i++ {
		cards[i] = deck.Draw(1)[0]
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Didn't get the expected cards from multiple draws.  Expected: %v Got: %v", expected, cards)
	}

}

// If you try to draw more cards then are left in the deck, it won't complain, but you only get what's left.
func TestCannotDrawMoreCardsThenAreLeftInTheDeck(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	deck.Draw(50)
	cards := deck.Draw(3)

	if len(cards) != 2 {
		t.Errorf("Should have gotten only 2 remaining cards, but got %v.", len(cards))
	}
}

// Even if you've drawn all the cards.
func TestCannotDrawCardsFromAnExhaustedDeck(t *testing.T) {
	deck := toggleDecks.CreateFullDeck()
	deck.Draw(52)

	if deck.Len() != 0 {
		t.Errorf("Deck should be exhausted after drawing 52 cards, but has %v cards left", deck.Len())
	}

	cards := deck.Draw(1)

	if len(cards) != 0 {
		t.Errorf("Should have gotten an empty slice from drawing from an empty deck, but got %v", cards)
	}

}
