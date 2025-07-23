package main

import (
	"fmt"
	"sshbook/controllers"
	"sshbook/models"

	ui "github.com/gizak/termui/v3"
)

func main() {
	if err := ui.Init(); err != nil {
		fmt.Printf("Error initializing termui: %v\n", err)
		return
	}
	defer ui.Close()

	appState := &models.AppState{
		SSHDirContents: controllers.SshDirContents(),
		ActivePane:     "hosts",
		Panes:          []string{"hosts", "keys", "groups", "help"},
		SelectedIndex: map[string]int{
			"hosts":  0,
			"keys":   0,
			"groups": 0,
			"help":   0,
		},
	}

	controllers.RenderUI(appState)

	for e := range ui.PollEvents() {
		switch e.Type {
		case ui.ResizeEvent:
			// Handle terminal resize
			controllers.RenderUI(appState)
		case ui.KeyboardEvent:
			controllers.HandleKeyEvent(appState, e.ID)
		}
	}
}
