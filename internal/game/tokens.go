package game

// cards required -> points awarded
type BonusTokens map[int][]int

func newBonusTokens() BonusTokens {
	return BonusTokens{
		3: {3, 3, 2, 2, 2, 1, 1},
		4: {6, 6, 5, 5, 4, 4},
		5: {10, 10, 8, 8, 6},
	}
}

func (b BonusTokens) Next(cardsScored int) int {
	if len(b[cardsScored]) == 0 {
		return 0
	}

	var score = b[cardsScored][0]
	b[cardsScored] = b[cardsScored][1:]
	return score
}

type ResourceTokens map[ResourceType][]int

func newResourceTokens() ResourceTokens {
	return ResourceTokens{
		Diamond: []int{7, 7, 5, 5, 5},
		Gold:    []int{6, 6, 5, 5, 5},
		Silver:  []int{5, 5, 5, 5, 5},
		Cloth:   []int{5, 3, 3, 2, 2, 1, 1},
		Spice:   []int{5, 3, 3, 2, 2, 1, 1},
		Leather: []int{4, 3, 2, 1, 1, 1, 1, 1, 1},
	}
}

func (r ResourceTokens) Next(rt ResourceType) int {
	if len(r[rt]) == 0 {
		return 0
	}

	var score = r[rt][0]
	r[rt] = r[rt][1:]
	return score
}
