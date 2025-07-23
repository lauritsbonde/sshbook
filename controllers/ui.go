package controllers

import (
	"log"
	"os"
	"sshbook/models"
	"sshbook/views"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func HandleKeyEvent(state *models.AppState, key string) {
	log.Printf("Key pressed: %s", key) // Debugging: Log the key pressed
	rerender := true                   // Flag to check if we need to re-render the UI
	switch key {
	case "q":
		ui.Close()
		os.Exit(0) // Exit the application
	case "h":
		state.ActivePane = "hosts"
	case "k":
		state.ActivePane = "keys"
	case "g":
		state.ActivePane = "groups"
	case "?", "+":
		state.ActivePane = "help"
	case "<Tab>":
		state.ActivePane = tabPaneShift(state, false) // Shift right on Tab
	case "<Escape>[Z":
		state.ActivePane = tabPaneShift(state, true) // Shift left on Shift-Tab
	case "<Up>":
		updateSelectedIndex(state, state.ActivePane, -1)
	case "<Down>":
		updateSelectedIndex(state, state.ActivePane, 1)
	case "<Enter>":
		handleEnterKey(state)

	default:
		log.Printf("Unhandled key: %s", key) // Debugging: Log unhandled keys
		rerender = false                     // No need to re-render for unhandled keys
	}

	if rerender {
		RenderUI(state)
	}
}

func RenderUI(state *models.AppState) {
	hosts := widgets.NewParagraph()
	sshKeys := widgets.NewParagraph()
	groups := widgets.NewParagraph()
	help := widgets.NewParagraph()

	termWidth, termHeight := ui.TerminalDimensions()

	views.SetupHosts(termWidth, termHeight, hosts, state.SSHDirContents.KnownHosts, isActivePane(state, "hosts"), state.SelectedIndex["hosts"])
	views.SetupSSHKeys(termWidth, termHeight, sshKeys, state.SSHDirContents.Keys, isActivePane(state, "keys"), state.SelectedIndex["keys"])
	views.SetupGroups(termWidth, termHeight, groups, isActivePane(state, "groups"))
	views.SetupHelp(termWidth, termHeight, help, isActivePane(state, "help"))

	grid := ui.NewGrid()
	grid.Set(
		ui.NewRow(0.4, hosts),
		ui.NewRow(0.4,
			ui.NewCol(0.5, sshKeys),
			ui.NewCol(0.5, groups),
		),
		ui.NewRow(0.2, help),
	)
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Clear() // Clear the UI to prevent overlapping renders
	ui.Render(grid)
}

func isActivePane(state *models.AppState, paneName models.Pane) bool {
	return state.ActivePane == paneName
}

func tabPaneShift(state *models.AppState, shiftPressed bool) models.Pane {
	index := -1
	for i, pane := range state.Panes {
		if pane == state.ActivePane {
			index = i
			break
		}
	}
	if index == -1 {
		log.Println("Active pane not found in state.Panes") // Debugging: Log if
		return state.ActivePane
	}
	if shiftPressed {
		index = (index - 1 + len(state.Panes)) % len(state.Panes) // Shift left
	} else {
		index = (index + 1) % len(state.Panes) // Shift right
	}
	return state.Panes[index] // Return the new active pane
}

func updateSelectedIndex(state *models.AppState, paneName models.Pane, direction int) {
	if state.SelectedIndex[paneName] == 0 && direction < 0 {
		return // Prevent going below 0
	}
	state.SelectedIndex[paneName] += direction
	log.Printf("Pane %s not found in state.Panes", paneName) // Debugging: Log if pane not found
}

func handleEnterKey(state *models.AppState) {

}

func sshConnection(state *models.AppState) {
	// Placeholder for SSH connection logic
	// This function would handle the SSH connection based on the selected host or key
	log.Println("SSH connection logic not implemented yet")
}
