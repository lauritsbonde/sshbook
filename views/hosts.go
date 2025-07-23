package views

import (
	"fmt"
	"strings"

	"github.com/gizak/termui/v3/widgets"
)

func SetupHosts(width int, height int, p *widgets.Paragraph, knownHosts []string, active bool, activeIndex int) {
	Outline(p, active, "Known Hosts")
	p.SetRect(0, 0, width, height/2)

	extraLines := 8                          // found by trial and error, adjust as needed
	visibleRows := (height / 2) - extraLines // title + padding
	if visibleRows < 1 {
		visibleRows = 1
	}

	start := 0
	if activeIndex >= visibleRows {
		start = activeIndex - visibleRows + 1
	}
	end := start + visibleRows
	if end > len(knownHosts) {
		end = len(knownHosts)
		start = end - visibleRows
		if start < 0 {
			start = 0
		}
	}

	p.Text = fmt.Sprintf("Total known hosts: %d\n", len(knownHosts))
	for i := start; i < end; i++ {
		shortHost := strings.Split(knownHosts[i], " ")[0]
		if i == activeIndex {
			shortHost = "➤ " + shortHost
		} else {
			shortHost = "- " + shortHost
		}
		p.Text += "\n" + shortHost
	}
	p.Border = true
}
