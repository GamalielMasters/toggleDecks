package toggleDecks

import (
	"math/rand"
	"time"
)

/*

 */
var SUITES = [4]string{"S", "D", "C", "H"}
var RANKS = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

type Card string
type Deck []Card

func (c Card) String() string {
	return string(c)
}

func (d Deck) String() string {
	s := ""
	for _, c := range d {
		s += c.String() + " "
	}

	return s[:len(s)-1]
}

func (d Deck) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Deck) Shuffle() Deck {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	//return Shuffle{r.Perm( len(d) ), &d}
	r.Shuffle(len(d), d.Swap)
	return d
}

func CreateDeck() (cards Deck) {
	cards = Deck(make([]Card, 52))
	idx := 0
	for _, suite := range SUITES {
		for _, rank := range RANKS {
			cards[idx] = Card(rank + suite)
			idx++
		}
	}

	return
}
