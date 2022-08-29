package internal

import (
	"strconv"
	"strings"
)

func FindColumnIndex(mainColumns string, requestedColumn string) int {
	columnsSplit := strings.Split(mainColumns, " ") // separator
	if columnNumber, err := strconv.Atoi(requestedColumn); err == nil {
		if columnNumber >= 0 && columnNumber < len(columnsSplit) {
			return columnNumber
		}
	}
	for idx, str := range columnsSplit {
		if str == requestedColumn {
			return idx
		}
	}
	return 0
}
