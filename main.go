package main

import (
	"fmt"
	"os"

	"sshbook/controllers"
	"sshbook/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.New(controllers.SshDirContents())

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
