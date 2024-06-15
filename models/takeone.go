package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type TakeOne struct {
	Game   *game.Game
	Cursor *int
}

func NewTakeOne(game *game.Game) TakeOne {
	g := game
	g.MarketCursor = 0
	return TakeOne{Game: g, Cursor: new(int)}
}

func (v TakeOne) Init() tea.Cmd {
	return nil
}

func (v TakeOne) View() string {
	var s = tui.TitleStyle.Render("Take one card")
	s += "\n"
	if len(v.Game.MarketSelected) == 1 {
		s += tui.TitleStyle.Render("Take " + v.Game.Market[v.Game.MarketSelected[0]].String() + "?")
		s += "\n"
	}
	s += tui.HelpStyle.Render("b = back, enter = confirm")

	return s
}

func (v TakeOne) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v TakeOne) confirming() bool {
	return len(v.Game.MarketSelected) == 1
}

func (v TakeOne) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	model := v
	if model.Game.MarketCursor == -1 {
		model.Game.MarketCursor = 0
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			if model.confirming() {
				model.Game.ToggleMarket(model.Game.MarketSelected[0])
				return model, nil, "", nil
			}
			model.Game.MarketCursor = -1
			return model, nil, "selectActionMenu", nil
		}
		if msg.String() == "left" && !model.confirming() {
			if *model.Cursor > 0 {
				*model.Cursor = *model.Cursor - 1
				model.Game.MarketCursor = *model.Cursor
			}
		}
		if msg.String() == "right" && !model.confirming() {
			if *model.Cursor < len(model.Game.Market)-1 {
				*model.Cursor = *model.Cursor + 1
				model.Game.MarketCursor = *model.Cursor
			}
		}
		if msg.String() == "enter" {
			if len(model.Game.MarketSelected) == 0 {
				model.Game.ToggleMarket(*model.Cursor)
				return model, nil, "", nil
			}
			if len(model.Game.MarketSelected) == 1 {
				model.Game.LastActionString = "took " + model.Game.Market[*model.Cursor].String()
				endRound, _ := model.Game.PlayerTakeOne(*model.Cursor)
				if endRound {
					return model, nil, "endRound", nil
				}
				return model, nil, "startTurn", nil
			}
		}
	}
	return model, nil, "", nil
}
