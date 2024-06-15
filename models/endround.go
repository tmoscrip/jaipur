package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type EndRound struct {
	Game *game.Game
}

func NewEndRound(game *game.Game) EndRound {
	return EndRound{game}
}

func (v EndRound) Init() tea.Cmd {
	return nil
}

func (v EndRound) View() string {
	var s = tui.TitleStyle.Render("End of round!")
	s += "\n\n"
	s += "Scores:\n"
	// player 1
	p1 := v.Game.Players.Get(0)
	p2 := v.Game.Players.Get(1)
	s += fmt.Sprintf("%s: %d\n", p1.Name, p1.Score)
	// player 2
	s += fmt.Sprintf("%s: %d\n", p2.Name, p2.Score)

	s += fmt.Sprintf("\nWinner: %s\n\n", v.Game.Players.HigestScoring().Name)
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
