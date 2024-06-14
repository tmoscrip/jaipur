package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/logger"
	"github.com/tmoscrip/jaipur/internal/tui"
	"github.com/tmoscrip/jaipur/models"
)

type MyMainModel interface {
	tea.Model
	MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error)
}

type MainModel struct {
	ActiveView   MyMainModel
	Game         *game.GameState
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
		m.ActiveView = models.NewSelectAction(m.Game)
	case "takeOneCard":
		m.ShowTopMenu = true
		m.ActiveView = models.NewTakeOne(m.Game)
	case "takeSeveralCards":
		m.ShowTopMenu = true
		m.ActiveView = models.NewTakeMultiple(m.Game)
	case "takeCamels":
		m.ShowTopMenu = true
		m.ActiveView = models.NewTakeCamels(m.Game)
	case "sellCards":
		m.ShowTopMenu = true
		m.ActiveView = models.NewSellCards(m.Game)
	case "endRound":
		m.ShowTopMenu = false
		m.ActiveView = models.NewEndRound(m.Game)
	}
	return m, cmd
}

var numberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

func (m MainModel) formatDiscarded() string {
	discarded := ""
	for _, resource := range []game.ResourceType{game.Diamond, game.Gold, game.Silver, game.Cloth, game.Spice, game.Leather, game.Camel} {
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

func (m MainModel) formatRemainingTokensColumn() string {
	s := ""
	for _, resource := range []game.ResourceType{game.Diamond, game.Gold, game.Silver, game.Cloth, game.Spice, game.Leather} {
		count := 0
		for _, token := range m.Game.ResourceTokens[resource] {
			if token != 0 {
				count++
			}
		}
		if count != 0 {
			var bars = ""
			for i := 0; i < count; i++ {
				bars += "|"
			}
			s += fmt.Sprintf("%s %s\n", bars, resource.String())
		}
	}
	return s
}

func (m MainModel) formatDiscardedColumn() string {
	s := ""
	for _, resource := range []game.ResourceType{game.Diamond, game.Gold, game.Silver, game.Cloth, game.Spice, game.Leather, game.Camel} {
		count := 0
		for _, card := range m.Game.Discarded {
			if card == resource {
				count++
			}
		}
		if count != 0 {
			var bars = ""
			for i := 0; i < count; i++ {
				bars += "|"
			}
			s += fmt.Sprintf("%s (%d) %s\n", bars, count, resource.String())
		}
	}
	return s
}

func (m MainModel) View() string {
	s := ""
	if m.ShowTopMenu {

		menuLeft := ""
		menuRight := ""

		menuLeft += fmt.Sprintf("Player %d: %s        Rounds: %d\n", *m.Game.ActivePlayerIdx+1, m.Game.ActivePlayer().Name, m.Game.ActivePlayer().Rounds)
		menuLeft += fmt.Sprintf("Score: %d            Herd: %d\n", m.Game.ActivePlayer().Score, m.Game.ActivePlayer().Herd)
		if len(m.Game.Discarded) > 0 {
			menuLeft += fmt.Sprintf("Discarded\n%s\n", m.formatDiscarded())
		}
		menuLeft += fmt.Sprintf("Market:\n%s\n", tui.RenderCards(m.Game.Market))
		menuLeft += fmt.Sprintf("Your hand:\n%s", tui.RenderCards(m.Game.ActivePlayer().Hand))

		columns := m.formatRemainingTokensColumn()
		rightStyle := lipgloss.NewStyle().Width(tui.Width - lipgloss.Width(menuLeft) - lipgloss.Width(columns)).Align(lipgloss.Right)

		menuRight += rightStyle.Render(fmt.Sprintf("Tokens\n%s\n", columns))

		s = tui.TopMenuStyle.Render(lipgloss.JoinHorizontal(0, menuLeft, menuRight))
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
	var g = game.NewGame()
	var v = models.NewNameEntry(&g)
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
