package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// When you visit the endpoint for decks with a get request, you get a list of all the decks in the system.
func TestListDeck(t *testing.T) {
	app.ClearTheDatabase()
	iids := make([]string, 3)
	for i := 0; i < 3; i++ {
		app.NewDeck("", true)
	}

	x := 0
	for i, _ := range app.TheDecks {
		iids[x] = i
		x++
	}

	actual, status := DoRequest(t, "GET", "/api/v1/decks")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := fmt.Sprintf(`{"decks":[{"deck_id":"%v"},{"deck_id":"%v"},{"deck_id":"%v"}]}`+"\n", iids[0], iids[1], iids[2])
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
