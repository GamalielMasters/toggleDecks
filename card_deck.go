/*
	Modeling for a standard deck of cards.

	A deck contains a number of cards (identified by a character code).  Decks are created in a standard order, and can
	be shuffled after creation.  They also know how many cards they have remaining, and if they are or are not shuffled.

	A deck maintains the current state of an un-played deck, meaning that When you draw from a deck, the cards that are
	returned are removed from the deck, and the remaining cards decrease appropriately.
*/

package toggleDecks

import (
	"math/rand"
	"strings"
	"time"
)

// The codes that represent a standard deck, in the initial order, which is by Ascending Rank with suits in the order
// Spades, Diamonds, Clubs, Hearts.  This can be used to generate a standard "french" deck of cards.
const STANDARD_DECK = "AS 2S 3S 4S 5S 6S 7S 8S 9S 10S JS QS KS AD 2D 3D 4D 5D 6D 7D 8D 9D 10D JD QD KD AC 2C 3C 4C 5C 6C 7C 8C 9C 10C JC QC KC AH 2H 3H 4H 5H 6H 7H 8H 9H 10H JH QH KH"

var SuiteMap = map[string]string{"S": "SPADES", "D": "DIAMONDS", "C": "CLUBS", "H": "HEARTS"}
var RankMap = map[string]string{"A": "ACE", "1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7", "8": "8", "9": "9", "10": "10", "J": "JACK", "Q": "QUEEN", "K": "KING"}

type Card string

func (c Card) String() string {
	return string(c)
}

func (c Card) Code() string {
	return c.String()
}

func (c Card) Rank() string {
	code := c.Code()
	return RankMap[code[:len(code)-1]]
}

func (c Card) Suite() string {
	code := c.Code()
	return SuiteMap[code[len(code)-1:]]
}

type Deck struct {
	Cards    []Card
	Shuffled bool
}

func (d Deck) Len() int {
	return len(d.Cards)
}

func (d *Deck) String() string {
	s := ""
	for _, c := range d.Cards {
		s += c.String() + " "
	}

	return s[:len(s)-1]
}

func (d *Deck) Swap(i, j int) {
	d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
}

func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(d.Len(), d.Swap)
	d.Shuffled = true
}

// Create a standard 52 card "French" Deck of playing cards.
func CreateFullDeck() Deck {
	return CreateDeck(STANDARD_DECK)
}

// Create a new card deck containing the specified cards.
// Does NOT check if the cards are "valid" card codes for any given type of deck, that should be done by the caller.
func CreateDeck(includedCards string) (cards Deck) {
	cardCodes := strings.Split(includedCards, " ")
	cards = Deck{make([]Card, len(cardCodes)), false}
	for idx, code := range cardCodes {
		cards.Cards[idx] = Card(code)
	}

	return
}
