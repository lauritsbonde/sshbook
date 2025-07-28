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
	PaneSSH    Pane = "ssh"
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

type Tabs struct {
	Sessions  []*SSHSession // one per host
	ActiveIdx int
}

type AppState struct {
	SSHDirContents SSHDirContents

	Panes      []Pane // List of all panes
	ActivePane Pane   // Currently active pane

	SelectedIndex map[Pane]int // Map to hold selected index for each pane
	SSHTabs       Tabs         // Tabs for SSH sessions

	InputMode   bool   // true = input mode (like Vim's insert mode)
	InputBuffer string // holds typed input while in input mode
}

func (s *AppState) CurrentSelectedHost() string {
	if s.ActivePane != PaneHosts {
		return ""
	}
	i := s.SelectedIndex[PaneHosts]
	if i >= 0 && i < len(s.SSHDirContents.KnownHosts) {
		// Only return the first token (the actual hostname)
		return strings.Fields(s.SSHDirContents.KnownHosts[i])[0]
	}
	return ""
}

func (t *Tabs) Add(host string, cmd *exec.Cmd, pty *os.File) {
	t.Sessions = append(t.Sessions, &SSHSession{
		Host: host,
		Cmd:  cmd,
		Pty:  pty,
	})
	t.ActiveIdx = len(t.Sessions) - 1
}

func (t *Tabs) Active() *SSHSession {
	if len(t.Sessions) == 0 || t.ActiveIdx < 0 || t.ActiveIdx >= len(t.Sessions) {
		return nil
	}
	return t.Sessions[t.ActiveIdx]
}

func (t *Tabs) Next() {
	if len(t.Sessions) > 0 {
		t.ActiveIdx = (t.ActiveIdx + 1) % len(t.Sessions)
	}
}

func (t *Tabs) Prev() {
	if len(t.Sessions) > 0 {
		t.ActiveIdx = (t.ActiveIdx - 1 + len(t.Sessions)) % len(t.Sessions)
	}
}
