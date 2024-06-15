package tui

import (
	lg "github.com/charmbracelet/lipgloss"
	"github.com/tmoscrip/jaipur/internal/game"
)

const Width = 80

const Black = lg.Color("#000000")
const Night = lg.Color("#111111")
const EerieBlack = lg.Color("#232323")
const Jet = lg.Color("#343434")
const OuterSpace = lg.Color("#464646")
const DavysGray = lg.Color("#575757")
const DimGray = lg.Color("#696969")
const Gray = lg.Color("#7a7a7a")
const Silver = lg.Color("#9d9d9d")
const WhiteSmoke = lg.Color("#c0c0c0")
const WhiteSmoke2 = lg.Color("#d3d3d3")
const White = lg.Color("#ffffff")

var TopMenuStyle = lg.NewStyle().
	Width(Width).
	Align(lg.Left).
	Padding(1, 3).
	Background(EerieBlack).
	Foreground(WhiteSmoke2).
	Border(lg.RoundedBorder()).
	BorderForeground(WhiteSmoke2)

var MenuOptionsStyle = lg.NewStyle().
	Align(lg.Center).
	Foreground(White).
	Background(EerieBlack).
	Border(lg.RoundedBorder()).
	BorderForeground(White).
	Padding(1, 2).
	Width(Width)

var ErrorStyle = lg.NewStyle().
	Align(lg.Center).
	Foreground(lg.Color("#DD3333")).
	Background(EerieBlack).
	Bold(true).
	Border(lg.RoundedBorder()).
	BorderTop(false).
	Padding(1).
	Width(Width)

var TitleStyle = lg.NewStyle().
	Foreground(lg.Color("#FFFFFF")).
	Background(lg.Color("#000000")).
	Bold(true).Border(lg.RoundedBorder()).
	Padding(0, 4)

var HelpStyle = TitleStyle.Bold(false)

type cardView struct {
	Resource game.ResourceType
	Selected bool
	Active   bool
}

func RenderCard(card cardView) string {
	style := lg.NewStyle().
		Border(lg.RoundedBorder()).
		Padding(1).
		Background(card.Resource.Color())

	if card.Selected {
		style = style.BorderStyle(lg.ThickBorder()).BorderForeground(lg.Color("#CC2222"))
	}
	if card.Active {
		style = style.BorderBackground(DimGray).Bold(true)
	}
	return style.Render(card.Resource.String())
}

func RenderCards(cards []game.ResourceType, selectedIndexes []int, cursor int) string {
	var cs = make([]string, 0)
	for rowIndex, card := range cards {
		selected, isCursor := false, false
		if rowIndex == cursor {
			isCursor = true
		}
		for _, selectedIdx := range selectedIndexes {
			if rowIndex == selectedIdx {
				selected = true
			}
		}
		cs = append(cs, RenderCard(cardView{card, selected, isCursor}))
	}

	return lg.JoinHorizontal(lg.Left, cs...)
}
