package tests

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"toggleDecks"
)

// Convert the deck into a string for comparison purposes.
func DeckToSSortedString(d toggleDecks.Deck) string {
	return SortDeckString(d.String())
}

// Convert a deck code string the same as a deck would be so that constant deck strings can be compared to decks.
func SortDeckString(expected string) string {
	sortedStrings := strings.Split(expected, " ")
	sort.Strings(sortedStrings)
	return strings.Join(sortedStrings, " ")
}

// Compare a deck to a string of deck codes to ensure that they contain the same cards, not necessarily in the same order.
func DeckContainsCards(d toggleDecks.Deck, expected string) (bool, string, string) {
	the_deck := DeckToSSortedString(d)
	expected_deck := SortDeckString(expected)
	return the_deck == expected_deck, expected_deck, the_deck
}

// Mocking support for GUID deck identifiers
type GuidMock struct{}

func (m GuidMock) GenerateIdentifier() string {
	return "a251071b-662f-44b6-ba11-e24863039c59"
}

// Patch the guid provider to return the same guid every time for testing purposes.
func PatchUID() {
	toggleDecks.TheGuidProvider = GuidMock{}
}

// Unpatch the guid provider so it goes back to providing random GUIDs.
func UnPatchUID() {
	toggleDecks.TheGuidProvider = toggleDecks.GuidIdProvider{}
}

// Empty the mock database.
func ClearTheDatabase() {
	toggleDecks.OurDecks = map[string]toggleDecks.Deck{}
}

// Create a custom test deck and file it into the proper place to mock having added it through the api.
func createStandardDeck() (iid string) {
	ClearTheDatabase()
	iid = "a251071b-662f-44b6-ba11-e24863039c59"
	toggleDecks.OurDecks[iid] = toggleDecks.CreateDeck("AS KH 8C")
	return
}

// Setup for creating a deck, and execute a deck creation request.
func DoCreateRequest(t *testing.T, method string, url string) (body string, result int) {
	ClearTheDatabase()
	PatchUID()
	defer UnPatchUID()

	return DoRequest(t, method, url)
}

// Execute a request and return the results.
func DoRequest(t *testing.T, method string, url string) (body string, result int) {
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
