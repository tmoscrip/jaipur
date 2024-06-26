package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type SelectAction struct {
	Game            *game.Game
	options         []string
	transitions     []string
	Cursor          *int
	SelectedOptions []int
}

type TooManyInHandError struct{}

func (e *TooManyInHandError) Error() string {
	return "You have too many cards in your hand to do that"
}

type NoGoodsInHandError struct{}

func (e *NoGoodsInHandError) Error() string {
	return "You have no goods in your hand to sell"
}

func NewSelectAction(game *game.Game) SelectAction {
	var selectedOptions = make([]int, 4)
	return SelectAction{
		Game: game,
		options: []string{
			"Take 1 resource",
			"Take multiple resources",
			"Take camels",
			"Sell goods",
		},
		transitions:     []string{"takeOneCard", "takeSeveralCards", "takeCamels", "sellCards"},
		SelectedOptions: selectedOptions,
		Cursor:          new(int),
	}
}

func (v SelectAction) Init() tea.Cmd {
	return nil
}

func (v SelectAction) Render() string {
	var s = ""
	for i, option := range v.options {
		style := tui.TitleStyle.Bold(false).BorderForeground(tui.Silver)
		if i == *v.Cursor {
			style = style.Foreground(tui.White).Background(tui.EerieBlack).Bold(true)
		}
		s += style.Render(option)
		if i != len(v.options)-1 {
			s += "\n"
		}
	}
	return s
}

func (v SelectAction) View() string {
	return v.Render()
}

func (v SelectAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v SelectAction) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	model := v
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "up" {
			if *v.Cursor > 0 {
				*v.Cursor = *v.Cursor - 1
			}
		}
		if msg.String() == "down" {
			if *v.Cursor < len(v.options)-1 {
				*v.Cursor = *v.Cursor + 1
			}
		}
		if msg.String() == "enter" {
			_, err := v.validate(*v.Cursor)
			if err != nil {
				return model, nil, "", err
			}
			return v, nil, v.transitions[*v.Cursor], nil
		}
	}
	return model, nil, "", nil
}

func (v SelectAction) validate(cursor int) (bool, error) {
	switch cursor {
	// take one
	case 0:
		if len(v.Game.Players.Active().Hand) == 7 {
			return false, &TooManyInHandError{}
		}
		// take multiple
	case 1:
		return true, nil
		// take camels
	case 2:
		if len(v.Game.Players.Active().Hand)+v.Game.Market.Count(game.Camel) > 7 {
			return false, &TooManyInHandError{}
		}
		if v.Game.Market.Count(game.Camel) == 0 {
			return false, game.GameError{Message: "There are no camels in the market to take"}
		}
		// sell goods
	case 3:
		if len(v.Game.Players.Active().Hand) == 0 {
			return false, &NoGoodsInHandError{}
		}
		return true, nil

	default:
		return false, nil
	}
	return true, nil
}
