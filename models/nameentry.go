package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
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
	var pString = fmt.Sprintf("Player %d", v.Game.Players.ActiveIdx+1)
	return fmt.Sprintf("%s player's name:\n> %s", pString, v.Game.Players.Active().Name)
}

func (v NameEntry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v NameEntry) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	var p = v.Game.Players.Active()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			v.Game.Players.Next()
			if v.Game.Players.Get(1).Name != "" {
				return v, nil, "selectActionMenu", nil
			}
			return v, nil, "", nil
		}
		p.Name += msg.String()
	}
	return v, nil, "", nil
}
