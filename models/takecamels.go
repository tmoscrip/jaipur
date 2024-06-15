package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type TakeCamels struct {
	Game *game.Game
}

func NewTakeCamels(g *game.Game) TakeCamels {
	g2 := g
	// iterate over market, set all camels selected in market
	for i := 0; i < len(g2.Market); i++ {
		if g2.Market[i] == game.Camel {
			g2.MarketSelected = append(g2.MarketSelected, i)
		}
	}
	return TakeCamels{Game: g2}
}

func (v TakeCamels) Init() tea.Cmd {
	return nil
}

func (v TakeCamels) View() string {
	var s = ""
	s += tui.TitleStyle.Render(fmt.Sprintf("Take %d camels?", v.Game.Market.Count(game.Camel)))
	s += "\n"
	s += tui.HelpStyle.Render("b = back, enter = confirm")
	return s
}

func (v TakeCamels) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v TakeCamels) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	model := v
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			model.Game.MarketSelected = []int{}
			return model, nil, "selectActionMenu", nil
		}
		if msg.String() == "enter" {
			model.Game.LastActionString = "took camels"
			endRound, error := v.Game.PlayerTakeCamels()
			if error != nil {
				return model, nil, "", error
			}
			if endRound {
				return model, nil, "endRound", nil
			}
			return model, nil, "startTurn", nil
		}
	}
	return model, nil, "", nil
}
