package game

import (
	"math/rand/v2"

	"github.com/tmoscrip/jaipur/internal/logger"
)

type Deck struct {
	cards []ResourceType
}

var countInDeck = map[ResourceType]int{
	Diamond: 6,
	Gold:    6,
	Silver:  6,
	Cloth:   8,
	Spice:   8,
	Leather: 10,
	Camel:   8,
}

type ErrNotEnoughCards struct{ error }

func NewDeck() Deck {
	d := Deck{}
	for card, count := range countInDeck {
		for i := 0; i < count; i++ {
			d.cards = append(d.cards, card)
		}
	}
	// 11 camels, 8 in deck and 3 in market
	for i := 0; i < 8; i++ {
		d.cards = append(d.cards, Camel)
	}

	logger.Message("Shuffling deck")
	for i := range d.cards {
		j := rand.IntN(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}

	return d
}

func (d *Deck) Draw(n int) ([]ResourceType, error) {
	if n > len(d.cards) {
		return nil, ErrNotEnoughCards{}
	}
	cards := d.cards[:n]
	d.cards = d.cards[n:]
	return cards, nil
}

func (d *Deck) Length() int {
	return len(d.cards)
}
