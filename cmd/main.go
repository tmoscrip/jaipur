package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
	"github.com/tmoscrip/jaipur/models"
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

func formatCardRow(market []game.ResourceType, title string) string {
	// cards := make([]string, len(market))
	// for i, card := range market {
	// 	cards[i] = tui.RenderCard(card)
	// }
	// items := make([]string, 0)
	// items = append(items, tui.TitleStyle.MarginLeft(1).Render(title))
	// items = append(items, cards...)
	// return lipgloss.JoinHorizontal(lipgloss.Center, items...) + "\n"
	s := tui.TitleStyle.MarginLeft(1).Render(title) + "\n"
	s += tui.RenderCards(market) + "\n"
	return s
}

func (m MainModel) View() string {
	s := ""
	if m.ShowTopMenu {

		menuLeft := ""
		menuRight := ""

		menuLeft += fmt.Sprintf("Player %d: %s        Rounds: %d\n", m.Game.Players.ActiveIdx+1, m.Game.Players.Active().Name, m.Game.Players.Active().Rounds)
		menuLeft += fmt.Sprintf("Score: %d            Herd: %d\n", m.Game.Players.Active().Score, m.Game.Players.Active().Herd)
		if len(m.Game.Discarded) > 0 {
			menuLeft += fmt.Sprintf("Discarded\n%s\n", m.formatDiscarded())
		}
		menuLeft += formatCardRow(m.Game.Market, "Market")
		menuLeft += formatCardRow(m.Game.Players.Active().Hand, "Hand")

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
