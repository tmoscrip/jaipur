package tui

import "github.com/charmbracelet/lipgloss"

const Width = 60

const Black = lipgloss.Color("#000000")
const Night = lipgloss.Color("#111111")
const EerieBlack = lipgloss.Color("#232323")
const Jet = lipgloss.Color("#343434")
const OuterSpace = lipgloss.Color("#464646")
const DavysGray = lipgloss.Color("#575757")
const DimGray = lipgloss.Color("#696969")
const Gray = lipgloss.Color("#7a7a7a")

// some corresponding lighter shades that compliment
const Silver = lipgloss.Color("#9d9d9d")
const WhiteSmoke = lipgloss.Color("#c0c0c0")
const WhiteSmoke2 = lipgloss.Color("#d3d3d3")
const White = lipgloss.Color("#ffffff")

var TopMenuStyle = lipgloss.NewStyle().
	Width(Width).
	Align(lipgloss.Left).
	Padding(1, 3).
	Background(EerieBlack).
	Foreground(WhiteSmoke2).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(WhiteSmoke2)

var MenuOptionsStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Foreground(White).
	Background(EerieBlack).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(White).
	Padding(1, 2).
	Width(Width)

var ErrorStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Foreground(lipgloss.Color("#DD3333")).
	Background(EerieBlack).
	Bold(true).
	Border(lipgloss.RoundedBorder()).
	BorderTop(false).
	Padding(1).
	Width(Width)

var TitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#000000")).
	Bold(true).Border(lipgloss.RoundedBorder()).
	Padding(0, 4)
