package models

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type TakeMultiple struct {
	Game      *game.Game
	Cursor    *int
	activeRow *int
}

func NewTakeMultiple(game *game.Game) TakeMultiple {
	g := game
	g.HandCursor = 0
	g.MarketCursor = -1
	return TakeMultiple{Game: g, Cursor: new(int), activeRow: new(int)}
}

func (v TakeMultiple) Init() tea.Cmd {
	return nil
}

func (v TakeMultiple) View() string {
	s := tui.TitleStyle.Render("Take Multiple")
	s += "\n"
	s += tui.HelpStyle.Render("b = back, c = confirm")
	s += "\n"

	return s
}

func (v TakeMultiple) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v TakeMultiple) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	model := v
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			model.Game.HandCursor = -1
			model.Game.MarketCursor = -1
			return model, nil, "selectActionMenu", nil
		}
		if msg.String() == "left" {
			if *model.Cursor > 0 {
				*model.Cursor = *model.Cursor - 1
				if *model.activeRow == 1 {
					model.Game.MarketCursor = *model.Cursor
					model.Game.HandCursor = -1
				} else {
					model.Game.HandCursor = *model.Cursor
					model.Game.MarketCursor = -1
				}
			}
		}
		if msg.String() == "right" {
			if *model.Cursor < len(model.Game.Market)-1 {
				*model.Cursor = *model.Cursor + 1
				if *model.activeRow == 1 {
					model.Game.MarketCursor = *model.Cursor
					model.Game.HandCursor = -1
				} else {
					model.Game.HandCursor = *model.Cursor
					model.Game.MarketCursor = -1
				}
			}
		}
		if msg.String() == "up" {
			if *model.activeRow == 0 {
				*model.activeRow = 1
				model.Game.MarketCursor = *model.Cursor
				model.Game.HandCursor = -1
			}
		}
		if msg.String() == "down" {
			if *model.activeRow == 1 {
				*model.activeRow = 0
				model.Game.HandCursor = *model.Cursor
				model.Game.MarketCursor = -1
			}
		}
		if msg.String() == "enter" {
			if *model.activeRow == 1 {
				model.Game.ToggleMarket(*model.Cursor)
			} else {
				model.Game.ToggleHand(*model.Cursor)
			}
		}
		if msg.String() == "c" {
			model.Game.LastActionString = "took multiple"
			endRound, err := model.Game.PlayerTakeMultiple(model.Game.MarketSelected, model.Game.HandSelected)
			if err != nil {
				return model, nil, "", err
			}
			if endRound {
				return model, nil, "endRound", nil
			}
			return model, nil, "startTurn", nil
		}
	}
	return model, nil, "", nil
}
