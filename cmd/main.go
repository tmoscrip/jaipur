package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tmoscrip/jaipur/internal/game"
	"github.com/tmoscrip/jaipur/models"
)

func initialModel() tea.Model {
	var g = game.NewGame()
	var v = models.NewNameEntry(&g)
	return models.MainModel{
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
