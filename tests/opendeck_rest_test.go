package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestOpenADeck(t *testing.T) {
	iid := createCustomDeck("AS KH 8C")
	actual, status := DoRequest(t, "GET", fmt.Sprintf("/api/v1/decks/%v", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := fmt.Sprintf(`{"deck_id":"%v","shuffled":false,"remaining":3,"cards":[{"value":"ACE","suite":"SPADES","code":"AS"},{"value":"KING","suite":"HEARTS","code":"KH"},{"value":"8","suite":"CLUBS","code":"8C"}]}`+"\n", iid)
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

func TestOpeNonExistentDeck(t *testing.T) {
	createCustomDeck("AS KH 8C")
	_, status := DoRequest(t, "GET", "/api/v1/decks/f6d6ccf0-b740-459d-9e90-4b3869e1985c")

	if status != http.StatusNotFound {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusNotFound, status)
	}

}

func TestOpenDeckWithMalformedId(t *testing.T) {
	_, status := DoRequest(t, "GET", "/api/v1/decks/v1/")

	if status != http.StatusNotFound {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusNotFound, status)
	}
}
