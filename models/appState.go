package models

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

type AppState struct {
	SSHDirContents SSHDirContents
	Panes          []Pane       // List of all panes
	ActivePane     Pane         // Currently active pane
	SelectedIndex  map[Pane]int // Map to hold selected index for each pane
}

func (s *AppState) CurrentSelectedHost() string {
	if s.ActivePane != PaneHosts {
		return ""
	}
	i := s.SelectedIndex[PaneHosts]
	if i >= 0 && i < len(s.SSHDirContents.KnownHosts) {
		return s.SSHDirContents.KnownHosts[i]
	}
	return ""
}
