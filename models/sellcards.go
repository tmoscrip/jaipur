package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type SellCards struct {
	Game   *game.Game
	Cursor *int
}

func NewSellCards(game *game.Game) SellCards {
	g := game
	g.HandCursor = 0
	return SellCards{Game: g, Cursor: new(int)}
}

func (v SellCards) Init() tea.Cmd {
	return nil
}

func (v SellCards) View() string {
	s := tui.TitleStyle.Render("Sell cards")
	s += "\n"
	return s
}

func (v SellCards) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v SellCards) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	model := v
	if model.Game.HandCursor == -1 {
		model.Game.HandCursor = 0
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			model.Game.HandCursor = -1
			return model, nil, "selectActionMenu", nil
		}
		if msg.String() == "left" {
			if *model.Cursor > 0 {
				*model.Cursor = *model.Cursor - 1
				model.Game.HandCursor = *model.Cursor
			}
		}
		if msg.String() == "right" {
			if *model.Cursor < len(model.Game.Players.Active().Hand)-1 {
				*model.Cursor = *model.Cursor + 1
				model.Game.HandCursor = *model.Cursor
			}
		}
		if msg.String() == "enter" {
			model.Game.ToggleHand(*model.Cursor)
		}
		if msg.String() == "n" {
			endRound, err := model.Game.PlayerSellCards(model.Game.HandSelected)
			if err != nil {
				return model, nil, "", err
			}
			if endRound {
				return model, nil, "endRound", nil
			}
			return model, nil, "selectActionMenu", nil
		}
	}
	return model, nil, "", nil
}
