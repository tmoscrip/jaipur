package takeone

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/logger"
	"github.com/tmoscrip/jaipur/models"
)

type TakeOneCard struct {
	Game       *models.GameState
	Cursor     *int
	confirming *bool
}

func New(game *models.GameState) TakeOneCard {
	return TakeOneCard{Game: game, Cursor: new(int), confirming: new(bool)}
}

func (v TakeOneCard) Init() tea.Cmd {
	return nil
}

func (v TakeOneCard) View() string {
	var s = ""
	s += fmt.Sprintf("Player: %s\n", v.Game.ActivePlayer().Name)
	s += fmt.Sprintf("Market: %s\n", v.Game.Market)
	s += fmt.Sprintf("Your hand: %s\n\n", v.Game.ActivePlayer().Hand)
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

func (v TakeOneCard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v TakeOneCard) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		logger.Message(fmt.Sprintf("cursor: %d", *v.Cursor))
		logger.Message(fmt.Sprintf("key: %s", msg.String()))
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
				logger.Message(fmt.Sprintf("new cursor: %d", *v.Cursor))
			}
		}
		if msg.String() == "enter" {
			logger.Message(fmt.Sprintf("selected good: %s", v.Game.Market[*v.Cursor]))
			if !*v.confirming {
				logger.Message("confirming")
				*v.confirming = true
				return v, nil, "", nil
			}
			if *v.confirming {
				logger.Message("comfirmed, next player turn")
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
