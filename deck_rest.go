package toggleDecks

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
	Rest API for the card deck.  Relatively independent of the actual deck functionality.  This exposes the deck
	as a rest API with the following endpoints.

	/api/v1/decks 						-> POST -- Creates a new deck and returns its salient details.
	/api/v1/decks						-> GET  -- Returns a list of decks currently in the system.
	/api/v1/decks/{id)					-> GET  -- Opens a deck, providing its details and the remaining cards in the deck.
	/api/v1/decks/{id}/draw?number=x	-> POST -- Draws x cards from the deck, returning them and removing them from the deck.

*/

var Router = mux.NewRouter()

func init() {
	Router.HandleFunc("/api/v1/decks", DeckCreateEndpoint).Methods("POST")
	Router.HandleFunc("/api/v1/decks/{deckId}", DeckOpenEndpoint).Methods("GET")
	Router.HandleFunc("/api/v1/decks/{deckId}/draw", DeckDrawEndpoint).Methods("POST")
}

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

// A fake database of decks.
var OurDecks = map[string]*Deck{}

// The object representing the deck information.  This is used both when we are and are not returning the cards in the deck.
type RestDeckMessage struct {
	Id        string     `json:"deck_id"`
	Shuffled  *bool      `json:"shuffled"`
	Remaining *int       `json:"remaining"`
	Cards     []RestCard `json:"cards,omitempty"`
}

// The object representing the draw of a number of cards.
type RestDrawMessage struct {
	Cards []RestCard `json:"cards"`
}

// Structure for JSON serialization of a Card.
type RestCard struct {
	Value string `json:"value"`
	Suite string `json:"suite"`
	Code  string `json:"code"`
}

// Return an array of RestCard structs from an array of Card structs.  Translation for JSON serialization.
func restCardsFromCards(cards []Card) []RestCard {
	restCards := make([]RestCard, len(cards))
	for i, c := range cards {
		restCards[i] = RestCard{c.Rank(), c.Suite(), c.Code()}
	}
	return restCards
}

// Indicate success and write json data.
func WriteSuccess(w http.ResponseWriter, rm interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if e := json.NewEncoder(w).Encode(rm); e != nil {
		_ = log.Output(1, "Error encoding data to json"+e.Error())
	}
}

// Retrieve a stored deck from our "database" based on the request parameter named "deckId".  If it doesn't work, write
// and error and return done=true.  Otherwise, return the IID and the deck* to the retrieved deck.
func GetDeckFromRequest(w http.ResponseWriter, r *http.Request) (iid string, deck *Deck, done bool) {
	pathParams := mux.Vars(r)

	iid, ok := pathParams["deckId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Deck ID Required")
		return "", nil, true
	}

	deck, ok = OurDecks[iid]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "ID %v is not a valid deck id.", iid)
		return "", nil, true
	}
	return iid, deck, false
}

// REST Endpoint for Creating a new deck
func DeckCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	iid := TheGuidProvider.GenerateIdentifier()

	query := r.URL.Query()
	shuffled := query.Get("shuffle")
	custom := query.Get("cards")

	var deck Deck

	if len(custom) != 0 {
		// Check if the cards are legal, which they are if the suite and ranks exist in the respective maps.
		// We don't care if there is more than one of each card, etc, just that the collection is of actual card ids.

		cardIds := strings.Split(custom, ",")
		for _, id := range cardIds {
			pos := len(id) - 1
			suite := id[pos:]
			rank := id[:pos]

			_, rankOk := RankMap[rank]
			_, suiteOk := SuiteMap[suite]

			if !(rankOk && suiteOk) {
				w.WriteHeader(http.StatusBadRequest)

				if !rankOk {
					_, _ = fmt.Fprintf(w, "%v is not a valid rank for a custom deck.", rank)
				}

				if !suiteOk {
					_, _ = fmt.Fprintf(w, "%v is not a valid suite for a custom deck.", suite)
				}

				return
			}
		}

		deck = CreateDeck(strings.Join(cardIds, " "))
	} else {
		deck = CreateFullDeck()
	}

	if shuffled == "true" {
		deck.Shuffle()
	}

	OurDecks[iid] = &deck

	remaining := deck.Len()
	rm := RestDeckMessage{iid, &deck.Shuffled, &remaining, []RestCard{}}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if e := json.NewEncoder(w).Encode(rm); e != nil {
		_ = log.Output(1, "Error encoding data to json"+e.Error())
	}
}

// REST endpoint for opening (i.e. listing) a deck.
func DeckOpenEndpoint(w http.ResponseWriter, r *http.Request) {
	iid, deck, done := GetDeckFromRequest(w, r)
	if done {
		return
	}

	numRemaining := deck.Len()
	rm := RestDeckMessage{iid, &deck.Shuffled, &numRemaining, restCardsFromCards(deck.Cards)}

	WriteSuccess(w, rm)
}

// REST endpoint for drawing cards from a deck.
func DeckDrawEndpoint(w http.ResponseWriter, r *http.Request) {
	_, deck, done := GetDeckFromRequest(w, r)
	if done {
		return
	}

	query := r.URL.Query()
	count, err := strconv.Atoi(query.Get("cards"))

	if err != nil {
		count = 1
	}

	rm := RestDrawMessage{restCardsFromCards(deck.Draw(count))}

	WriteSuccess(w, rm)
}
