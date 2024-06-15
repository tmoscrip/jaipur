package game

type Player struct {
	Name   string
	Herd   int
	Hand   []ResourceType
	Score  int
	Rounds int
}

func (p *Player) ResourcesInHand(rt ResourceType) int {
	var count = 0
	for _, card := range p.Hand {
		if card == rt {
			count++
		}
	}
	return count
}

func (p *Player) MoveCamelsToHerd() int {
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

func (p *Player) RemoveIndexesFromHand(indexes []int) []ResourceType {
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

func (p *Player) AddScore(score int) {
	p.Score = p.Score + score
}

func (p *Player) WonRound() {
	p.Rounds++
}

func (p *Player) setHand(hand []ResourceType) {
	p.Hand = hand
}

type Players struct {
	players   []Player
	ActiveIdx int
}

func (p *Players) Active() *Player {
	return &p.players[p.ActiveIdx]
}

func (p *Players) Get(i int) *Player {
	return &p.players[i]
}

func (p *Players) Next() *Player {
	p.ActiveIdx = (p.ActiveIdx + 1) % 2
	return &p.players[p.ActiveIdx]
}

func (p *Players) Add(player Player) {
	p.players = append(p.players, player)
	if len(p.players) > 2 {
		panic("Too many players")
	}
}

func (p *Players) Herd(i int) int {
	return p.players[i].Herd
}

// Returns the player with the highest score, or nil if the scores are tied.
func (p *Players) HigestScoring() *Player {
	if p.players[0].Score == p.players[1].Score {
		return nil
	}

	var winner = &p.players[0]
	for i := 1; i < 2; i++ {
		if p.players[i].Score > winner.Score {
			winner = &p.players[i]
		}
	}
	return winner
}

func (p *Players) AddScore(idx int, score int) {
	p.players[idx].Score += score
}
