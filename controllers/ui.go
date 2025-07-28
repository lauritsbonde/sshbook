package controllers

import (
	"log"
	"os"
	"os/exec"
	"sshbook/models"
	"sshbook/views"
	"strconv"
	"strings"

	"github.com/creack/pty"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func HandleKeyEvent(state *models.AppState, key string) {
	log.Printf("Key pressed: %s", key)
	rerender := true

	if state.InputMode {
		switch key {
		case "<Escape>":
			state.InputMode = false
			state.InputBuffer = ""
		case "<Enter>":
			index, err := strconv.Atoi(state.InputBuffer)
			if err == nil && state.ActivePane == models.PaneHosts && index >= 0 && index < len(state.SSHDirContents.KnownHosts) {
				state.SelectedIndex[models.PaneHosts] = index
			}
			state.InputMode = false
			state.InputBuffer = ""
		default:
			if len(key) == 1 {
				state.InputBuffer += key
			} else if key == "<Backspace>" {
				if len(state.InputBuffer) > 0 {
					state.InputBuffer = state.InputBuffer[:len(state.InputBuffer)-1]
				}
			} else {
				log.Printf("Unhandled input in input mode: %s", key)
				rerender = false // Don't rerender for unhandled keys in input mode
			}
		}
	} else {
		switch key {
		case "q":
			ui.Close()
			os.Exit(0)
		case "h":
			state.ActivePane = "hosts"
		case "k":
			state.ActivePane = "keys"
		case "g":
			state.ActivePane = "groups"
		case "?", "+":
			state.ActivePane = "help"
		case "<Tab>":
			state.ActivePane = tabPaneShift(state, false)
		case "<Escape>[Z":
			state.ActivePane = tabPaneShift(state, true)
		case "<Up>":
			updateSelectedIndex(state, state.ActivePane, -1)
		case "<Down>":
			updateSelectedIndex(state, state.ActivePane, 1)
		case "<Enter>":
			handleEnterKey(state)
		case "i":
			state.InputMode = true
			state.InputBuffer = ""
		default:
			log.Printf("Unhandled key: %s", key)
			rerender = false
		}
	}

	if rerender {
		RenderUI(state)
	}
}

func RenderUI(state *models.AppState) {

	var grid *ui.Grid

	switch state.ActivePane {
	case "hosts", "keys", "groups", "help":
		grid = renderStartScreen(state)
	case "ssh":
		grid = renderSSHCon(state)
	}

	ui.Clear() // Clear the UI to prevent overlapping renders
	ui.Render(grid)
}

func renderStartScreen(state *models.AppState) *ui.Grid {
	hosts := widgets.NewParagraph()
	sshKeys := widgets.NewParagraph()
	groups := widgets.NewParagraph()
	help := widgets.NewParagraph()
	status := widgets.NewParagraph()

	termWidth, termHeight := ui.TerminalDimensions()

	views.SetupHosts(termWidth, termHeight, hosts, state.SSHDirContents.KnownHosts, isActivePane(state, "hosts"), state.SelectedIndex["hosts"])
	views.SetupSSHKeys(termWidth, termHeight, sshKeys, state.SSHDirContents.Keys, isActivePane(state, "keys"), state.SelectedIndex["keys"])
	views.SetupGroups(termWidth, termHeight, groups, isActivePane(state, "groups"))
	views.SetupHelp(termWidth, termHeight, help, isActivePane(state, "help"))

	// ── Status bar content
	if state.InputMode {
		status.Text = "[INPUT MODE] Buffer: " + state.InputBuffer
		status.Border = false
		status.TextStyle.Fg = ui.ColorYellow
	} else {
		status.Text = "[NORMAL MODE] Press 'i' to enter input mode"
		status.Border = false
		status.TextStyle.Fg = ui.ColorCyan
	}

	grid := ui.NewGrid()
	grid.Set(
		ui.NewRow(0.4, hosts),
		ui.NewRow(0.4,
			ui.NewCol(0.5, sshKeys),
			ui.NewCol(0.5, groups),
		),
		ui.NewRow(0.15, help),
		ui.NewRow(0.05, status),
	)
	grid.SetRect(0, 0, termWidth, termHeight)

	return grid
}

func renderSSHCon(state *models.AppState) *ui.Grid {
	sshTab := widgets.NewParagraph()

	termWidth, termHeight := ui.TerminalDimensions()

	views.SetupSSHCon(termWidth, termHeight, sshTab, state)

	grid := ui.NewGrid()
	grid.Set(
		ui.NewRow(1.0, sshTab),
	)
	grid.SetRect(0, 0, termWidth, termHeight)
	return grid
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
	if state.ActivePane == models.PaneHosts {
		host := state.CurrentSelectedHost()
		if host != "" {
			go startSSHSession(state, host)
			state.ActivePane = models.PaneSSH
		}
	}
}

func startSSHSession(state *models.AppState, host string) {
	cmd := exec.Command("ssh", host)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("Failed to start SSH: %v", err)
		return
	}

	session := &models.SSHSession{
		Host:        host,
		Cmd:         cmd,
		Pty:         ptmx,
		OutputLines: []string{},
	}
	state.SSHTabs.Sessions = append(state.SSHTabs.Sessions, session)
	state.SSHTabs.ActiveIdx = len(state.SSHTabs.Sessions) - 1

	go func() {
		defer ptmx.Close()
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				output := string(buf[:n])
				lines := splitLines(output)
				session.OutputLines = append(session.OutputLines, lines...)

				ui.Clear()
				RenderUI(state)
			}
			if err != nil {
				session.OutputLines = append(session.OutputLines, "[SSH session ended]")
				break
			}
		}
	}()
}

func splitLines(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}
