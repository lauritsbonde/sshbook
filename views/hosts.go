package views

import (
	"fmt"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func SetupHosts(width int, height int, p *widgets.Paragraph, knownHosts []string) {
	p.Title = "Known Hosts"
	p.SetRect(0, 0, width, height/2)

	p.Text += "\n\n" + fmt.Sprintf("Total known hosts: %d", len(knownHosts))
	for _, host := range knownHosts {
		shortHost := strings.Split(host, " ")[0] // just "github.com,192.30.255.113"
		p.Text += "\n- " + shortHost
	}
	p.BorderStyle.Fg = ui.ColorCyan

	p.Border = true
}
