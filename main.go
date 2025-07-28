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
		DashboardPane:  models.PaneHosts,
		SelectedIndex: map[models.Pane]int{
			models.PaneHosts:  0,
			models.PaneKeys:   0,
			models.PaneGroups: 0,
			models.PaneHelp:   0,
		},
		Tabs: []models.Tab{
			{Type: models.TabDashboard},
		},
		ActiveIdx: 0,
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
