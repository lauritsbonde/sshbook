package views

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func SetupSSHKeys(width int, height int, p *widgets.Paragraph, keys []string) {
	p.Title = "SSH Keys"
	p.SetRect(0, height/2, width, height)
	p.Text += "\n\n" + fmt.Sprintf("Total SSH keys: %d", len(keys))
	for _, key := range keys {
		p.Text += "\n- " + key
	}
	p.BorderStyle.Fg = ui.ColorGreen
	p.Border = true
}
