package tui

import (
	"fmt"
	"strings"

	"sshbook/models"

	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

// Pane identifiers.
const (
	paneHosts  = "hosts"
	paneKeys   = "keys"
	paneGroups = "groups"
	paneHelp   = "help"
)

var panes = []string{paneHosts, paneKeys, paneGroups, paneHelp}

// Styles.
var (
	activeBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("42")) // green

	inactiveBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")) // grey

	titleStyle = lipgloss.NewStyle().Bold(true)

	selectedRow = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
)

// Model is the root bubbletea model.
type Model struct {
	ssh        models.SSHDirContents
	active     string
	selected   map[string]int
	width      int
	height     int
	statusLine string
}

// New builds the initial model from the SSH directory contents.
func New(ssh models.SSHDirContents) Model {
	return Model{
		ssh:    ssh,
		active: paneHosts,
		selected: map[string]int{
			paneHosts:  0,
			paneKeys:   0,
			paneGroups: 0,
			paneHelp:   0,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// itemsFor returns the selectable items for a pane (nil for non-list panes).
func (m Model) itemsFor(pane string) []string {
	switch pane {
	case paneHosts:
		return m.ssh.KnownHosts
	case paneKeys:
		return m.ssh.Keys
	default:
		return nil
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h":
			m.active = paneHosts
		case "k":
			m.active = paneKeys
		case "g":
			m.active = paneGroups
		case "?", "+":
			m.active = paneHelp
		case "tab":
			m.active = shift(m.active, 1)
		case "shift+tab":
			m.active = shift(m.active, -1)
		case "up":
			m.move(-1)
		case "down":
			m.move(1)
		case "enter":
			m.handleEnter()
		}
	}
	return m, nil
}

// move changes the selected index of the active pane, clamped to its items.
func (m *Model) move(dir int) {
	items := m.itemsFor(m.active)
	if len(items) == 0 {
		return
	}
	idx := m.selected[m.active] + dir
	if idx < 0 {
		idx = 0
	}
	if idx > len(items)-1 {
		idx = len(items) - 1
	}
	m.selected[m.active] = idx
}

// handleEnter is the placeholder for SSH connection logic.
func (m *Model) handleEnter() {
	items := m.itemsFor(m.active)
	if len(items) == 0 {
		return
	}
	sel := items[m.selected[m.active]]
	m.statusLine = fmt.Sprintf("selected %s: %s (connect not implemented)", m.active, firstField(sel))
}

// shift cycles to the next/previous pane.
func shift(active string, dir int) string {
	idx := 0
	for i, p := range panes {
		if p == active {
			idx = i
			break
		}
	}
	idx = (idx + dir + len(panes)) % len(panes)
	return panes[idx]
}

func firstField(s string) string {
	return strings.Split(s, " ")[0]
}
