package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/tui"
	"github.com/tmoscrip/jaipur/logger"
	"github.com/tmoscrip/jaipur/models"
	"github.com/tmoscrip/jaipur/views/endround"
	"github.com/tmoscrip/jaipur/views/nameentry"
	"github.com/tmoscrip/jaipur/views/selectaction"
	"github.com/tmoscrip/jaipur/views/sellcards"
	"github.com/tmoscrip/jaipur/views/takecamels"
	"github.com/tmoscrip/jaipur/views/takemultiple"
	"github.com/tmoscrip/jaipur/views/takeone"
)

type MyMainModel interface {
	tea.Model
	MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error)
}

type MainModel struct {
	ActiveView   MyMainModel
	Game         *models.GameState
	ErrorMessage string
	ShowTopMenu  bool
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.ErrorMessage = ""

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			logger.Message("quitting")
			return m, tea.Quit
		}
	}

	_, cmd, transition, err := m.ActiveView.MyUpdate(msg)
	if transition != "" {
		logger.Message(fmt.Sprintf("transition: %s", transition))
	}

	if err != nil {
		m.ErrorMessage = err.Error()
		return m, cmd
	}

	switch transition {
	case "selectActionMenu":
		m.ShowTopMenu = true
		m.ActiveView = selectaction.New(m.Game)
	case "takeOneCard":
		m.ShowTopMenu = true
		m.ActiveView = takeone.New(m.Game)
	case "takeSeveralCards":
		m.ShowTopMenu = true
		m.ActiveView = takemultiple.New(m.Game)
	case "takeCamels":
		m.ShowTopMenu = true
		m.ActiveView = takecamels.New(m.Game)
	case "sellCards":
		m.ShowTopMenu = true
		m.ActiveView = sellcards.New(m.Game)
	case "endRound":
		m.ShowTopMenu = false
		m.ActiveView = endround.New(m.Game)
	}
	return m, cmd
}

var numberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

func (m MainModel) FormatDiscarded() string {
	discarded := ""
	for _, resource := range []models.ResourceType{models.Diamond, models.Gold, models.Silver, models.Cloth, models.Spice, models.Leather, models.Camel} {
		count := 0
		for _, card := range m.Game.Discarded {
			if card == resource {
				count++
			}
		}
		if count != 0 {
			discarded += numberStyle.Render(fmt.Sprintf("%dx%s ", count, resource))
		}
	}
	return discarded
}

func (m MainModel) FormatRemainingTokens() string {
	remaining := ""
	for _, resource := range []models.ResourceType{models.Diamond, models.Gold, models.Silver, models.Cloth, models.Spice, models.Leather} {
		count := 0
		for _, token := range m.Game.ResourceTokens[resource] {
			if token != 0 {
				count++
			}
		}
		if count != 0 {
			remaining += numberStyle.Render(fmt.Sprintf("%dx%s ", count, resource))
		}
	}
	return remaining
}

func (m MainModel) View() string {
	s := ""
	borderStyle := lipgloss.NewStyle().Width(tui.Width).Align(lipgloss.Left).Padding(1, 3).Background(tui.EerieBlack).Foreground(tui.WhiteSmoke2).Border(lipgloss.RoundedBorder()).BorderBackground(tui.DavysGray).BorderForeground(tui.Night)
	if m.ShowTopMenu {
		s += fmt.Sprintf("Player %d: %s     Rounds: %d\n", *m.Game.ActivePlayerIdx+1, m.Game.ActivePlayer().Name, m.Game.ActivePlayer().Rounds)
		s += fmt.Sprintf("Score: %d    Herd: %d\n", m.Game.ActivePlayer().Score, m.Game.ActivePlayer().Herd)
		s += fmt.Sprintf("Market: %s\n", m.Game.Market)
		s += fmt.Sprintf("Discarded: %s\n", m.FormatDiscarded())
		s += fmt.Sprintf("Remaining tokens: %s\n", m.FormatRemainingTokens())
		s += fmt.Sprintf("Your hand: %s", m.Game.ActivePlayer().Hand)
	}
	if len(s) > 0 {
		s = borderStyle.Render(s)
		s += "\n"
	}

	viewStyle := tui.MenuOptionsStyle
	if len(m.ErrorMessage) > 0 {
		viewStyle = viewStyle.BorderBottom(false)
	}

	s += viewStyle.Render(m.ActiveView.View())
	if m.ErrorMessage != "" {
		s += "\n"
		s += tui.ErrorStyle.Render(m.ErrorMessage)
	}

	return s
}

func initialModel() tea.Model {
	var g = models.NewGame()
	var v = nameentry.New(&g)
	return MainModel{
		ActiveView:  v,
		Game:        &g,
		ShowTopMenu: false,
	}
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
