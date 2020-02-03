package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// When you visit the endpoint for decks with a get request, you get a list of all the decks in the system.
func TestListDeck(t *testing.T) {
	app.ClearTheDatabase()
	id := app.NewDeck("", true)

	actual, status := DoRequest(t, "GET", "/api/v1/decks")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := fmt.Sprintf(`{"decks":[{"deck_id":"%v"}]}`+"\n", id)
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// But if there aren't any, you get an empty list, of course.
func TestListNoDeck(t *testing.T) {
	app.ClearTheDatabase()
	actual, status := DoRequest(t, "GET", "/api/v1/decks")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"decks":[]}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}
