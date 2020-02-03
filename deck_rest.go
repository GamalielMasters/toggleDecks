package toggleDecks

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
}

// Provides an ID for our api objects
type RestIdProvider interface {
	GenerateIdentifier() string
}

type GuidIdProvider struct{}

func (i GuidIdProvider) GenerateIdentifier() string {
	return uuid.New().String()
}

var TheGuidProvider RestIdProvider = GuidIdProvider{}

var OurDecks = map[string]Deck{}

type RestMessage struct {
	Id        string     `json:"deck_id"`
	Shuffled  bool       `json:"shuffled"`
	Remaining int        `json:"remaining"`
	Cards     []RestCard `json:"cards,omitempty"`
}

type RestCard struct {
	Value string `json:"value"`
	Suite string `json:"suite"`
	Code  string `json:"code"`
}

func RestCardFromCard(card Card) (restCard RestCard) {
	restCard.Code = card.Code()
	restCard.Suite = card.Suite()
	restCard.Value = card.Rank()

	return
}

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

	OurDecks[iid] = deck

	rm := RestMessage{iid, deck.Shuffled, deck.Len(), []RestCard{}}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if e := json.NewEncoder(w).Encode(rm); e != nil {
		_ = log.Output(1, "Error encoding data to json"+e.Error())
	}
}

func DeckOpenEndpoint(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)

	iid, ok := pathParams["deckId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "Deck ID Required")
		return
	}

	deck, ok := OurDecks[iid]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "ID %v is not a valid deck id.", iid)
		return
	}

	cards := make([]RestCard, deck.Len())
	for i, c := range deck.Cards {
		cards[i] = RestCardFromCard(c)
	}

	rm := RestMessage{iid, deck.Shuffled, deck.Len(), cards}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if e := json.NewEncoder(w).Encode(rm); e != nil {
		_ = log.Output(1, "Error encoding data to json"+e.Error())
	}
}
