package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type mainModel interface {
	tea.Model
	MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error)
}

type MainModel struct {
	ActiveView   mainModel
	Game         *game.Game
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
			return m, tea.Quit
		}
	}

	_, cmd, transition, err := m.ActiveView.MyUpdate(msg)

	if err != nil {
		m.ErrorMessage = err.Error()
		return m, cmd
	}

	switch transition {
	case "selectActionMenu":
		m.ShowTopMenu = true
		m.ActiveView = NewSelectAction(m.Game)
	case "takeOneCard":
		m.ShowTopMenu = true
		m.ActiveView = NewTakeOne(m.Game)
	case "takeSeveralCards":
		m.ShowTopMenu = true
		m.ActiveView = NewTakeMultiple(m.Game)
	case "takeCamels":
		m.ShowTopMenu = true
		m.ActiveView = NewTakeCamels(m.Game)
	case "sellCards":
		m.ShowTopMenu = true
		m.ActiveView = NewSellCards(m.Game)
	case "endRound":
		m.ShowTopMenu = false
		m.ActiveView = NewEndRound(m.Game)
	}
	return m, cmd
}

var numberStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

func (m MainModel) formatRemainingTokensColumn() string {
	s := "Tokens\n"
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
	return s + "\n"
}

func (m MainModel) formatDiscardedColumn() string {
	s := "Discarded\n"
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

func formatCardRow(cards []game.ResourceType, selectedIndexes []int, cursor int, title string) string {
	s := tui.TitleStyle.MarginLeft(1).Render(title) + "\n"

	s += tui.RenderCards(cards, selectedIndexes, cursor) + "\n"
	return s
}

func (m MainModel) View() string {
	s := ""
	if m.ShowTopMenu {

		menuLeft := ""
		menuRight := ""

		menuLeft += fmt.Sprintf("Player %d: %s        Rounds: %d\n", m.Game.Players.ActiveIdx+1, m.Game.Players.Active().Name, m.Game.Players.Active().Rounds)
		menuLeft += fmt.Sprintf("Score: %d            Herd: %d\n", m.Game.Players.Active().Score, m.Game.Players.Active().Herd)
		menuLeft += formatCardRow(m.Game.Market, m.Game.MarketSelected, m.Game.MarketCursor, "Market")
		menuLeft += formatCardRow(m.Game.Players.Active().Hand, m.Game.HandSelected, m.Game.HandCursor, "Hand")

		menuRight += m.formatRemainingTokensColumn()
		menuRight += m.formatDiscardedColumn()
		rightStyle := lipgloss.NewStyle().Width(tui.Width - lipgloss.Width(menuLeft) - lipgloss.Width(menuRight)).Align(lipgloss.Right)
		menuRight = rightStyle.Render(menuRight)
		joined := lipgloss.JoinHorizontal(lipgloss.Top, menuLeft, menuRight)
		s = tui.TopMenuStyle.Render(joined)
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
