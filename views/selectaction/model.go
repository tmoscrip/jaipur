package selectaction

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/models"
)

type SelectActionMenu struct {
	Game            *models.GameState
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

func New(game *models.GameState) SelectActionMenu {
	var selectedOptions = make([]int, 4)
	return SelectActionMenu{
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

func (v SelectActionMenu) Init() tea.Cmd {
	return nil
}

func (v SelectActionMenu) View() string {
	var s = ""

	for i, option := range v.options {
		var x = " "
		if i == *v.Cursor {
			x = "x"
		}
		s += fmt.Sprintf("[%s] %s\n", x, option)
	}
	return s
}

func (v SelectActionMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v SelectActionMenu) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
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
				return v, nil, "", err
			}
			return v, nil, v.transitions[*v.Cursor], nil
		}
	}
	return v, nil, "", nil
}

func (v SelectActionMenu) validate(cursor int) (bool, error) {
	switch cursor {
	// take one
	case 0:
		if len(v.Game.ActivePlayer().Hand) == 7 {
			return false, &TooManyInHandError{}
		}
		// take multiple
	case 1:
		return true, nil
		// take camels
	case 2:
		if len(v.Game.ActivePlayer().Hand)+v.Game.MarketCamelCount() > 7 {
			return false, &TooManyInHandError{}
		}
		if v.Game.MarketCamelCount() == 0 {
			return false, &models.NoCamelsInMarketError{}
		}
		// sell goods
	case 3:
		if len(v.Game.ActivePlayer().Hand) == 0 {
			return false, &NoGoodsInHandError{}
		}
		return true, nil

	default:
		return false, nil
	}
	return true, nil
}
