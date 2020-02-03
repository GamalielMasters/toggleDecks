package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"toggleDecks"
)

// Create a custom test deck and file it into the proper place to mock having added it through the api.
func createStandardDeck() (iid string) {
	ClearTheDatabase()
	iid = "a251071b-662f-44b6-ba11-e24863039c59"
	toggleDecks.OurDecks[iid] = toggleDecks.CreateDeck("AS KH 8C")
	return
}

// Setup and deploy an Open request to the endpoint.
func DoOpenRequest(t *testing.T, method string, url string) (body string, result int) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	toggleDecks.Router.ServeHTTP(rr, req)
	result = rr.Code

	body = rr.Body.String()
	return
}

func TestOpenADeck(t *testing.T) {
	iid := createStandardDeck()
	actual, status := DoOpenRequest(t, "GET", fmt.Sprintf("/api/v1/decks/%v", iid))

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":3,"cards":[{"value":"ACE","suite":"SPADES","code":"AS"},{"value":"KING","suite":"HEARTS","code":"KH"},{"value":"8","suite":"CLUBS","code":"8C"}]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

func TestOpeNonExistentDeck(t *testing.T) {
	createStandardDeck()
	_, status := DoOpenRequest(t, "GET", "/api/v1/decks/f6d6ccf0-b740-459d-9e90-4b3869e1985c")

	if status != http.StatusNotFound {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusNotFound, status)
	}

}

func TestOpenDeckWithMalformedId(t *testing.T) {
	_, status := DoOpenRequest(t, "GET", "/api/v1/decks/v1/")

	if status != http.StatusNotFound {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusNotFound, status)
	}
}
