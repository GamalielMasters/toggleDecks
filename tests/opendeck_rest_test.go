package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// When you open a deck, you get all the information about it, including all the cards it has in it.
func TestOpenADeck(t *testing.T) {
	iid := app.NewDeck("AS KH 8C", false)

	actual, status := DoRequest(t, "GET", fmt.Sprintf("/api/v1/decks/%v", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := fmt.Sprintf(`{"deck_id":"%v","shuffled":false,"remaining":3,"cards":[{"value":"ACE","suite":"SPADES","code":"AS"},{"value":"KING","suite":"HEARTS","code":"KH"},{"value":"8","suite":"CLUBS","code":"8C"}]}`+"\n", iid)
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// But that doesn't exhaust the deck like drawing from it does... the cards are all still there ready to be drawn.
func TestOpenDoesNotConsumeCards(t *testing.T) {
	iid := app.NewDeck("AS KH 8C", false)

	_, status := DoRequest(t, "GET", fmt.Sprintf("/api/v1/decks/%v", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	actual, status := DoRequest(t, "POST", fmt.Sprintf("/api/v1/decks/%v/draw?count=3", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"cards":[{"value":"ACE","suite":"SPADES","code":"AS"},{"value":"KING","suite":"HEARTS","code":"KH"},{"value":"8","suite":"CLUBS","code":"8C"}]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// Except of course for decks that don't exist.
func TestOpeNonExistentDeck(t *testing.T) {
	app.NewDeck("AS KH 8C", false)
	_, status := DoRequest(t, "GET", "/api/v1/decks/f6d6ccf0-b740-459d-9e90-4b3869e1985c")

	if status != http.StatusNotFound {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusNotFound, status)
	}

}

// And if you try to open a deck with something that isn't a deck id.... well, that's a no-no too.
func TestOpenDeckWithMalformedId(t *testing.T) {
	_, status := DoRequest(t, "GET", "/api/v1/decks/v1/")

	if status != http.StatusNotFound {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusNotFound, status)
	}
}
