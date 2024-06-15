package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type TakeOne struct {
	Game       *game.Game
	Cursor     *int
	confirming *bool
}

func NewTakeOne(game *game.Game) TakeOne {
	return TakeOne{Game: game, Cursor: new(int), confirming: new(bool)}
}

func (v TakeOne) Init() tea.Cmd {
	return nil
}

func (v TakeOne) View() string {
	var s = tui.TitleStyle.Render("Take one card")
	s += "\n"
	confirm := ""
	if *v.confirming {
		confirm = fmt.Sprintf(" (confirm %s)", v.Game.Market[*v.Cursor])
		confirm += "\nb for back"
	}
	s += fmt.Sprintf("Select a good to take%s\n", confirm)

	for i, resource := range v.Game.Market {
		var x = " "
		if *v.Cursor == i {
			x = "x"
		}
		s += fmt.Sprintf("[%s] %s\n", x, resource)
	}

	return s
}

func (v TakeOne) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v TakeOne) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			if *v.confirming {
				*v.confirming = false
				return v, nil, "", nil
			}
			return v, nil, "selectActionMenu", nil
		}
		if msg.String() == "up" && !*v.confirming {
			if *v.Cursor > 0 {
				*v.Cursor = *v.Cursor - 1
			}
		}
		if msg.String() == "down" && !*v.confirming {
			if *v.Cursor < len(v.Game.Market)-1 {
				*v.Cursor = *v.Cursor + 1
			}
		}
		if msg.String() == "enter" {
			if !*v.confirming {
				*v.confirming = true
				return v, nil, "", nil
			}
			if *v.confirming {
				*v.confirming = false
				endRound, _ := v.Game.PlayerTakeOne(*v.Cursor)
				if endRound {
					return v, nil, "endRound", nil
				}
				return v, nil, "selectActionMenu", nil
			}
		}
	}
	return v, nil, "", nil
}
