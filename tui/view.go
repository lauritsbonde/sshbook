package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	// Row heights (borders eat 2 rows each; leave 1 for status line).
	statusH := 1
	usable := m.height - statusH
	topH := usable * 4 / 10
	midH := usable * 4 / 10
	botH := usable - topH - midH

	full := m.width
	half := full / 2

	hosts := m.listPane(paneHosts, "Known Hosts", m.hostRows(visibleRows(topH)), full, topH)

	keys := m.listPane(paneKeys, "SSH Keys", m.keyRows(visibleRows(midH)), half, midH)
	groups := m.textPane(paneGroups, "Groups", groupsText, full-half, midH)
	mid := lipgloss.JoinHorizontal(lipgloss.Top, keys, groups)

	help := m.textPane(paneHelp, "Help", helpText, full, botH)

	status := m.statusLine
	if status == "" {
		status = "h hosts · k keys · g groups · ? help · tab cycle · ↑/↓ move · enter select · q quit"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		hosts,
		mid,
		help,
		lipgloss.NewStyle().Faint(true).Render(truncate(status, full)),
	)
}

// style returns the border style for a pane based on active state.
func (m Model) style(pane string, w, h int) lipgloss.Style {
	s := inactiveBorder
	if m.active == pane {
		s = activeBorder
	}
	// Account for border (2 cols, 2 rows).
	return s.Width(w - 2).Height(h - 2)
}

func (m Model) title(pane, label string) string {
	if m.active == pane {
		return titleStyle.Foreground(lipgloss.Color("42")).Render("➤ " + label + " (active)")
	}
	return titleStyle.Render(label)
}

// listPane renders a bordered box with a title and pre-rendered rows.
func (m Model) listPane(pane, label string, body string, w, h int) string {
	content := m.title(pane, label) + "\n" + body
	return m.style(pane, w, h).Render(content)
}

func (m Model) textPane(pane, label, text string, w, h int) string {
	content := m.title(pane, label) + "\n\n" + text
	return m.style(pane, w, h).Render(content)
}

func (m Model) hostRows(visible int) string {
	hosts := m.ssh.KnownHosts
	var b strings.Builder
	fmt.Fprintf(&b, "Total known hosts: %d\n", len(hosts))
	b.WriteString(rows(mapFirstField(hosts), m.selected[paneHosts], visible))
	return b.String()
}

func (m Model) keyRows(visible int) string {
	keys := m.ssh.Keys
	var b strings.Builder
	fmt.Fprintf(&b, "Total SSH keys: %d\n", len(keys))
	b.WriteString(rows(keys, m.selected[paneKeys], visible))
	return b.String()
}

// visibleRows is how many list rows fit in a box of height h
// (2 border rows, 1 title row, 1 total-count row).
func visibleRows(h int) int {
	v := h - 4
	if v < 1 {
		return 1
	}
	return v
}

// rows renders a windowed list with a selection marker on the active index.
func rows(items []string, sel, visible int) string {
	start := 0
	if sel >= visible {
		start = sel - visible + 1
	}
	end := start + visible
	if end > len(items) {
		end = len(items)
		start = end - visible
		if start < 0 {
			start = 0
		}
	}

	var b strings.Builder
	for i := start; i < end; i++ {
		if i == sel {
			b.WriteString(selectedRow.Render("➤ "+items[i]) + "\n")
		} else {
			b.WriteString("- " + items[i] + "\n")
		}
	}
	return b.String()
}

func mapFirstField(items []string) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = firstField(it)
	}
	return out
}

func truncate(s string, w int) string {
	if lipgloss.Width(s) <= w {
		return s
	}
	return string([]rune(s)[:w])
}

const groupsText = "Groups organize SSH hosts and keys.\n\nCurrently, no groups are defined."

const helpText = "q quit · h hosts · k keys · g groups · ? help\n" +
	"tab / shift+tab cycle panes · ↑/↓ navigate · enter select"
