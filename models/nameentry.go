package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
)

type NameEntry struct {
	Game *game.GameState
}

func NewNameEntry(game *game.GameState) NameEntry {
	return NameEntry{Game: game}
}

func (v NameEntry) Init() tea.Cmd {
	return nil
}

func (v NameEntry) View() string {
	var pString = fmt.Sprintf("Player %d", *v.Game.ActivePlayerIdx+1)
	return fmt.Sprintf("%s player's name:\n> %s", pString, v.Game.ActivePlayer().Name)
}

func (v NameEntry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v NameEntry) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	var p = v.Game.ActivePlayer()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			*v.Game.ActivePlayerIdx++
			if *v.Game.ActivePlayerIdx > 1 {
				*v.Game.ActivePlayerIdx = 0
				return v, nil, "selectActionMenu", nil
			}
			return v, nil, "", nil
		}
		p.Name += msg.String()
	}
	return v, nil, "", nil
}
