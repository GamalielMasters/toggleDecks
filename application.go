/*
	A card deck. Provides both a standard deck of 52 cards as well as custom decks made up of selected cards.
	Exposes the deck as a rest API with the following endpoints.

	/api/v1/decks 						-> POST -- Creates a new deck and returns its salient details.
	/api/v1/decks						-> GET  -- Returns a list of decks currently in the system.
	/api/v1/decks/{id)					-> GET  -- Opens a deck, providing its details and the remaining cards in the deck.
	/api/v1/decks/{id}/draw?number=x	-> POST -- Draws x cards from the deck, returning them and removing them from the deck.
*/

package toggleDecks

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Main application of the toggleDecks server.  Initializes the database and router, and optionally starts the server.
type App struct {
	Router *mux.Router

	// In memory storage for all created decks.
	TheDecks map[string]*Deck
}

// Create and initialize a new app (and database and router)
func NewApp() *App {
	a := App{mux.NewRouter(), map[string]*Deck{}}
	a.Router.HandleFunc("/api/v1/decks", a.DeckCreateEndpoint).Methods("POST")
	a.Router.HandleFunc("/api/v1/decks/{deckId}", a.DeckOpenEndpoint).Methods("GET")
	a.Router.HandleFunc("/api/v1/decks/{deckId}/draw", a.DeckDrawEndpoint).Methods("POST")

	return &a
}

// Run the server on the passed address.
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Empty the mock database.
func (a *App) ClearTheDatabase() {
	a.TheDecks = map[string]*Deck{}
}

// Create a deck and file it the decks database.
func (a *App) NewDeck(cards string, shuffle bool) (iid string) {
	iid = TheGuidProvider.GenerateIdentifier()
	var deck Deck
	if len(cards) == 0 {
		deck = CreateFullDeck()
	} else {
		deck = CreateDeck(cards)
	}

	if shuffle {
		deck.Shuffle()
	}

	a.TheDecks[iid] = &deck

	return
}

// Fetch a deck by it's ID.
func (a *App) GetDeck(iid string) (deck *Deck, ok bool) {
	deck, ok = a.TheDecks[iid]
	return
}

// Retrieve a stored deck from our "database" based on the request parameter named "deckId".  If it doesn't work, write
// and error and return done=true.  Otherwise, return the IID and the deck* to the retrieved deck.
// This is meant to be called as a helper from REST endpoints, it is not an endpoint itself.
func (a *App) getDeckFromRequest(w http.ResponseWriter, r *http.Request) (iid string, deck *Deck, err error) {
	pathParams := mux.Vars(r)

	iid, ok := pathParams["deckId"]
	if !ok {
		WriteError(w, http.StatusBadRequest, "Deck ID Required")
		return "", nil, fmt.Errorf("no Deck ID")
	}

	deck, ok = a.GetDeck(iid)
	if !ok {
		WriteError(w, http.StatusNotFound, fmt.Sprintf("ID %v is not a valid deck id.", iid))
		return "", nil, fmt.Errorf("deck ID does not reference a deck")
	}
	return iid, deck, nil
}

// REST Endpoint for Creating a new deck
func (a *App) DeckCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	shuffled := query.Get("shuffle") == "true"
	custom := query.Get("cards")

	var iid string

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
				WriteError(w, http.StatusBadRequest, "Invalid Card Identifier.")

				if !rankOk {
					_, _ = fmt.Fprintf(w, "%v is not a valid rank for a custom deck.", rank)
				}

				if !suiteOk {
					_, _ = fmt.Fprintf(w, "%v is not a valid suite for a custom deck.", suite)
				}

				return
			}
		}

		iid = a.NewDeck(strings.Join(cardIds, " "), shuffled)
	} else {
		iid = a.NewDeck("", shuffled)
	}
	deck, _ := a.GetDeck(iid)

	WriteSuccess(w, NewRestDeckMessage(iid, deck, false))
}

// REST endpoint for opening (i.e. listing) a deck.
func (a *App) DeckOpenEndpoint(w http.ResponseWriter, r *http.Request) {
	iid, deck, err := a.getDeckFromRequest(w, r)
	if err != nil {
		return
	}

	WriteSuccess(w, NewRestDeckMessage(iid, deck, true))
}

// REST endpoint for drawing cards from a deck.
func (a *App) DeckDrawEndpoint(w http.ResponseWriter, r *http.Request) {
	_, deck, err := a.getDeckFromRequest(w, r)
	if err != nil {
		return
	}

	query := r.URL.Query()
	count, err := strconv.Atoi(query.Get("count"))
	if err != nil {
		count = 1
	}

	WriteSuccess(w, NewRestDrawMessage(deck.Draw(count)))
}
