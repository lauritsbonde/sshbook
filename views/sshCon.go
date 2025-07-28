package views

import (
	"sshbook/models"
	"strings"

	"github.com/gizak/termui/v3/widgets"
)

func SetupSSHCon(width int, height int, p *widgets.Paragraph, state *models.AppState) {
	if state.ActivePane == "ssh" {
		active := state.SSHTabs.Active()
		if active != nil {
			p.Title = "SSH: " + active.Host
			p.Text = lastNLines(active.OutputLines, 20)
		} else {
			p.Title = "SSH"
			p.Text = "No active SSH session"
		}
	}

	p.SetRect(0, 0, width, height)
	p.Border = true
}

func lastNLines(lines []string, n int) string {
	if len(lines) <= n {
		return joinLines(lines)
	}
	return joinLines(lines[len(lines)-n:])
}

func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
