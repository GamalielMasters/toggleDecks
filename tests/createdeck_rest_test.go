package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test the Create Deck Endpoint.

// A default deck is 52 un-shuffled cards
func TestCreateNewDefaultDeck(t *testing.T) {
	app.ClearTheDatabase()
	actual, status := DoCreateRequest(t, "POST", "/api/v1/decks")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":52}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	if len(app.TheDecks) != 1 {
		t.Error("Deck was not entered in internal map.")
	}
}

// But you can request that it be shuffled if you like.
func TestCreateNewShuffledDeck(t *testing.T) {
	actual, status := DoCreateRequest(t, "POST", "/api/v1/decks?shuffle=true")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":true,"remaining":52}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

// Or you can ask for it to contain only specific cards you select. You get a custom deck back in the order you provide the cards.
func TestCreateNewCustomDeck(t *testing.T) {
	actual, status := DoCreateRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,KH")
	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":5}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	expectedDeck := "AS KD AC 2C KH"
	generatedDeck, _ := app.GetDeck("a251071b-662f-44b6-ba11-e24863039c59")
	actualDeck := generatedDeck.String()

	if actualDeck != expectedDeck {
		t.Errorf("Deck is not configured with the correct cards.\n\tExpected: %v\n\tGot:     %v", expectedDeck, actualDeck)
	}
}

// Unless you request that it be shuffled.
func TestCreateNewCustomShuffledDeck(t *testing.T) {
	actual, status := DoCreateRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,KH&shuffle=true")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":true,"remaining":5}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	expectedDeck := "AS KD AC 2C KH"
	generatedDeck := app.TheDecks["a251071b-662f-44b6-ba11-e24863039c59"]
	actualDeck := generatedDeck.String()

	if actualDeck == expectedDeck {
		t.Error("Deck is not shuffled.")
	}

	if SortDeckString(actualDeck) != SortDeckString(expectedDeck) {
		t.Errorf("Deck is not configured with the correct cards.\n\tExpected: %v\n\tGot:     %v", expectedDeck, actualDeck)
	}
}

// It's fine to ask for more then one of the same card in your custom deck.
func TestCreateCustomDeckWithRepeatedCards(t *testing.T) {
	actual, status := DoCreateRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,KH,AS,KD,AC,2C,KH&shuffle=false")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":10}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	expectedDeck := "AS KD AC 2C KH AS KD AC 2C KH"
	generatedDeck := app.TheDecks["a251071b-662f-44b6-ba11-e24863039c59"]
	actualDeck := generatedDeck.String()

	if actualDeck != expectedDeck {
		t.Errorf("Deck is not configured with the correct cards.\n\tExpected: %v\n\tGot:     %v", expectedDeck, actualDeck)
	}
}

// But you'll get an error if you ask for things that aren't valid cards in the normal 52 card deck.
func TestCreateCustomDeckWithInvalidCards(t *testing.T) {
	// Wrong Rank - 11 of Dimonds is not a legal card.
	_, status := DoCreateRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C, 11D")

	if status != http.StatusBadRequest {
		t.Errorf("Did not respond with proper error to bad card rank.  Expected %v but got %v", http.StatusBadRequest, status)
	}

	//Invalid Suite - There is no suite "B"
	_, status = DoCreateRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,9B")

	if status != http.StatusBadRequest {
		t.Errorf("Did not respond with proper error to bad card suite.  Expected %v but got %v", http.StatusBadRequest, status)
	}
}

// If you ask for more then one deck, each one gets a different ID.
func TestMultipleDecksGetDifferentIds(t *testing.T) {
	// Cannot use DoCreateRequest here because it installs the GUID mock, which defeats what we are testing here.
	app.ClearTheDatabase()
	req, err := http.NewRequest("POST", "/api/v1/decks", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	for i := 0; i < 3; i++ {
		app.Router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, rr.Code)
		}
	}

	if len(app.TheDecks) != 3 {
		t.Errorf("Not all decks were added.  Expected 3 but got %v\nThis might mean multiple decks were created with the same ID.", len(app.TheDecks))
	}
}
