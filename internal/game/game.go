package game

import (
	"fmt"
	"math/rand/v2"

	"github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/logger"
)

type GameState struct {
	Deck            []ResourceType
	Discarded       []ResourceType
	BonusTokens     map[int][]int // cards required -> points awarded
	ResourceTokens  map[ResourceType][]int
	Market          []ResourceType
	Players         []PlayerState
	ActivePlayerIdx *int
}

func (g *GameState) WinningPlayer() *PlayerState {
	if g.Players[0].Score == g.Players[1].Score {
		return nil
	}

	var winner = &g.Players[0]
	for i := 1; i < 2; i++ {
		if g.Players[i].Score > winner.Score {
			winner = &g.Players[i]
		}
	}
	return winner
}

func NewGame() GameState {
	var g = GameState{}
	g.ActivePlayerIdx = new(int)
	g.Players = make([]PlayerState, 2)
	g.Deck = make([]ResourceType, 0)
	g.Discarded = make([]ResourceType, 0)
	for i := 0; i < 6; i++ {
		g.Deck = append(g.Deck, Diamond)
	}
	for i := 0; i < 6; i++ {
		g.Deck = append(g.Deck, Gold)
	}
	for i := 0; i < 6; i++ {
		g.Deck = append(g.Deck, Silver)
	}
	for i := 0; i < 8; i++ {
		g.Deck = append(g.Deck, Cloth)
	}
	for i := 0; i < 8; i++ {
		g.Deck = append(g.Deck, Spice)
	}
	for i := 0; i < 10; i++ {
		g.Deck = append(g.Deck, Leather)
	}
	// 11 camels, 8 in deck and 3 in market
	for i := 0; i < 8; i++ {
		g.Deck = append(g.Deck, Camel)
	}
	// shuffle deck
	logger.Message("Shuffling deck")
	for i := range g.Deck {
		j := rand.IntN(i + 1)
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	}

	g.Market = make([]ResourceType, 0)
	for i := 0; i < 3; i++ {
		g.Market = append(g.Market, Camel)
	}

	g.Market = append(g.Market, g.Deck[0], g.Deck[1])
	g.Deck = g.Deck[2:]

	g.BonusTokens = make(map[int][]int)
	g.BonusTokens[3] = []int{3, 3, 2, 2, 2, 1, 1}
	g.BonusTokens[4] = []int{6, 6, 5, 5, 4, 4}
	g.BonusTokens[5] = []int{10, 10, 8, 8, 6}

	g.ResourceTokens = make(map[ResourceType][]int)
	g.ResourceTokens[Diamond] = []int{7, 7, 5, 5, 5}
	g.ResourceTokens[Gold] = []int{6, 6, 5, 5, 5}
	g.ResourceTokens[Silver] = []int{5, 5, 5, 5, 5}
	g.ResourceTokens[Cloth] = []int{5, 3, 3, 2, 2, 1, 1}
	g.ResourceTokens[Spice] = []int{5, 3, 3, 2, 2, 1, 1}
	g.ResourceTokens[Leather] = []int{4, 3, 2, 1, 1, 1, 1, 1, 1}

	for i := 0; i < 2; i++ {
		var player = PlayerState{}
		player.Hand = g.Deck[0:5]
		logger.Message(fmt.Sprintf("Player %d hand: %s", i, player.Hand))
		g.Deck = g.Deck[5:]
		var camels = player.MoveCamelsToHerd()
		// refill hand accounting for removed camels
		for j := 0; j < camels; j++ {
			player.Hand = append(player.Hand, g.Deck[0])
			g.Deck = g.Deck[1:]
		}
		g.Players[i] = player
	}

	return g
}

func (g *GameState) ActivePlayer() *PlayerState {
	return &g.Players[*g.ActivePlayerIdx]
}

func (g *GameState) MarketCamelCount() int {
	var count = 0
	for _, card := range g.Market {
		if card == Camel {
			count++
		}
	}
	return count
}

type TooManyInHandError struct{}

func (e *TooManyInHandError) Error() string {
	return "Your hand would have more than 7 cards"
}

func (g *GameState) PlayerTakeOne(marketIndex int) (bool, error) {
	g.ActivePlayer().Hand = append(g.ActivePlayer().Hand, g.Market[marketIndex])
	g.Market[marketIndex] = g.Deck[0]
	g.Deck = g.Deck[1:]
	if g.nextPlayer() {
		return true, nil
	}
	return false, nil
}

type NoCamelsInMarketError struct{}

func (e *NoCamelsInMarketError) Error() string {
	return "There are no camels in the market to take"
}

func (g *GameState) PlayerTakeCamels() (bool, error) {
	if g.MarketCamelCount() == 0 {
		return false, &NoCamelsInMarketError{}
	}
	if g.MarketCamelCount()+len(g.ActivePlayer().Hand) > 7 {
		return false, &TooManyInHandError{}
	}
	originalMarket := g.Market
	herd := g.ActivePlayer().Herd
	newMarket := make([]ResourceType, 0)
	for i := 0; i < len(originalMarket); i++ {
		if originalMarket[i] == Camel {
			herd++
		}
		if originalMarket[i] != Camel {
			newMarket = append(newMarket, originalMarket[i])
		}
	}
	g.ActivePlayer().Herd = herd
	g.Market = newMarket
	if g.nextPlayer() {
		return true, nil
	}
	return false, nil
}

type TakeMultipleCountMismatch struct{}

func (e *TakeMultipleCountMismatch) Error() string {
	return "You must take the same number of cards from your hand and the market"
}

/*
PlayerTakeMultiple takes cards from the player's hand and the market and swaps them.
Parameters are indexes of the cards to take from the hand and the market.
*/
func (g *GameState) PlayerTakeMultiple(hand []int, market []int) (bool, error) {
	if len(hand) != len(market) {
		return false, &TakeMultipleCountMismatch{}
	}

	newHand := g.ActivePlayer().Hand
	for i := 0; i < len(hand); i++ {
		handCard := g.ActivePlayer().Hand[hand[i]]
		marketCard := g.Market[market[i]]
		newHand[hand[i]] = marketCard
		g.Market[market[i]] = handCard
	}

	g.ActivePlayer().Hand = newHand
	if g.nextPlayer() {
		return true, nil
	}
	return false, nil
}

type MustSellTwoCardsError struct{}

func (e *MustSellTwoCardsError) Error() string {
	return "You must sell at least 2 cards for that resource"
}

type NotEnoughOfResourceError struct{}

func (e *NotEnoughOfResourceError) Error() string {
	return "You don't have enough of that resource"
}

type SellCardsMismatchedResourcesError struct{}

func (e *SellCardsMismatchedResourcesError) Error() string {
	return "You can only sell one type of resource at a time"
}

func (g *GameState) PlayerSellCards(indexes []int) (bool, error) {
	mismatch := false
	for i := 0; i < len(indexes)-1; i++ {
		if g.ActivePlayer().Hand[indexes[i]] != g.ActivePlayer().Hand[indexes[i+1]] {
			mismatch = true
			break
		}
	}

	if mismatch {
		return false, &SellCardsMismatchedResourcesError{}
	}

	rt := g.ActivePlayer().Hand[indexes[0]]

	if rt == Diamond || rt == Gold || rt == Silver {
		if len(indexes) < 2 {
			return false, &MustSellTwoCardsError{}
		}
	}
	removedResources := g.ActivePlayer().RemoveIndexesFromHand(indexes)
	g.AddToDiscard(removedResources)

	for i := 0; i < len(indexes); i++ {
		g.ActivePlayer().AddScore(g.nextResourceToken(rt))
	}
	g.ActivePlayer().AddScore(g.nextBonusToken(len(indexes)))
	if g.nextPlayer() {
		return true, nil
	}
	return false, nil
}

func (g *GameState) AddToDiscard(cards []ResourceType) {
	g.Discarded = append(g.Discarded, cards...)
}

func (g *GameState) ShouldRoundEnd() bool {
	depletedResourceTokens := 0
	for _, rt := range []ResourceType{Diamond, Gold, Silver, Cloth, Spice, Leather} {
		if len(g.ResourceTokens[rt]) == 0 {
			depletedResourceTokens++
		}
	}

	cardsRemaining := 0
	cardsRemaining += len(g.Deck)
	cardsRemaining += len(g.Market)

	return depletedResourceTokens >= 3 || cardsRemaining < 5 || len(g.Deck) == 0
}

func (g *GameState) nextPlayer() bool {
	g.ActivePlayer().MoveCamelsToHerd()
	newIdx := (*g.ActivePlayerIdx + 1) % 2
	*g.ActivePlayerIdx = newIdx

	// refill market
	newMarket := g.Market
	// if newmarket is less than 5, add cards from deck
	for len(newMarket) < 5 {
		newMarket = append(newMarket, g.Deck[0])
		g.Deck = g.Deck[1:]
	}
	g.Market = newMarket

	if g.ShouldRoundEnd() {
		// score camels, player with most gets 5 points
		player0Camels := g.Players[0].Herd
		player1Camels := g.Players[1].Herd
		if player0Camels > player1Camels {
			g.Players[0].AddScore(5)
		} else if player1Camels > player0Camels {
			g.Players[1].AddScore(5)
		}
		if g.WinningPlayer() != nil {
			g.WinningPlayer().Rounds++
		}
		return true
	}
	return false
}

func (g *GameState) StartRound() {
	newGame := NewGame()
	players := newGame.Players
	players[0].Rounds = g.Players[0].Rounds
	players[1].Rounds = g.Players[1].Rounds
	newGame.Players = players
	*g = newGame
}

func (g *GameState) nextResourceToken(rt ResourceType) int {
	if len(g.ResourceTokens[rt]) == 0 {
		return 0
	}

	var score = g.ResourceTokens[rt][0]
	g.ResourceTokens[rt] = g.ResourceTokens[rt][1:]
	return score
}

func (g *GameState) nextBonusToken(cardsScored int) int {
	if len(g.BonusTokens[cardsScored]) == 0 {
		return 0
	}

	var score = g.BonusTokens[cardsScored][0]
	g.BonusTokens[cardsScored] = g.BonusTokens[cardsScored][1:]
	return score
}

type PlayerState struct {
	Name   string
	Herd   int
	Hand   []ResourceType
	Score  int
	Rounds int
}

func (p *PlayerState) ResourcesInHand(rt ResourceType) int {
	var count = 0
	for _, card := range p.Hand {
		if card == rt {
			count++
		}
	}
	return count
}

func (p *PlayerState) MoveCamelsToHerd() int {
	camelIndexes := make([]int, 0)
	for j, card := range p.Hand {
		if card == Camel {
			camelIndexes = append(camelIndexes, j)
		}
	}
	var removed = p.RemoveIndexesFromHand(camelIndexes)
	p.Herd += len(removed)
	return len(removed)
}

func (p *PlayerState) RemoveIndexesFromHand(indexes []int) []ResourceType {
	originalHand := p.Hand
	newHand := make([]ResourceType, 0)
	removedResources := make([]ResourceType, 0)

	indexMap := make(map[int]struct{}, len(indexes))
	for _, idx := range indexes {
		indexMap[idx] = struct{}{}
	}

	for i, resource := range originalHand {
		if _, found := indexMap[i]; found {
			removedResources = append(removedResources, resource)
		} else {
			newHand = append(newHand, resource)
		}
	}

	p.Hand = newHand
	return removedResources
}

func (p *PlayerState) AddScore(score int) {
	p.Score = p.Score + score
}

type ResourceType int

const (
	Diamond ResourceType = 0
	Gold    ResourceType = 1
	Silver  ResourceType = 2
	Cloth   ResourceType = 3
	Spice   ResourceType = 4
	Leather ResourceType = 5
	Camel   ResourceType = 6
)

func (c ResourceType) Color() lipgloss.Color {
	switch c {
	case Diamond:
		return lipgloss.Color("#00FFFF")
	case Gold:
		return lipgloss.Color("#FFD700")
	case Silver:
		return lipgloss.Color("#C0C0C0")
	case Cloth:
		return lipgloss.Color("#FF00AA")
	case Spice:
		return lipgloss.Color("#FFA500")
	case Leather:
		return lipgloss.Color("#8B4513")
	case Camel:
		return lipgloss.Color("#FFD777")
	default:
		return lipgloss.Color("#FFFFFF")
	}
}

func (c ResourceType) String() string {
	switch c {
	case Diamond:
		return lipgloss.NewStyle().Foreground(c.Color()).Render("Dia")
	case Gold:
		return lipgloss.NewStyle().Foreground(c.Color()).Render("Gld")
	case Silver:
		return lipgloss.NewStyle().Foreground(c.Color()).Render("Slv")
	case Cloth:
		return lipgloss.NewStyle().Foreground(c.Color()).Render("Cth")
	case Spice:
		return lipgloss.NewStyle().Foreground(c.Color()).Render("Spi")
	case Leather:
		return lipgloss.NewStyle().Foreground(c.Color()).Render("Lth")
	case Camel:
		return lipgloss.NewStyle().Foreground(c.Color()).Render("Cml")
	}
	return ""
}
