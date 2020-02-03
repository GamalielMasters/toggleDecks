package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// Drawing from a deck returns just the cards as an array of card objects.  By default you get one card.
func TestDrawOneCardFromADeck(t *testing.T) {
	iid := app.NewDeck("", false)
	actual, status := DoRequest(t, "POST", fmt.Sprintf("/api/v1/decks/%v/draw", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"cards":[{"value":"ACE","suite":"SPADES","code":"AS"}]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// You can draw more cards by specifying the number in the cards parameter to the request.
func TestDrawThreeCardsFromADeck(t *testing.T) {
	iid := app.NewDeck("", false)
	actual, status := DoRequest(t, "POST", fmt.Sprintf("/api/v1/decks/%v/draw?count=3", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"cards":[{"value":"ACE","suite":"SPADES","code":"AS"},{"value":"2","suite":"SPADES","code":"2S"},{"value":"3","suite":"SPADES","code":"3S"}]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// If you try to draw more cards then are left in the deck, you get whatever was left, not the number you asked for.
func TestDrawTwoCardsFromADeckContainingOne(t *testing.T) {
	iid := app.NewDeck("AS", false)
	actual, status := DoRequest(t, "POST", fmt.Sprintf("/api/v1/decks/%v/draw?count=2", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"cards":[{"value":"ACE","suite":"SPADES","code":"AS"}]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// Trying to draw from an exhausted deck gives you back an empty array... no cards at all.
func TestDrawCardFromAnExhaustedDeck(t *testing.T) {
	iid := app.NewDeck("AS QH", false)
	_, status := DoRequest(t, "POST", fmt.Sprintf("/api/v1/decks/%v/draw?count=2", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	actual, status := DoRequest(t, "POST", fmt.Sprintf("/api/v1/decks/%v/draw?count=1", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"cards":[]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// If you provide anything but an integer for the number to draw, you get the default of one card.
func TestDrawNaNCardsFromADeck(t *testing.T) {
	iid := app.NewDeck("", false)
	actual, status := DoRequest(t, "POST", fmt.Sprintf("/api/v1/decks/%v/draw?count=NaN", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"cards":[{"value":"ACE","suite":"SPADES","code":"AS"}]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// Trying to draw cards from an invalid deck id gives you a big fat 404, good buddy!
func TestDrawCardsFromAnInvalidDeck(t *testing.T) {
	_, status := DoRequest(t, "POST", "/api/v1/decks/INVALID_ID/draw?count=1")

	if status != http.StatusNotFound {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}
}
