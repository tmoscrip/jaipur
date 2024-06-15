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
	case "startTurn":
		m.ShowTopMenu = false
		m.ActiveView = NewStartTurn(m.Game)
	}
	return m, cmd
}

func formatRemainingTokensColumn(tokens game.ResourceTokens) string {
	s := "Tokens\n"
	for _, resource := range []game.ResourceType{game.Diamond, game.Gold, game.Silver, game.Cloth, game.Spice, game.Leather} {
		count := 0
		for _, token := range tokens[resource] {
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

func formatDiscardedColumn(discarded []game.ResourceType) string {
	s := "Discarded\n"
	for _, resource := range []game.ResourceType{game.Diamond, game.Gold, game.Silver, game.Cloth, game.Spice, game.Leather, game.Camel} {
		count := 0
		for _, card := range discarded {
			if card == resource {
				count++
			}
		}
		if count != 0 {
			s += fmt.Sprintf("(%d) %s\n", count, resource.String())
		}
	}
	return s
}

func formatCardRow(cards []game.ResourceType, selectedIndexes []int, cursor int, title string) string {
	s := tui.TitleStyle.MarginLeft(2).Render(title)
	s += "\n"
	s += tui.RenderCards(cards, selectedIndexes, cursor) + "\n"
	return s
}

func formatPlayerInfo(p *game.Player) string {
	s := p.Name + " "
	for i := 0; i < p.Rounds; i++ {
		s += "•"
	}
	return tui.TitleStyle.Render(s) + "\n"
}

func formatPlayerScores(p *game.Player) string {
	s := fmt.Sprintf("( %d ) ( 🐫 %d )", p.Score, p.Herd)
	return tui.TitleStyle.Render(s)
}

func renderTopMenu(g *game.Game) string {
	s, menuLeft, menuRight := "", "", ""

	playerInfo := formatPlayerInfo(g.Players.Active())
	scoreInfo := formatPlayerScores(g.Players.Active())

	menuLeft += lipgloss.JoinHorizontal(lipgloss.Top, playerInfo, scoreInfo)
	menuLeft += "\n"

	menuLeft += formatCardRow(g.Market, g.MarketSelected, g.MarketCursor, "Market")
	menuLeft += formatCardRow(g.Players.Active().Hand, g.HandSelected, g.HandCursor, "Hand")

	menuRight += formatRemainingTokensColumn(g.ResourceTokens)
	menuRight += formatDiscardedColumn(g.Discarded)
	rightStyle := lipgloss.NewStyle().Width(tui.Width - lipgloss.Width(menuLeft) - lipgloss.Width(menuRight)).Align(lipgloss.Right)
	menuLeft = lipgloss.NewStyle().Background(tui.EerieBlack).Render(menuLeft)
	menuRight = rightStyle.Height(lipgloss.Height(menuLeft)).Background(tui.EerieBlack).Render(menuRight)
	joined := lipgloss.JoinHorizontal(lipgloss.Top, menuLeft, menuRight)
	s = tui.TopMenuStyle.Render(joined)
	s += "\n"
	return s
}

func (m MainModel) View() string {
	s := ""
	if m.ShowTopMenu {
		s += renderTopMenu(m.Game)
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
