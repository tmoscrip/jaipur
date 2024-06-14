package models

import (
	"reflect"
	"testing"
)

func TestMoveCamelsToHerd(t *testing.T) {
	tests := []struct {
		initialHand  []ResourceType
		initialHerd  int
		expectedHand []ResourceType
		expectedHerd int
	}{
		{
			initialHand:  []ResourceType{Camel, Cloth, Cloth, Cloth, Cloth},
			initialHerd:  2,
			expectedHand: []ResourceType{Cloth, Cloth, Cloth, Cloth},
			expectedHerd: 3,
		},
		{
			initialHand:  []ResourceType{Cloth, Cloth, Cloth, Cloth, Cloth},
			initialHerd:  2,
			expectedHand: []ResourceType{Cloth, Cloth, Cloth, Cloth, Cloth},
			expectedHerd: 2,
		},
		{
			initialHand:  []ResourceType{Camel, Camel, Camel, Camel, Camel},
			initialHerd:  2,
			expectedHand: []ResourceType{},
			expectedHerd: 7,
		},
	}

	for _, test := range tests {
		p := &PlayerState{Hand: test.initialHand, Herd: test.initialHerd}
		p.MoveCamelsToHerd()

		if !reflect.DeepEqual(p.Hand, test.expectedHand) {
			t.Errorf("expected hand %v, got %v", test.expectedHand, p.Hand)
		}
		if p.Herd != test.expectedHerd {
			t.Errorf("expected herd %d, got %d", test.expectedHerd, p.Herd)
		}
	}
}

func TestPlayerSellCards(t *testing.T) {
	tests := []struct {
		indexesToSell     []int
		initialHand       []ResourceType
		initialScore      int
		initialDiscarded  []ResourceType
		expectedHand      []ResourceType
		expectedScore     int
		expectedDiscarded []ResourceType
	}{
		{
			indexesToSell:     []int{0, 1, 2},
			initialHand:       []ResourceType{Diamond, Diamond, Diamond, Gold, Gold, Gold, Gold},
			initialScore:      0,
			initialDiscarded:  []ResourceType{},
			expectedHand:      []ResourceType{Gold, Gold, Gold, Gold},
			expectedScore:     7 + 7 + 5 + 3,
			expectedDiscarded: []ResourceType{Diamond, Diamond, Diamond},
		},
		// TODO: override resource tokens from default, sell more resources than tokens
		{
			indexesToSell:     []int{0, 1, 2},
			initialHand:       []ResourceType{Diamond, Diamond, Diamond},
			initialScore:      0,
			initialDiscarded:  []ResourceType{Spice},
			expectedHand:      []ResourceType{},
			expectedScore:     7 + 7 + 5 + 3,
			expectedDiscarded: []ResourceType{Spice, Diamond, Diamond, Diamond},
		},
	}

	for _, test := range tests {
		g := NewGame()

		player := g.Players[0]
		player.Hand = test.initialHand
		player.Score = test.initialScore

		g.Discarded = test.initialDiscarded

		g.Players[0] = player

		g.PlayerSellCards(test.indexesToSell)

		actualHand := g.Players[0].Hand
		if !reflect.DeepEqual(actualHand, test.expectedHand) {
			t.Errorf("expected hand %v, got %v", test.expectedHand, actualHand)
		}
		actualScore := g.Players[0].Score
		if actualScore != test.expectedScore {
			t.Errorf("expected score %d, got %d", test.expectedScore, actualScore)
		}

		actualDiscardedLen := len(g.Discarded)
		if actualDiscardedLen != len(test.expectedDiscarded) {
			t.Errorf("expected discarded len %d, got %d", len(test.expectedDiscarded), actualDiscardedLen)
		}
	}
}

func TestPlayerTakeMultiple(t *testing.T) {
	tests := []struct {
		giveFromHand   []int
		takeFromMarket []int
		initialHand    []ResourceType
		initialMarket  []ResourceType
		expectedHand   []ResourceType
		expectedMarket []ResourceType
	}{
		{
			giveFromHand:   []int{0, 1},
			takeFromMarket: []int{0, 1},
			initialHand:    []ResourceType{Gold, Gold},
			initialMarket:  []ResourceType{Diamond, Diamond},
			expectedHand:   []ResourceType{Diamond, Diamond},
			expectedMarket: []ResourceType{Gold, Gold, Spice, Spice, Spice},
		},
		{
			giveFromHand:   []int{0, 5},
			takeFromMarket: []int{0, 1},
			initialHand:    []ResourceType{Cloth, Gold, Gold, Gold, Gold, Spice},
			initialMarket:  []ResourceType{Diamond, Diamond, Leather},
			expectedHand:   []ResourceType{Diamond, Gold, Gold, Gold, Gold, Diamond},
			expectedMarket: []ResourceType{Cloth, Spice, Leather, Spice, Spice},
		},
	}

	for _, test := range tests {
		g := NewGame()

		player := g.Players[0]
		player.Hand = test.initialHand

		g.Players[0] = player
		g.Market = test.initialMarket

		g.PlayerTakeMultiple(test.giveFromHand, test.takeFromMarket)

		actualHand := g.Players[0].Hand
		if !reflect.DeepEqual(actualHand, test.expectedHand) {
			t.Errorf("expected hand %v, got %v", test.expectedHand, actualHand)
		}

		actualMarket := g.Market
		if len(actualMarket) != len(test.expectedMarket) {
			t.Errorf("expected market len %d, got %d", len(test.expectedMarket), len(actualMarket))
		}

	}
}

func TestPlayerTakeOne(t *testing.T) {
	tests := []struct {
		takeFromMarket    int
		initialHand       []ResourceType
		initialMarket     []ResourceType
		initialDeckCount  int
		expectedHand      []ResourceType
		expectedMarket    []ResourceType
		expectedDeckCount int
	}{
		{
			takeFromMarket:   0,
			initialHand:      []ResourceType{Gold},
			initialMarket:    []ResourceType{Diamond, Diamond, Spice, Spice, Spice},
			initialDeckCount: 30,
			expectedHand:     []ResourceType{Gold, Diamond},
			// actually should be random card added but we just check the length
			// also should never be less than 5 cards in market
			expectedMarket:    []ResourceType{Diamond, Diamond, Spice, Spice, Spice},
			expectedDeckCount: 29,
		},
	}

	for _, test := range tests {
		g := NewGame()

		player := g.Players[0]
		player.Hand = test.initialHand

		g.Players[0] = player
		g.Market = test.initialMarket

		g.Deck = make([]ResourceType, test.initialDeckCount)

		g.PlayerTakeOne(test.takeFromMarket)

		actualHand := g.Players[0].Hand
		if !reflect.DeepEqual(actualHand, test.expectedHand) {
			t.Errorf("expected hand %v, got %v", test.expectedHand, actualHand)
		}
		actualMarket := g.Market
		if len(actualMarket) != len(test.expectedMarket) {
			t.Errorf("expected market len %d, got %d", len(test.expectedMarket), len(actualMarket))
		}

		actualDeckCount := len(g.Deck)
		if actualDeckCount != test.expectedDeckCount {
			t.Errorf("expected deck count %d, got %d", test.expectedDeckCount, actualDeckCount)
		}
	}
}

func TestPlayerTakeCamels(t *testing.T) {
	tests := []struct {
		initialHand    []ResourceType
		initialMarket  []ResourceType
		initialHerd    int
		expectedHand   []ResourceType
		expectedMarket []ResourceType
		expectedHerd   int
	}{
		{
			initialHand:    []ResourceType{Gold},
			initialMarket:  []ResourceType{Camel, Camel, Spice, Spice, Spice},
			initialHerd:    0,
			expectedHand:   []ResourceType{Gold},
			expectedMarket: []ResourceType{Diamond, Diamond, Spice, Spice, Spice},
			expectedHerd:   2,
		},
	}

	for _, test := range tests {
		g := NewGame()

		player := g.Players[0]
		player.Hand = test.initialHand
		player.Herd = test.initialHerd

		g.Players[0] = player
		g.Market = test.initialMarket

		g.PlayerTakeCamels()

		actualHand := g.Players[0].Hand
		if !reflect.DeepEqual(actualHand, test.expectedHand) {
			t.Errorf("expected hand %v, got %v", test.expectedHand, actualHand)
		}

		actualMarket := g.Market
		if len(actualMarket) != len(test.expectedMarket) {
			t.Errorf("expected market len %d, got %d", len(test.expectedMarket), len(actualMarket))
		}

		actualHerd := g.Players[0].Herd
		if actualHerd != test.expectedHerd {
			t.Errorf("expected herd %d, got %d", test.expectedHerd, actualHerd)
		}
	}
}

func TestAddToDiscard(t *testing.T) {
	tests := []struct {
		initialDiscarded  []ResourceType
		toDiscard         []ResourceType
		expectedDiscarded []ResourceType
	}{
		{
			initialDiscarded:  []ResourceType{Gold, Gold, Gold},
			toDiscard:         []ResourceType{Diamond},
			expectedDiscarded: []ResourceType{Gold, Gold, Gold, Diamond},
		},
	}

	for _, test := range tests {
		g := NewGame()
		g.Discarded = test.initialDiscarded

		g.AddToDiscard(test.toDiscard)

		actualDiscarded := g.Discarded
		if !reflect.DeepEqual(actualDiscarded, test.expectedDiscarded) {
			t.Errorf("expected discarded %v, got %v", test.expectedDiscarded, actualDiscarded)
		}
	}
}
