package game

type Market []ResourceType

func NewMarket() Market {
	return Market{
		Camel, Camel, Camel,
	}
}
func (m *Market) Count(rt ResourceType) int {
	var count = 0
	for _, card := range *m {
		if card == rt {
			count++
		}
	}
	return count
}
