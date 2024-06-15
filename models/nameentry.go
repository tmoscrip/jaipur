package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type NameEntry struct {
	Game *game.Game
}

func NewNameEntry(game *game.Game) NameEntry {
	return NameEntry{Game: game}
}

func (v NameEntry) Init() tea.Cmd {
	return nil
}

func (v NameEntry) View() string {
	var s = fmt.Sprintf("Player %d name\n", v.Game.Players.ActiveIdx+1)
	s += tui.TitleStyle.Width(30).AlignHorizontal(lipgloss.Center).Render(v.Game.Players.Active().Name)
	return s
}

func (v NameEntry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v NameEntry) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	model := v
	var p = v.Game.Players.Active()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if v.Game.Players.Active().Name == "" {
				return model, nil, "", nil
			}
			v.Game.Players.Next()
			if v.Game.Players.Get(1).Name != "" {
				return model, nil, "startTurn", nil
			}
			return model, nil, "", nil
		}
		p.Name += msg.String()
	}
	return model, nil, "", nil
}
