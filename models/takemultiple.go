package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/internal/tui"
)

type takeMultipleOption struct {
	column   int
	Label    string
	Selected bool
	Index    int
}

var TableBorder = lipgloss.NewStyle().Foreground(tui.Silver).Background(lipgloss.Color("#000000"))

func (m takeMultipleOption) FormatRight(activeCol int, activeCursor int) string {
	var cursor = fmt.Sprintf(" ")
	if m.column == activeCol && m.Index == activeCursor {
		cursor = fmt.Sprintf(">")
	}
	var checked = fmt.Sprintf(" ")
	if m.Selected {
		checked = fmt.Sprintf("x")
	}
	return fmt.Sprintf("%s [%s] %s", cursor, checked, m.Label)
}

func (m takeMultipleOption) FormatLeft(activeCol int, activeCursor int) string {
	var cursor = fmt.Sprintf(" ")
	if m.column == activeCol && m.Index == activeCursor {
		cursor = fmt.Sprintf("<")
	}
	var checked = fmt.Sprintf(" ")
	if m.Selected {
		checked = fmt.Sprintf("x")
	}
	return fmt.Sprintf("%s [%s] %s", m.Label, checked, cursor)
}

func (m takeMultipleOption) CursorActive(cursorIdx int) bool {
	return cursorIdx == m.Index
}

func (m takeMultipleOption) columnActive(columnIdx int) bool {
	return columnIdx == m.column
}

type TakeMultiple struct {
	Game         *game.Game
	columns      map[int]map[int]takeMultipleOption
	Cursor       *int
	activecolumn *int
}

func (v TakeMultiple) Activecolumn() map[int]takeMultipleOption {
	return v.columns[*v.activecolumn]
}

func NewTakeMultiple(game *game.Game) TakeMultiple {
	market := make(map[int]takeMultipleOption)
	hand := make(map[int]takeMultipleOption)

	columns := make(map[int]map[int]takeMultipleOption)
	columns[0] = hand
	columns[1] = market
	for i, card := range game.ActivePlayer().Hand {
		hand[i] = takeMultipleOption{column: 0, Label: card.String(), Selected: false, Index: i}
	}
	for i, card := range game.Market {
		market[i] = takeMultipleOption{column: 1, Label: card.String(), Selected: false, Index: i}
	}
	return TakeMultiple{Game: game, Cursor: new(int), activecolumn: new(int), columns: columns}
}

func (v TakeMultiple) Init() tea.Cmd {
	return nil
}

func (v TakeMultiple) View() string {

	s := tui.TitleStyle.Render("Take Multiple")
	s += "\n"

	var max = 0
	if len(v.Game.ActivePlayer().Hand) > len(v.Game.Market) {
		max = len(v.Game.ActivePlayer().Hand)
	} else {
		max = len(v.Game.Market)
	}

	t := table.New().
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return lipgloss.NewStyle().Align(lipgloss.Center)
			}
			return lipgloss.NewStyle().Align(lipgloss.Center)
		}).Width(40).Headers("Hand", "Market").Border(lipgloss.NormalBorder()).BorderStyle(TableBorder)
	for i := 0; i < max; i++ {
		// rows = append(rows, []string{v.columns[0][i].FormatLeft(*v.activecolumn, *v.Cursor), v.columns[1][i].FormatRight(*v.activecolumn, *v.Cursor)})
		t.Row(v.columns[0][i].FormatLeft(*v.activecolumn, *v.Cursor), v.columns[1][i].FormatRight(*v.activecolumn, *v.Cursor))
	}

	s += t.Render()

	return s
}

func (v TakeMultiple) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return v, nil
}

func (v TakeMultiple) MyUpdate(msg tea.Msg) (tea.Model, tea.Cmd, string, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "b" {
			if *v.activecolumn == 1 {
				*v.activecolumn = 0
				return v, nil, "", nil
			}
			return v, nil, "selectActionMenu", nil
		}
		if msg.String() == "up" {
			if *v.Cursor > 0 {
				*v.Cursor = *v.Cursor - 1
			}
		}
		if msg.String() == "down" {
			if *v.Cursor < len(v.Game.Market)-1 {
				*v.Cursor = *v.Cursor + 1
			}
		}
		if msg.String() == "left" {
			if *v.activecolumn > 0 {
				*v.activecolumn = *v.activecolumn - 1
			}
		}
		if msg.String() == "right" {
			if *v.activecolumn < 1 {
				*v.activecolumn = *v.activecolumn + 1
			}
		}
		if msg.String() == "enter" {
			col := v.columns[*v.activecolumn]
			item := col[*v.Cursor]
			item.Selected = !item.Selected
			col[*v.Cursor] = item
		}
		if msg.String() == "n" {

			endRound, err := v.Game.PlayerTakeMultiple(v.selectedHand(), v.selectedMarket())
			if err != nil {
				return v, nil, "", err
			}
			if endRound {
				return v, nil, "endRound", nil
			}
			return v, nil, "selectActionMenu", nil
		}
	}
	return v, nil, "", nil
}

func (v TakeMultiple) selectedHand() []int {
	var selected []int
	for i, card := range v.columns[0] {
		if card.Selected {
			selected = append(selected, i)
		}
	}
	return selected
}

func (v TakeMultiple) selectedMarket() []int {
	var selected []int
	for i, card := range v.columns[1] {
		if card.Selected {
			selected = append(selected, i)
		}
	}
	return selected
}
