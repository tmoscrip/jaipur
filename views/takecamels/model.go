package takecamels

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/tui"
	"github.com/tmoscrip/jaipur/models"
)

type TakeCamels struct {
	Game *models.GameState
}

func New(game *models.GameState) TakeCamels {
	return TakeCamels{Game: game}
}

func (v TakeCamels) Init() tea.Cmd {
	return nil
}

func (v TakeCamels) View() string {
	var s = ""
	s += tui.TitleStyle.Render(fmt.Sprintf("Take %d camels?", v.Game.MarketCamelCount()))
	return s
}

func (v TakeCamels) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v TakeCamels) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			return v, nil, "selectActionMenu", nil
		}
		if msg.String() == "enter" {
			endRound, error := v.Game.PlayerTakeCamels()
			if error != nil {
				return v, nil, "", error
			}
			if endRound {
				return v, nil, "endRound", nil
			}
			return v, nil, "selectActionMenu", nil
		}
	}
	return v, nil, "", nil
}
