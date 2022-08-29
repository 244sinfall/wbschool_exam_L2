package up

import (
	"errors"
	"strings"
	"unicode"
	"unpacker/internal"
)

func Unpack(s string) (string, error) {
	// Превращаем строку в руны для итерации
	if len(s) == 0 {
		return "", errors.New("empty string given")
	}
	// Создаем билдер для вывода
	output := strings.Builder{}
	// Сначала будем считывать руну не цифру
	var currentRune rune
	// Билдер строки для числа (в числе может быть более 1 руны)
	currentLength := strings.Builder{}
	// Для escape последовательностей. Даем одному символу избежать проверки на символ
	var escape bool
	// Итерируемся по рунам на входе
	for idx, r := range s {
		if r == '\\' && !escape {
			escape = true
			continue
		}
		if unicode.IsDigit(r) && !escape { // Если цифра - сразу записываем в Length и выходим
			currentLength.WriteRune(r)
		} else {
			if escape {
				escape = false
			}
			if currentRune != 0 { // Уже какая-то руна была записана
				internal.WriteSequence(currentRune, internal.GetAmount(&currentLength), &output)
				currentLength.Reset()
			}
			currentRune = r
		}
		if idx == len(s)-1 {
			if currentRune != 0 {
				internal.WriteSequence(currentRune, internal.GetAmount(&currentLength), &output)
				currentLength.Reset()
			}

		}
	}
	if currentLength.Len() > 0 {
		return output.String(), errors.New("result might be wrong due to incorrect string given")
	}
	return output.String(), nil
}
