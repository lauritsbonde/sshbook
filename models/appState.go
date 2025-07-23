package models

import (
	"fmt"
)

type SSHDirContents struct {
	Keys       []string
	Config     string
	KnownHosts []string
}

func (s SSHDirContents) String() string {
	result := "SSH Directory Contents:\n"
	result += "Keys:\n"
	for _, key := range s.Keys {
		result += fmt.Sprintf("- %s\n", key)
	}
	result += fmt.Sprintf("Config file: %s\n", s.Config)
	result += fmt.Sprintf("Known hosts file: %s\n", s.KnownHosts)
	return result
}

type AppState struct {
	SSHDirContents SSHDirContents
	ActivePane     string         // "hosts", "keys", "groups", "help" - (thinking about making this an index)
	Panes          []string       // List of all panes
	SelectedIndex  map[string]int // Map to hold selected index for each pane
}
