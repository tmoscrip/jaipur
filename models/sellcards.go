package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type MenuOption struct {
	Label    string
	Selected bool
	Index    int
}

func (m MenuOption) Format(activeCursor int) string {
	var cursor = " "
	if m.Index == activeCursor {
		cursor = ">"
	}
	var checked = " "
	if m.Selected {
		checked = "x"
	}
	return fmt.Sprintf("%s [%s] %s", cursor, checked, m.Label)
}

type SellCards struct {
	Game    *game.GameState
	options []MenuOption
	Cursor  *int
}

func NewSellCards(game *game.GameState) SellCards {
	options := make([]MenuOption, len(game.ActivePlayer().Hand))
	for i, card := range game.ActivePlayer().Hand {
		options[i] = MenuOption{Index: i, Label: card.String(), Selected: false}
	}
	return SellCards{Game: game, options: options, Cursor: new(int)}
}

func (v SellCards) Init() tea.Cmd {
	return nil
}

func (v SellCards) View() string {
	s := tui.TitleStyle.Render("Sell cards")
	s += "\n"

	for _, card := range v.options {
		s += fmt.Sprintf("%s\n", card.Format(*v.Cursor))
	}
	return s
}

func (v SellCards) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v SellCards) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			return v, nil, "selectActionMenu", nil
		}
		if msg.String() == "up" {
			if *v.Cursor > 0 {
				*v.Cursor = *v.Cursor - 1
			}
		}
		if msg.String() == "down" {
			if *v.Cursor < len(v.Game.ActivePlayer().Hand)-1 {
				*v.Cursor = *v.Cursor + 1
			}
		}
		if msg.String() == "enter" {
			option := v.options[*v.Cursor]
			option.Selected = !option.Selected
			v.options[*v.Cursor] = option
		}
		if msg.String() == "n" {
			selected := make([]int, 0)
			for i, option := range v.options {
				if option.Selected {
					selected = append(selected, i)
				}
			}
			endRound, err := v.Game.PlayerSellCards(selected)
			if err != nil {
				return v, nil, "", err
			}
			if endRound {
				return v, nil, "endRound", nil
			}
			return v, nil, "selectActionMenu", nil
		}
	}
	return v, nil, "", nil
}
