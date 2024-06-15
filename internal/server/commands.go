package server

import "github.com/tmoscrip/jaipur/internal/game"

type Action string

const (
	SetPlayerName    Action = "setPlayerName"
	TakeOneCard      Action = "takeOneCard"
	TakeSeveralCards Action = "takeSeveralCards"
	TakeCamels       Action = "takeCamels"
	SellCards        Action = "sellCards"
)

type GameCommand struct {
	Action Action
	Params map[string]interface{}
}

type GameCommandWrapper struct {
	Player  int
	GameId  int
	Command GameCommand
}

func NewSetPlayerNameCommand(player int, name string) GameCommandWrapper {
	return GameCommandWrapper{
		Player: player,
		Command: GameCommand{
			Action: SetPlayerName,
			Params: map[string]interface{}{
				"name": name,
			},
		},
	}
}

func NewTakeOneCardCommand(player int, marketIndex int) GameCommandWrapper {
	return GameCommandWrapper{
		Player: player,
		Command: GameCommand{
			Action: TakeOneCard,
			Params: map[string]interface{}{
				"marketIndex": marketIndex,
			},
		},
	}
}

func NewTakeSeveralCardsCommand(player int, handIndexes []int, marketIndexes []int) GameCommandWrapper {
	return GameCommandWrapper{
		Player: player,
		Command: GameCommand{
			Action: TakeSeveralCards,
			Params: map[string]interface{}{
				"handIndexes":   handIndexes,
				"marketIndexes": marketIndexes,
			},
		},
	}
}

func NewTakeCamelsCommand(player int) GameCommandWrapper {
	return GameCommandWrapper{
		Player: player,
		Command: GameCommand{
			Action: TakeCamels,
			Params: map[string]interface{}{},
		},
	}
}

func NewSellCardsCommand(player int, handIndexes []int) GameCommandWrapper {
	return GameCommandWrapper{
		Player: player,
		Command: GameCommand{
			Action: SellCards,
			Params: map[string]interface{}{
				"handIndexes": handIndexes,
			},
		},
	}
}

func (gcw GameCommandWrapper) Run(g game.Game) (game.Game, error) {
	params := gcw.Command.Params
	if gcw.Player != g.Players.ActiveIdx {
		return g, game.GameError{Message: "Not this player's turn"}
	}

	switch gcw.Command.Action {
	case SetPlayerName:
		_, err := g.SetPlayerName(params["name"].(string))
		if err != nil {
			return g, err
		}

	case TakeOneCard:
		_, err := g.PlayerTakeOne(params["marketIndex"].(int))
		if err != nil {
			return g, err
		}

	case TakeSeveralCards:
		_, err := g.PlayerTakeMultiple(params["handIndexes"].([]int), params["marketIndexes"].([]int))
		if err != nil {
			return g, err
		}

	case TakeCamels:
		_, err := g.PlayerTakeCamels()
		if err != nil {
			return g, err
		}
	case SellCards:
		_, err := g.PlayerSellCards(params["handIndexes"].([]int))
		if err != nil {
			return g, err
		}

	default:
		return g, nil
	}

	return g, nil
}
