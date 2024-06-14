package game

import (
	"fmt"
	"math/rand/v2"

	"github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/logger"
)

type Game struct {
	Deck            []ResourceType
	Discarded       []ResourceType
	BonusTokens     map[int][]int // cards required -> points awarded
	ResourceTokens  map[ResourceType][]int
	Market          []ResourceType
	Players         []Player
	ActivePlayerIdx *int
}

func (g *Game) WinningPlayer() *Player {
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

func NewGame() Game {
	var g = Game{}
	g.ActivePlayerIdx = new(int)
	g.Players = make([]Player, 2)
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
		var player = Player{}
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

func (g *Game) ActivePlayer() *Player {
	return &g.Players[*g.ActivePlayerIdx]
}

func (g *Game) MarketCamelCount() int {
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

func (g *Game) PlayerTakeOne(marketIndex int) (bool, error) {
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

func (g *Game) PlayerTakeCamels() (bool, error) {
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
func (g *Game) PlayerTakeMultiple(hand []int, market []int) (bool, error) {
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

func (g *Game) PlayerSellCards(indexes []int) (bool, error) {
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

func (g *Game) AddToDiscard(cards []ResourceType) {
	g.Discarded = append(g.Discarded, cards...)
}

func (g *Game) ShouldRoundEnd() bool {
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

func (g *Game) nextPlayer() bool {
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

func (g *Game) StartRound() {
	newGame := NewGame()
	players := newGame.Players
	players[0].Rounds = g.Players[0].Rounds
	players[1].Rounds = g.Players[1].Rounds
	newGame.Players = players
	*g = newGame
}

func (g *Game) nextResourceToken(rt ResourceType) int {
	if len(g.ResourceTokens[rt]) == 0 {
		return 0
	}

	var score = g.ResourceTokens[rt][0]
	g.ResourceTokens[rt] = g.ResourceTokens[rt][1:]
	return score
}

func (g *Game) nextBonusToken(cardsScored int) int {
	if len(g.BonusTokens[cardsScored]) == 0 {
		return 0
	}

	var score = g.BonusTokens[cardsScored][0]
	g.BonusTokens[cardsScored] = g.BonusTokens[cardsScored][1:]
	return score
}

type ResourceType int

const (
	Diamond ResourceType = iota
	Gold
	Silver
	Cloth
	Spice
	Leather
	Camel
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

type StyleShortString struct {
	Style       lipgloss.Style
	ShortString string
}

var resourceTypeStyles = map[ResourceType]StyleShortString{
	Diamond: {
		Style:       lipgloss.NewStyle().Foreground(Diamond.Color()),
		ShortString: "Dia",
	},
	Gold: {
		Style:       lipgloss.NewStyle().Foreground(Gold.Color()),
		ShortString: "Gld",
	},
	Silver: {
		Style:       lipgloss.NewStyle().Foreground(Silver.Color()),
		ShortString: "Slv",
	},
	Cloth: {
		Style:       lipgloss.NewStyle().Foreground(Cloth.Color()),
		ShortString: "Cth",
	},
	Spice: {
		Style:       lipgloss.NewStyle().Foreground(Spice.Color()),
		ShortString: "Spi",
	},
	Leather: {
		Style:       lipgloss.NewStyle().Foreground(Leather.Color()),
		ShortString: "Lth",
	},
	Camel: {
		Style:       lipgloss.NewStyle().Foreground(Camel.Color()),
		ShortString: "Cml",
	},
}

func (c ResourceType) Style() lipgloss.Style {
	return resourceTypeStyles[c].Style
}

func (c ResourceType) ShortString() string {
	return resourceTypeStyles[c].ShortString
}

func (c ResourceType) String() string {
	return c.Style().Render(c.ShortString())
}
