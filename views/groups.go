package views

import (
	"github.com/gizak/termui/v3/widgets"
)

func SetupGroups(width int, height int, p *widgets.Paragraph, active bool) {
	Outline(p, active, "Groups")

	p.Text = "Groups\n\n" +
		"Groups are used to organize SSH hosts and keys.\n\n" +
		"Currently, no groups are defined.\n\n"

	// Set the rectangle dimensions
	p.SetRect(width/2, height/2, width, height)
}
