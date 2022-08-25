package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func writeSequence(character rune, amount int, builder *strings.Builder) {
	for i := 0; i < amount; i++ {
		builder.WriteRune(character)
	}
}

func getAmount(builder *strings.Builder) int {
	if builder.Len() == 0 {
		return 1
	}
	amount, err := strconv.Atoi(builder.String())
	if err != nil {
		return 0
	}
	return amount
}

func unpack(s string) (string, error) {
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
				writeSequence(currentRune, getAmount(&currentLength), &output)
				currentLength.Reset()
			}
			currentRune = r
		}
		if idx == len(s)-1 {
			if currentRune != 0 {
				writeSequence(currentRune, getAmount(&currentLength), &output)
				currentLength.Reset()
			}

		}
	}
	if currentLength.Len() > 0 {
		return output.String(), errors.New("result might be wrong due to incorrect string given")
	}
	return output.String(), nil
}

func main() {
	fmt.Println("Enter packed string:")
	var packedString string
	_, err := fmt.Scanln(&packedString)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, "error while parsing packed string: "+err.Error())
		os.Exit(1)
	}
	unpackedString, err := unpack(packedString)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, "error unpacking string: "+err.Error())
		os.Exit(2)
	}
	fmt.Println(unpackedString)
	os.Exit(0)
}
