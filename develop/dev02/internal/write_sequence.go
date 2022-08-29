package internal

import "strings"

func WriteSequence(character rune, amount int, builder *strings.Builder) {
	for i := 0; i < amount; i++ {
		builder.WriteRune(character)
	}
}
