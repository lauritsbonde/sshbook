package views

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func SetupGroups(width int, height int, p *widgets.Paragraph) {
	p.Title = "Groups"
	p.SetRect(width/2, height/2, width, height)
	p.Text += "\n\n" + "No groups available."
	p.BorderStyle.Fg = ui.ColorYellow
	p.Border = true
}
