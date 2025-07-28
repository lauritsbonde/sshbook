package views

import (
	"sshbook/models"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func TabPane(appState *models.AppState) *widgets.TabPane {
	termWidth, termHeight := ui.TerminalDimensions()

	tabTitles := []string{"Dashboard"}
	for _, session := range appState.Tabs {
		tabTitles = append(tabTitles, session.Host)
	}

	tabpane := widgets.NewTabPane(tabTitles...)

	tabpane.SetRect(0, 0, termWidth, termHeight)
	tabpane.Border = true

	tabpane.PaddingBottom = 1

	if appState.ActiveIdx > 0 {
		tabpane.ActiveTabIndex = appState.ActiveIdx + 1 // offset for "Dashboard"
	} else {
		tabpane.ActiveTabIndex = 0 // default to dashboard
	}

	return tabpane
}
