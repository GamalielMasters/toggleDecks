package toggleDecks

import (
	"math/rand"
	"time"
)

/*

 */

var SUITES = [4]string{"S", "D", "C", "H"}
var RANKS = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func LegalCodes() (cards [52]string) {
	idx := 0

	for _, suite := range SUITES {
		for _, rank := range RANKS {
			cards[idx] = rank + suite
			idx++
		}
	}

	return
}

type Card string

func (c Card) String() string {
	return string(c)
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

func CreateDeck() (cards Deck) {
	cards = Deck{make([]Card, 52), false}
	for idx, code := range LegalCodes() {
		cards.Cards[idx] = Card(code)
	}

	return
}
