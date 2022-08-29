package internal

import (
	"fmt"
	"strconv"
)

var suffixes = map[string]int{"K": 1000, "M": 1_000_000, "G": 1_000_000_000, "T": 1_000_000_000_000}

func GetSuffixValue(number string) int {
	n1, err := strconv.Atoi(number)
	if err == nil {
		fmt.Println(n1)
		return n1
	}
	num, suf := number[:len(number)-1], number[len(number)-1:]
	n2, _ := strconv.Atoi(num)
	fmt.Println(num, suf)
	return n2 * suffixes[suf]
}
