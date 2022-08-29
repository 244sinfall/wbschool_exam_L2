package main

import (
	"fmt"
	"io"
	"os"
	"unpacker/pkg/up"
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

func main() {
	fmt.Println("Enter packed string:")
	var packedString string
	_, err := fmt.Scanln(&packedString)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, "error while parsing packed string: "+err.Error())
		os.Exit(1)
	}
	unpackedString, err := up.Unpack(packedString)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, "error unpacking string: "+err.Error())
		os.Exit(2)
	}
	fmt.Println(unpackedString)
	os.Exit(0)
}
