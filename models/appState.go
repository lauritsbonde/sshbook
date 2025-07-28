package models

import (
	"os"
	"os/exec"
	"strings"
)

type Pane string

const (
	PaneHosts  Pane = "hosts"
	PaneKeys   Pane = "keys"
	PaneGroups Pane = "groups"
	PaneHelp   Pane = "help"
)

type SSHConfigEntry struct {
	Host         string
	HostName     string
	User         string
	IdentityFile string
}

type SSHDirContents struct {
	Keys             []string
	ConfigPath       string
	KnownHosts       []string
	SSHConfigEntries []SSHConfigEntry
}

type SSHSession struct {
	Host        string    // e.g. "github.com"
	OutputLines []string  // streamed output
	Cmd         *exec.Cmd // the SSH process
	Pty         *os.File  // PTY for interactive sessions
}

type TabType string

const (
	TabDashboard TabType = "dashboard"
	TabSSH       TabType = "ssh"
)

type Tab struct {
	Type TabType

	// Only for TabSSH
	Host        string
	OutputLines []string
	Cmd         *exec.Cmd
	Pty         *os.File
}

type AppState struct {
	Tabs      []Tab
	ActiveIdx int

	InputMode   bool   // whether we are in input mode (e.g. for selecting hosts)
	InputBuffer string // user input buffer

	SSHDirContents SSHDirContents // ssh dir contents
	SelectedIndex  map[Pane]int   // which index is selected in each dashboard pane
	DashboardPane  Pane           // which dashboard sub-pane is active
}

func (state *AppState) AddSSHSession(host string, cmd *exec.Cmd, pty *os.File) {
	state.Tabs = append(state.Tabs, Tab{
		Type:        TabSSH,
		Host:        host,
		Cmd:         cmd,
		Pty:         pty,
		OutputLines: []string{},
	})
	state.ActiveIdx = len(state.Tabs) - 1
}

func (state *AppState) ActiveTab() *Tab {
	if state.ActiveIdx >= 0 && state.ActiveIdx < len(state.Tabs) {
		return &state.Tabs[state.ActiveIdx]
	}
	return nil
}

func (s *AppState) CurrentSelectedHost() string {
	if s.ActiveTab().Type != TabDashboard || s.DashboardPane != PaneHosts {
		return ""
	}
	i := s.SelectedIndex[PaneHosts]
	if i >= 0 && i < len(s.SSHDirContents.KnownHosts) {
		// Only return the first token (the actual hostname)
		return strings.Fields(s.SSHDirContents.KnownHosts[i])[0]
	}
	return ""
}

func (s *AppState) NextTab() {
	if len(s.Tabs) > 0 {
		s.ActiveIdx = (s.ActiveIdx + 1) % len(s.Tabs)
	}
}

func (s *AppState) PrevTab() {
	if len(s.Tabs) > 0 {
		s.ActiveIdx = (s.ActiveIdx - 1 + len(s.Tabs)) % len(s.Tabs)
	}
}

func (s *AppState) CloseCurrentTab() {
	if s.ActiveIdx == 0 {
		// Never close dashboard
		return
	}
	s.Tabs = append(s.Tabs[:s.ActiveIdx], s.Tabs[s.ActiveIdx+1:]...)
	if s.ActiveIdx >= len(s.Tabs) {
		s.ActiveIdx = len(s.Tabs) - 1
	}
}
