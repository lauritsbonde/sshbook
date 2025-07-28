package views

import (
	"fmt"
	"sshbook/utils"
	"strconv"

	"github.com/gizak/termui/v3/widgets"
)

func RenderFilteredList(p *widgets.Paragraph, filteredRows utils.FilteredRows, activeIndex int, rendererFunc func(string) string) {
	maxIndex := filteredRows.Start + len(filteredRows.Rows) - 1
	indexDigits := len(strconv.Itoa(maxIndex)) // number of digits, e.g. 3 for "201"
	indexPrefixWidth := indexDigits + 2        // space for "###: "

	for i := 0; i < len(filteredRows.Rows); i++ {
		value := rendererFunc(filteredRows.Rows[i])
		idx := i + filteredRows.Start

		if i == activeIndex-filteredRows.Start {
			// Arrow with same width as the index line
			arrowPrefix := fmt.Sprintf("%*s", indexPrefixWidth, "➤")
			p.Text += "\n" + arrowPrefix + " " + value
		} else {
			indexPrefix := fmt.Sprintf("%*d: ", indexDigits, idx)
			p.Text += "\n " + indexPrefix + value
		}
	}
}
