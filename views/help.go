package views

import (
	"github.com/gizak/termui/v3/widgets"
)

// SetupHelp initializes the help view with instructions for the user.
func SetupHelp(width int, height int, p *widgets.Paragraph, active bool) {
	Outline(p, active, "Help")

	p.Text = "Help\nPress 'q' to quit.\n\n" +
		"Navigation:\n" +
		"- Use arrow keys to navigate through the interface.\n" +
		"- Press 'h' for help.\n" +
		"- Press 'q' to quit the application.\n\n" +
		"Known Hosts:\n" +
		"- Displays a list of known SSH hosts.\n\n" +
		"SSH Keys:\n" +
		"- Displays a list of SSH keys.\n\n" +
		"Groups:\n" +
		"- Displays available groups (if any)."
	p.SetRect(0, 0, width, height)

	p.Border = true

}
