package tui

import (
	"fmt"
	"os/exec"

	"sshbook/controllers"
	"sshbook/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Pane identifiers.
const (
	paneConnections = "connections"
	paneKeys        = "keys"
	paneHelp        = "help"
)

var panes = []string{paneConnections, paneKeys, paneHelp}

// mode is the top-level UI mode.
type mode int

const (
	modeList mode = iota
	modeForm
)

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
	ssh      models.SSHDirContents
	active   string
	selected map[string]int
	width    int
	height   int
	status   string
	mode     mode
	form     form
}

// New builds the initial model from the SSH directory contents.
func New(ssh models.SSHDirContents) Model {
	return Model{
		ssh:    ssh,
		active: paneConnections,
		selected: map[string]int{
			paneConnections: 0,
			paneKeys:        0,
			paneHelp:        0,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// connectFinishedMsg is delivered after an interactive ssh session exits.
type connectFinishedMsg struct {
	name string
	err  error
}

// itemCount returns how many selectable items a pane has.
func (m Model) itemCount(pane string) int {
	switch pane {
	case paneConnections:
		return len(m.ssh.Connections)
	case paneKeys:
		return len(m.ssh.Keys)
	default:
		return 0
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case connectFinishedMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("ssh %s failed: %v", msg.name, msg.err)
		} else {
			m.status = fmt.Sprintf("ssh %s session ended", msg.name)
		}
		return m, nil

	case tea.KeyMsg:
		if m.mode == modeForm {
			return m.updateForm(msg)
		}
		return m.updateList(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "c":
		m.active = paneConnections
	case "k":
		m.active = paneKeys
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
	case "a":
		m.form = newForm(m.ssh.Keys)
		m.mode = modeForm
		m.status = ""
	case "enter":
		return m, m.connect()
	}
	return m, nil
}

func (m Model) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f, result, cmd := m.form.update(msg)
	m.form = f
	switch result {
	case formCancelled:
		m.mode = modeList
		return m, nil
	case formSubmitted:
		conn := m.form.connection()
		if err := controllers.AppendConnection(m.ssh.Config, conn); err != nil {
			m.form.err = err.Error()
			return m, nil
		}
		m.ssh = controllers.SshDirContents() // reload with the new entry
		m.mode = modeList
		m.active = paneConnections
		m.status = fmt.Sprintf("added connection %q", conn.Name)
		return m, nil
	}
	return m, cmd
}

// connect suspends the TUI and runs an interactive ssh session for the
// selected connection.
func (m Model) connect() tea.Cmd {
	if m.active != paneConnections || len(m.ssh.Connections) == 0 {
		return nil
	}
	conn := m.ssh.Connections[m.selected[paneConnections]]
	c := exec.Command("ssh", conn.Name)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return connectFinishedMsg{name: conn.Name, err: err}
	})
}

// move changes the selected index of the active pane, clamped to its items.
func (m *Model) move(dir int) {
	n := m.itemCount(m.active)
	if n == 0 {
		return
	}
	idx := m.selected[m.active] + dir
	if idx < 0 {
		idx = 0
	}
	if idx > n-1 {
		idx = n - 1
	}
	m.selected[m.active] = idx
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
