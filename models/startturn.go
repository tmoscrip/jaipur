package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type StartTurn struct {
	game *game.Game
}

func NewStartTurn(g *game.Game) StartTurn {
	return StartTurn{g}
}

func (v StartTurn) Init() tea.Cmd {
	return nil
}

func (v StartTurn) View() string {
	s := ""
	start := tui.TitleStyle.Render("Start turn " + v.game.Players.Active().Name)
	if v.game.LastActionString != "" {
		start += "\n" + tui.TitleStyle.Render(v.game.LastActionString)
	}

	s += start
	return s
}

func (v StartTurn) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v StartTurn) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	model := v
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			return model, nil, "selectActionMenu", nil
		}
	}
	return model, nil, "", nil
}
