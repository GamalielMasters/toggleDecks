package toggleDecks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"net/http"
)

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
	Id        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
	Cards     []Card `json:"cards,omitempty"`
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

	rm := RestMessage{iid, deck.Shuffled, deck.Len(), []Card{}}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rm)
}
