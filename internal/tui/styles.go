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

func RenderCard(card game.ResourceType, active bool, selected bool) string {
	style := lg.NewStyle().
		Border(lg.RoundedBorder()).
		Padding(1).
		Background(card.Color())

	if active {
		style = style.BorderBackground(lg.Color("#FF0000"))
	}
	return style.Render(card.String())
}

func RenderCards(cards []game.ResourceType) string {
	var cs = make([]string, 0)
	for _, card := range cards {
		cs = append(cs, RenderCard(card, false, false))
	}

	return lg.JoinHorizontal(lg.Left, cs...)
}
