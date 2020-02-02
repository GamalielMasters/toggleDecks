package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"toggleDecks"
)

type GuidMock struct{}

// Return the same guid every time for testing purposes.
func (m GuidMock) GenerateIdentifier() string {
	return "a251071b-662f-44b6-ba11-e24863039c59"
}

func PatchUID() {
	toggleDecks.TheGuidProvider = GuidMock{}
}

func UnPatchUID() {
	toggleDecks.TheGuidProvider = toggleDecks.GuidIdProvider{}
}

func ClearTheDatabase() {
	toggleDecks.OurDecks = map[string]toggleDecks.Deck{}
}

func DoRequest(t *testing.T, method string, url string) (body string, result int) {
	ClearTheDatabase()
	PatchUID()
	defer UnPatchUID()

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(toggleDecks.DeckCreateEndpoint)
	handler.ServeHTTP(rr, req)
	result = rr.Code

	body = rr.Body.String()
	return
}

func TestCreateNewDefaultDeck(t *testing.T) {
	actual, status := DoRequest(t, "POST", "/api/v1/decks")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":52}` + "\n"
	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	if len(toggleDecks.OurDecks) != 1 {
		t.Error("Deck was not entered in internal map.")
	}
}

func TestCreateNewShuffledDeck(t *testing.T) {
	actual, status := DoRequest(t, "POST", "/api/v1/decks?shuffle=true")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":true,"remaining":52}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}
}

func TestCreateNewCustomDeck(t *testing.T) {
	actual, status := DoRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,KH")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":5}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	expectedDeck := "AS KD AC 2C KH"
	generatedDeck := toggleDecks.OurDecks["a251071b-662f-44b6-ba11-e24863039c59"]
	actualDeck := generatedDeck.String()

	if actualDeck != expectedDeck {
		t.Errorf("Deck is not configured with the correct cards.\n\tExpected: %v\n\tGot:     %v", expectedDeck, actualDeck)
	}
}

func TestCreateNewCustomShuffledDeck(t *testing.T) {
	actual, status := DoRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,KH&shuffle=true")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":true,"remaining":5}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	expectedDeck := "AS KD AC 2C KH"
	generatedDeck := toggleDecks.OurDecks["a251071b-662f-44b6-ba11-e24863039c59"]
	actualDeck := generatedDeck.String()

	if actualDeck == expectedDeck {
		t.Error("Deck is not shuffled.")
	}

	if SortDeckString(actualDeck) != SortDeckString(expectedDeck) {
		t.Errorf("Deck is not configured with the correct cards.\n\tExpected: %v\n\tGot:     %v", expectedDeck, actualDeck)
	}
}

func TestCreateCustomDeckWithRepeatedCardsOk(t *testing.T) {
	actual, status := DoRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,KH,AS,KD,AC,2C,KH&shuffle=false")

	if status != http.StatusOK {
		t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, status)
	}

	expected := `{"deck_id":"a251071b-662f-44b6-ba11-e24863039c59","shuffled":false,"remaining":10}` + "\n"

	if expected != actual {
		t.Errorf("Wrong result returned.\n\tExpected : %v\n\tGot      : %v", expected, actual)
	}

	expectedDeck := "AS KD AC 2C KH AS KD AC 2C KH"
	generatedDeck := toggleDecks.OurDecks["a251071b-662f-44b6-ba11-e24863039c59"]
	actualDeck := generatedDeck.String()

	if actualDeck != expectedDeck {
		t.Errorf("Deck is not configured with the correct cards.\n\tExpected: %v\n\tGot:     %v", expectedDeck, actualDeck)
	}
}

func TestCreateCustomDeckWithInvalidCards(t *testing.T) {
	// Wrong Rank - 11 of Dimonds is not a legal card.
	_, status := DoRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C, 11D")

	if status != http.StatusBadRequest {
		t.Errorf("Did not respond with proper error to bad card rank.  Expected %v but got %v", http.StatusBadRequest, status)
	}

	//Invalid Suite - There is no suite "B"
	_, status = DoRequest(t, "POST", "/api/v1/decks?cards=AS,KD,AC,2C,9B")

	if status != http.StatusBadRequest {
		t.Errorf("Did not respond with proper error to bad card suite.  Expected %v but got %v", http.StatusBadRequest, status)
	}
}

func TestMultipleDecksGetDifferentIds(t *testing.T) {
	// Cannot use DoRequest here because it installs the GUID mock.
	ClearTheDatabase()
	req, err := http.NewRequest("POST", "/api/v1/decks", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(toggleDecks.DeckCreateEndpoint)

	for i := 0; i < 3; i++ {
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Recived wrong status code. Expected %v, got %v.", http.StatusOK, rr.Code)
		}
	}

	if len(toggleDecks.OurDecks) != 3 {
		t.Errorf("Not all decks were added.  Expected 3 but got %v\nThis might mean multiple decks were created with the same ID.", len(toggleDecks.OurDecks))
	}
}
