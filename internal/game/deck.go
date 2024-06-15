package game

import (
	"math/rand/v2"
)

type Deck struct {
	Cards []ResourceType
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
			d.Cards = append(d.Cards, card)
		}
	}
	// 11 camels, 8 in deck and 3 in market
	for i := 0; i < 8; i++ {
		d.Cards = append(d.Cards, Camel)
	}

	for i := range d.Cards {
		j := rand.IntN(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}

	return d
}

func (d *Deck) Draw(n int) ([]ResourceType, error) {
	if n > len(d.Cards) {
		return nil, ErrNotEnoughCards{}
	}
	cards := d.Cards[:n]
	d.Cards = d.Cards[n:]
	return cards, nil
}

func (d *Deck) Length() int {
	return len(d.Cards)
}
