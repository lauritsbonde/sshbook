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
			if err == nil && state.DashboardPane == models.PaneHosts && index >= 0 && index < len(state.SSHDirContents.KnownHosts) {
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
				rerender = false
			}
		}
	} else {
		switch key {
		case "q":
			activeTab := state.ActiveTab()
			if activeTab.Type == models.TabSSH {
				if activeTab != nil && activeTab.Type == models.TabSSH {
					activeTab.Pty.Close()   // Close the PTY
					state.CloseCurrentTab() // Remove the SSH tab
				}
			} else {
				ui.Close()
				os.Exit(0)
			}

		case "h":
			state.DashboardPane = models.PaneHosts
		case "k":
			state.DashboardPane = models.PaneKeys
		case "g":
			state.DashboardPane = models.PaneGroups
		case "?", "+":
			state.DashboardPane = models.PaneHelp
		case "<Tab>":
			state.DashboardPane = tabPaneShift(state, false)
		case "<Escape>[Z":
			state.DashboardPane = tabPaneShift(state, true)
		case "<Up>":
			updateSelectedIndex(state, state.DashboardPane, -1)
		case "<Down>":
			updateSelectedIndex(state, state.DashboardPane, 1)
		case "<Enter>":
			handleEnterKey(state)
		case "i":
			state.InputMode = true
			state.InputBuffer = ""
		case "<Right>":
			state.NextTab()
		case "<Left>":
			state.PrevTab()
		case "x":
			state.CloseCurrentTab()
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
	var grid *ui.Grid = ui.NewGrid()

	tabPane := views.TabPane(state)

	var content *ui.Grid

	if tab := state.ActiveTab(); tab != nil && tab.Type == models.TabSSH {
		content = renderSSHCon(state)
	} else {
		content = renderStartScreen(state)
	}

	grid.Set(
		ui.NewRow(0.05, tabPane),
		ui.NewRow(0.95, content),
	)

	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	ui.Clear()
	ui.Render(grid)
}

func isActivePane(state *models.AppState, paneName models.Pane) bool {
	return state.DashboardPane == paneName
}

func renderStartScreen(state *models.AppState) *ui.Grid {
	hosts := widgets.NewParagraph()
	sshKeys := widgets.NewParagraph()
	groups := widgets.NewParagraph()
	help := widgets.NewParagraph()
	status := widgets.NewParagraph()

	termWidth, termHeight := ui.TerminalDimensions()

	views.SetupHosts(termWidth, termHeight, hosts, state.SSHDirContents.KnownHosts, isActivePane(state, models.PaneHosts), state.SelectedIndex[models.PaneHosts])
	views.SetupSSHKeys(termWidth, termHeight, sshKeys, state.SSHDirContents.Keys, isActivePane(state, models.PaneKeys), state.SelectedIndex[models.PaneKeys])
	views.SetupGroups(termWidth, termHeight, groups, isActivePane(state, models.PaneGroups))
	views.SetupHelp(termWidth, termHeight, help, isActivePane(state, models.PaneHelp))

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
		ui.NewRow(0.3,
			ui.NewCol(0.5, sshKeys),
			ui.NewCol(0.5, groups),
		),
		ui.NewRow(0.15, help),
		ui.NewRow(0.1, status),
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

func tabPaneShift(state *models.AppState, shiftPressed bool) models.Pane {
	panes := []models.Pane{models.PaneHosts, models.PaneKeys, models.PaneGroups, models.PaneHelp}
	index := -1

	for i, pane := range panes {
		if pane == state.DashboardPane {
			index = i
			break
		}
	}

	if index == -1 {
		log.Println("Active pane not found in dashboard panes")
		return state.DashboardPane
	}

	if shiftPressed {
		index = (index - 1 + len(panes)) % len(panes)
	} else {
		index = (index + 1) % len(panes)
	}

	return panes[index]
}

func updateSelectedIndex(state *models.AppState, paneName models.Pane, direction int) {
	if state.SelectedIndex[paneName] == 0 && direction < 0 {
		return // Prevent going below 0
	}
	state.SelectedIndex[paneName] += direction
	log.Printf("Pane %s not found in state.Panes", paneName) // Debugging: Log if pane not found
}

func handleEnterKey(state *models.AppState) {
	if state.DashboardPane == models.PaneHosts {
		host := state.CurrentSelectedHost()
		if host != "" {
			go startSSHSession(state, host)
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

	state.AddSSHSession(host, cmd, ptmx)

	go func() {
		defer ptmx.Close()
		buf := make([]byte, 1024)

		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				output := string(buf[:n])
				lines := splitLines(output)

				tab := state.ActiveTab()
				if tab != nil && tab.Type == models.TabSSH && tab.Pty == ptmx {
					tab.OutputLines = append(tab.OutputLines, lines...)
				}

				ui.Clear()
				RenderUI(state)
			}
			if err != nil {
				tab := state.ActiveTab()
				if tab != nil && tab.Type == models.TabSSH && tab.Pty == ptmx {
					tab.OutputLines = append(tab.OutputLines, "[SSH session ended]")
				}
				break
			}
		}
	}()

	RenderUI(state)
}

func splitLines(s string) []string {
	return strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
}
