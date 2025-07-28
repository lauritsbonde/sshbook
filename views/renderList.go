package views

import (
	"sshbook/utils"
	"strconv"

	"github.com/gizak/termui/v3/widgets"
)

func RenderFilteredList(p *widgets.Paragraph, filteredRows utils.FilteredRows, activeIndex int, rendererFunc func(string) string) {
	for i := 0; i < len(filteredRows.Rows); i++ {

		value := rendererFunc(filteredRows.Rows[i])

		if i == activeIndex {
			p.Text += "\n➤ " + value
		} else {
			p.Text += "\n " + strconv.Itoa(i+filteredRows.Start) + ": " + value
		}
	}
}
