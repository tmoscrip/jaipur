package server

import (
	"fmt"
	"testing"

	"github.com/tmoscrip/jaipur/internal/game"
)

func TestBothPlayersCanSetName(t *testing.T) {
	commands := []GameCommandWrapper{}
	for i := 0; i < 2; i++ {
		commands = append(commands, GameCommandWrapper{
			Player: i,
			Command: GameCommand{
				Action: SetPlayerName,
				Params: map[string]interface{}{
					"name": fmt.Sprintf("test %d", i),
				},
			},
		})
	}

	game := game.NewGame()
	for _, command := range commands {
		newGame, err := command.Run(game)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		game = newGame
	}

	for i := 0; i < 2; i++ {
		if game.Players.Get(i).Name != fmt.Sprintf("test %d", i) {
			t.Errorf("Expected name to be 'test', got %v", game.Players.Get(i).Name)
		}
	}
}
