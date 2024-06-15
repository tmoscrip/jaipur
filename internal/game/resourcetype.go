package game

import "github.com/charmbracelet/lipgloss"

type ResourceType int

const (
	Diamond ResourceType = iota
	Gold
	Silver
	Cloth
	Spice
	Leather
	Camel
)

func (c ResourceType) Color() lipgloss.Color {
	switch c {
	case Diamond:
		return lipgloss.Color("#00FFFF")
	case Gold:
		return lipgloss.Color("#FFD700")
	case Silver:
		return lipgloss.Color("#C0C0C0")
	case Cloth:
		return lipgloss.Color("#FF00AA")
	case Spice:
		return lipgloss.Color("#FFA500")
	case Leather:
		return lipgloss.Color("#8B4513")
	case Camel:
		return lipgloss.Color("#FFD777")
	default:
		return lipgloss.Color("#FFFFFF")
	}
}

type StyleShortString struct {
	Style       lipgloss.Style
	ShortString string
}

var resourceTypeStyles = map[ResourceType]StyleShortString{
	Diamond: {
		Style:       lipgloss.NewStyle().Foreground(Diamond.Color()),
		ShortString: "Dia",
	},
	Gold: {
		Style:       lipgloss.NewStyle().Foreground(Gold.Color()),
		ShortString: "Gld",
	},
	Silver: {
		Style:       lipgloss.NewStyle().Foreground(Silver.Color()),
		ShortString: "Slv",
	},
	Cloth: {
		Style:       lipgloss.NewStyle().Foreground(Cloth.Color()),
		ShortString: "Cth",
	},
	Spice: {
		Style:       lipgloss.NewStyle().Foreground(Spice.Color()),
		ShortString: "Spi",
	},
	Leather: {
		Style:       lipgloss.NewStyle().Foreground(Leather.Color()),
		ShortString: "Lth",
	},
	Camel: {
		Style:       lipgloss.NewStyle().Foreground(Camel.Color()),
		ShortString: "Cml",
	},
}

func (c ResourceType) Style() lipgloss.Style {
	return resourceTypeStyles[c].Style
}

func (c ResourceType) ShortString() string {
	return resourceTypeStyles[c].ShortString
}

func (c ResourceType) String() string {
	return c.Style().Render(c.ShortString())
}
