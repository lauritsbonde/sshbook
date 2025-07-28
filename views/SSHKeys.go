package views

import (
	"fmt"
	"sshbook/utils"

	"github.com/gizak/termui/v3/widgets"
)

func SetupSSHKeys(width int, height int, p *widgets.Paragraph, keys []string, active bool, activeIndex int) {
	Outline(p, active, "SSH Keys")
	p.SetRect(0, height/2, width/2, height)

	filteredRows := utils.FilterRows(keys, nil, height, activeIndex)

	p.Text = fmt.Sprintf("Total SSH keys: %d\n", len(keys))

	RenderFilteredList(p, filteredRows, activeIndex, func(key string) string {
		return key
	})

	p.Border = true
}
