package views

import (
	"fmt"
	"sshbook/utils"
	"strings"

	"github.com/gizak/termui/v3/widgets"
)

func SetupHosts(width int, height int, p *widgets.Paragraph, knownHosts []string, active bool, activeIndex int) {
	Outline(p, active, "Known Hosts")
	p.SetRect(0, 0, width, height/2)

	visibleRows := utils.FilterRows(knownHosts, nil, height, activeIndex)

	p.Text = fmt.Sprintf("Total known hosts: %d\n", len(knownHosts))

	RenderFilteredList(p, visibleRows, activeIndex, func(host string) string {
		shortHost := strings.Split(host, " ")[0]
		return shortHost
	})

	// for i := 0; i < len(visibleRows.Rows); i++ {
	// 	shortHost := strings.Split(knownHosts[i], " ")[0]
	// 	if i+visibleRows.Start == activeIndex {
	// 		shortHost = "➤ " + shortHost
	// 	} else {
	// 		shortHost = strconv.Itoa(i+visibleRows.Start) + ": " + shortHost
	// 	}
	// 	p.Text += "\n" + shortHost
	// }
	p.Border = true
}
