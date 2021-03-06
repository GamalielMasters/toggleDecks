package toggleDecks

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

/*
	Helper functions for implementing the RestAPI including the rest structures and handling functions.
*/

// Interface to an ID provider for our api objects
type RestIdProvider interface {
	GenerateIdentifier() string
}

// The actual ID provider that returns random GUIDs
type GuidIdProvider struct{}

// Implement the RestIdProvider interface
func (i GuidIdProvider) GenerateIdentifier() string {
	return uuid.New().String()
}

// The actual IID generator hook used to grab an ID for a new deck.  This is so we can mock it in tests.
var TheGuidProvider RestIdProvider = GuidIdProvider{}

// The object used to list decks.
type ListDeckMessage struct {
	Decks []RestDeckMessage `json:"decks"`
}

// The object representing the deck information.  This is used both when we are and are not returning the cards in the deck.
type RestDeckMessage struct {
	Id        string     `json:"deck_id"`
	Shuffled  *bool      `json:"shuffled,omitempty"`
	Remaining *int       `json:"remaining,omitempty"`
	Cards     []RestCard `json:"cards,omitempty"`
}

// Create a new RestDockMessage from the iid and *Deck.  It can include or exclude the actual cards.
func NewRestDeckMessage(iid string, deck *Deck, includeCards bool) (rdm RestDeckMessage) {
	remaining := deck.Len()
	shuffled := deck.Shuffled

	var cards []RestCard
	if includeCards {
		cards = NewRestDrawMessage(deck.Cards).Cards
	} else {
		cards = []RestCard{}
	}
	return RestDeckMessage{iid, &shuffled, &remaining, cards}
}

// Structure for JSON serialization of a Card.
type RestCard struct {
	Value string `json:"value"`
	Suite string `json:"suite"`
	Code  string `json:"code"`
}

// The object representing the draw of a number of cards.
type RestDrawMessage struct {
	Cards []RestCard `json:"cards"`
}

// Return a RestDrawMessage from  an array of Card structs.  Translation for JSON serialization.
func NewRestDrawMessage(cards []Card) RestDrawMessage {
	restCards := make([]RestCard, len(cards))
	for i, c := range cards {
		restCards[i] = RestCard{c.Rank(), c.Suite(), c.Code()}
	}
	return RestDrawMessage{restCards}
}

// Indicate success and write json data.
func WriteSuccess(w http.ResponseWriter, rm interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if e := json.NewEncoder(w).Encode(rm); e != nil {
		_ = log.Output(1, "Error encoding data to json"+e.Error())
	}
}

// Indicate a error and write an error message.
func WriteError(w http.ResponseWriter, errorCode int, message string) {
	w.WriteHeader(errorCode)
	_, _ = fmt.Fprint(w, message)

}
