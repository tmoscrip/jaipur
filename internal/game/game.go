package game

type Game struct {
	Deck             Deck
	Discarded        []ResourceType
	BonusTokens      BonusTokens
	ResourceTokens   ResourceTokens
	Market           Market
	Players          Players
	MarketSelected   []int
	MarketCursor     int
	HandSelected     []int
	HandCursor       int
	LastActionString string
}

func NewGame() Game {
	var g = Game{}
	g.HandCursor = -1
	g.MarketCursor = -1
	g.Players = Players{}
	g.Deck = NewDeck()
	g.Discarded = make([]ResourceType, 0)

	g.Market = NewMarket()

	drawn, _ := g.Deck.Draw(2)
	g.Market = append(g.Market, drawn...)

	g.BonusTokens = newBonusTokens()
	g.ResourceTokens = newResourceTokens()

	for i := 0; i < 2; i++ {
		var player = Player{}
		drawn, _ := g.Deck.Draw(5)
		player.Hand = drawn
		player.MoveCamelsToHerd()
		g.Players.Add(player)
	}

	return g
}

func (g *Game) ToggleMarket(index int) {
	// if the index is already in the selected list, remove it
	for i := 0; i < len(g.MarketSelected); i++ {
		if g.MarketSelected[i] == index {
			g.MarketSelected = append(g.MarketSelected[:i], g.MarketSelected[i+1:]...)
			return
		}
	}
	// if the index is not in the selected list, add it
	g.MarketSelected = append(g.MarketSelected, index)
}

func (g *Game) ToggleHand(index int) {
	// if the index is already in the selected list, remove it
	for i := 0; i < len(g.HandSelected); i++ {
		if g.HandSelected[i] == index {
			g.HandSelected = append(g.HandSelected[:i], g.HandSelected[i+1:]...)
			return
		}
	}
	// if the index is not in the selected list, add it
	g.HandSelected = append(g.HandSelected, index)
}

type TooManyInHandError struct{}

func (e *TooManyInHandError) Error() string {
	return "Your hand would have more than 7 cards"
}

func (g *Game) PlayerTakeOne(marketIndex int) (bool, error) {
	g.Players.Active().Hand = append(g.Players.Active().Hand, g.Market[marketIndex])
	drawn, _ := g.Deck.Draw(1)
	g.Market[marketIndex] = drawn[0]
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
	if g.Market.Count(Camel) == 0 {
		return false, &NoCamelsInMarketError{}
	}
	if g.Market.Count(Camel)+len(g.Players.Active().Hand) > 7 {
		return false, &TooManyInHandError{}
	}
	originalMarket := g.Market
	herd := g.Players.Active().Herd
	newMarket := make([]ResourceType, 0)
	for i := 0; i < len(originalMarket); i++ {
		if originalMarket[i] == Camel {
			herd++
		}
		if originalMarket[i] != Camel {
			newMarket = append(newMarket, originalMarket[i])
		}
	}
	g.Players.Active().Herd = herd
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

	newHand := g.Players.Active().Hand
	for i := 0; i < len(hand); i++ {
		handCard := g.Players.Active().Hand[hand[i]]
		marketCard := g.Market[market[i]]
		newHand[hand[i]] = marketCard
		g.Market[market[i]] = handCard
	}

	g.Players.Active().Hand = newHand
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

type NoResourceSelectedError struct{}

func (e *NoResourceSelectedError) Error() string {
	return "You must select at least one resource to sell"
}

func (g *Game) PlayerSellCards(indexes []int) (bool, error) {
	if len(indexes) == 0 {
		return false, &NoResourceSelectedError{}
	}
	mismatch := false
	for i := 0; i < len(indexes)-1; i++ {
		if g.Players.Active().Hand[indexes[i]] != g.Players.Active().Hand[indexes[i+1]] {
			mismatch = true
			break
		}
	}

	if mismatch {
		return false, &SellCardsMismatchedResourcesError{}
	}

	rt := g.Players.Active().Hand[indexes[0]]

	if rt == Diamond || rt == Gold || rt == Silver {
		if len(indexes) < 2 {
			return false, &MustSellTwoCardsError{}
		}
	}
	removedResources := g.Players.Active().RemoveIndexesFromHand(indexes)
	g.AddToDiscard(removedResources)

	for i := 0; i < len(indexes); i++ {
		g.Players.Active().AddScore(g.ResourceTokens.Next(rt))
	}
	g.Players.Active().AddScore(g.BonusTokens.Next(len(indexes)))
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
	cardsRemaining += g.Deck.Length()
	cardsRemaining += len(g.Market)

	return depletedResourceTokens >= 3 || cardsRemaining < 5 || g.Deck.Length() == 0
}

func (g *Game) nextPlayer() bool {
	g.MarketCursor = -1
	g.HandCursor = -1
	g.MarketSelected = make([]int, 0)
	g.HandSelected = make([]int, 0)

	g.Players.Active().MoveCamelsToHerd()
	g.Players.Next()

	// refill market
	newMarket := g.Market
	// if newmarket is less than 5, add cards from deck
	for len(newMarket) < 5 {
		drawn, _ := g.Deck.Draw(1)
		newMarket = append(newMarket, drawn[0])
	}
	g.Market = newMarket

	if g.ShouldRoundEnd() {
		// score camels, player with most gets 5 points
		player0Camels := g.Players.Herd(0)
		player1Camels := g.Players.Herd(1)
		if player0Camels > player1Camels {
			g.Players.AddScore(0, 5)
		} else if player1Camels > player0Camels {
			g.Players.AddScore(1, 5)
		}
		if g.Players.HigestScoring() != nil {
			g.Players.HigestScoring().WonRound()
		}
		return true
	}
	return false
}

func (g *Game) StartRound() {
	newGame := NewGame()
	players := newGame.Players
	players.Get(0).Rounds = g.Players.Get(0).Rounds
	players.Get(1).Rounds = g.Players.Get(1).Rounds
	newGame.Players = players
	*g = newGame
}
