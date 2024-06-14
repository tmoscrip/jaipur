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
