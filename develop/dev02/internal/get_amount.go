package internal

import (
	"strconv"
	"strings"
)

func GetAmount(builder *strings.Builder) int {
	if builder.Len() == 0 {
		return 1
	}
	amount, err := strconv.Atoi(builder.String())
	if err != nil {
		return 0
	}
	return amount
}
