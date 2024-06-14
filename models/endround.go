package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type EndRound struct {
	Game *game.GameState
}

func NewEndRound(game *game.GameState) EndRound {
	return EndRound{game}
}

func (v EndRound) Init() tea.Cmd {
	return nil
}

func (v EndRound) View() string {
	var s = tui.TitleStyle.Render("End of round!")
	s += "\n\n"
	s += "Scores:\n"
	for _, player := range v.Game.Players {
		s += fmt.Sprintf("%s: %d\n", player.Name, player.Score)
	}

	s += fmt.Sprintf("\nWinner: %s\n\n", v.Game.WinningPlayer().Name)

	s += "Press enter to start the next round"
	return s
}

func (v EndRound) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v EndRound) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			v.Game.StartRound()
			return v, nil, "selectActionMenu", nil
		}
	}
	return v, nil, "", nil
}
