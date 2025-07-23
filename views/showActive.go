package views

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func Outline(p *widgets.Paragraph, active bool, pane string) {
	if active {
		p.Title = "➤ " + pane + " (Active)" // Add an indicator to the title
		p.BorderStyle.Fg = ui.ColorGreen
	} else {
		p.Title = pane
		p.BorderStyle.Fg = ui.ColorWhite
	}

}
