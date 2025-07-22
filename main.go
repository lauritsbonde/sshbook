package main

import (
	"fmt"
	"sshbook/controllers"
	"sshbook/views"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		fmt.Printf("Error initializing termui: %v\n", err)
		return
	}
	defer ui.Close()

	sshContents := controllers.SshDirContents()

	hosts := widgets.NewParagraph()
	sshKeys := widgets.NewParagraph()
	groups := widgets.NewParagraph()
	help := widgets.NewParagraph()

	termWidth, termHeight := ui.TerminalDimensions()

	views.SetupHosts(termWidth, termHeight, hosts, sshContents.KnownHosts)
	views.SetupSSHKeys(termWidth, termHeight, sshKeys, sshContents.Keys)
	views.SetupGroups(termWidth, termHeight, groups)
	views.SetupHelp(termWidth, termHeight, help)

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

	ui.Render(grid)

	for e := range ui.PollEvents() {
		switch e.Type {
		case ui.ResizeEvent:
			// Handle terminal resize
			termWidth, termHeight = ui.TerminalDimensions()
			grid.SetRect(0, 0, termWidth, termHeight)
			ui.Clear() // Clear the UI to prevent overlapping renders
			ui.Render(grid)
		case ui.KeyboardEvent:
			// Check if the 'q' key is pressed to quit
			if e.ID == "q" {
				return
			}
		}
	}
}
