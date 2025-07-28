package utils

import (
	"sshbook/models"
)

type FilteredRows struct {
	Rows  []string
	Start int
	End   int
}

func FilterRows(rows []string, appState *models.AppState, height int, activeIndex int) FilteredRows {

	shownLines := rowsForHeight(rows, height, activeIndex)

	return shownLines
}

func rowsForHeight(rows []string, height int, activeIndex int) FilteredRows {
	extraLines := 9                          // found by trial and error, adjust as needed
	visibleRows := (height / 2) - extraLines // title + padding
	if visibleRows < 1 {
		visibleRows = 1
	}

	start := 0
	if activeIndex >= visibleRows {
		start = activeIndex - visibleRows + 1
	}
	end := start + visibleRows
	if end > len(rows) {
		end = len(rows)
		start = end - visibleRows
		if start < 0 {
			start = 0
		}
	}

	return FilteredRows{
		Rows:  rows[start : end+1],
		Start: start,
		End:   end,
	}
}
