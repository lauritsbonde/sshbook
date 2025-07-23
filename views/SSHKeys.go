package views

import (
	"fmt"

	"github.com/gizak/termui/v3/widgets"
)

func SetupSSHKeys(width int, height int, p *widgets.Paragraph, keys []string, active bool, activeIndex int) {
	Outline(p, active, "SSH Keys")
	p.SetRect(0, height/2, width/2, height)

	extraLines := 8 // One for "Total SSH keys", one for spacing
	visibleRows := (height - (height / 2)) - extraLines
	if visibleRows < 1 {
		visibleRows = 1
	}

	start := 0
	if activeIndex >= visibleRows {
		start = activeIndex - visibleRows + 1
	}
	end := start + visibleRows
	if end > len(keys) {
		end = len(keys)
		start = end - visibleRows
		if start < 0 {
			start = 0
		}
	}

	p.Text = fmt.Sprintf("Total SSH keys: %d\n", len(keys))
	for i := start; i < end; i++ {
		if i == activeIndex {
			p.Text += "\n➤ " + keys[i]
		} else {
			p.Text += "\n- " + keys[i]
		}
	}

	p.Border = true
}
