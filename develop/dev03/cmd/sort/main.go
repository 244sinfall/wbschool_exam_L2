package main

import (
	"flag"
	"fmt"
	"io"
	"man_sort/pkg/sorter"
	"os"
	"time"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

# Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

# Дополнительное

# Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	column := flag.String("k", "", "column that are being sorted")                           // С колонками или с начала строки
	integerValueSort := flag.Bool("n", false, "should it sort by integer value")             // С суфиксами или без них.
	reverseSort := flag.Bool("r", false, "sort in reverse order")                            // В самом конце перевернуть слайс перед записью в файл
	noRepeatedStrings := flag.Bool("u", false, "should not display repeated strings")        // На этапе сканирования файла
	monthNameSort := flag.Bool("M", false, "should compare months")                          // Приоритетнее чем колонки
	ignoreTailSpaces := flag.Bool("b", false, "should ignore spaces in the end of the line") // На этапе сканирования файла
	checkIfSorted := flag.Bool("c", false, "should check if input is already sorted")
	integerAndSuffixValueSort := flag.Bool("h", false, "should sort by integer value including suffixes") // Приоритетнее чем n
	flag.Parse()
	options := sorter.SortingOptions{
		SelectedColumn:     *column,
		IntValueSort:       *integerValueSort,
		IntSuffixValueSort: *integerAndSuffixValueSort,
		ReversedSort:       *reverseSort,
		NoRepeatedString:   *noRepeatedStrings,
		MonthNameSort:      *monthNameSort,
		IgnoreTailSpaces:   *ignoreTailSpaces,
		CheckOnly:          *checkIfSorted,
	}
	fileName := os.Args[len(os.Args)-1]
	file, err := os.Open(fileName)
	if err != nil {
		errStr := fmt.Sprintf("unable to open a file with error: %v. Provided path: %v\n", err.Error(), fileName)
		_, _ = os.Stderr.WriteString(errStr)
		os.Exit(1)
	}

	sortObj := sorter.FileSorter{
		File:    file,
		Options: options,
	}
	sortObj.ScanFile()
	_ = sortObj.File.Close()
	if sortObj.Data == nil || len(sortObj.Data) == 0 {
		_, _ = io.WriteString(os.Stderr, "file is empty\n")
		os.Exit(1)
	}
	sortObj.Sort()
	newFileName := fmt.Sprintf("output-%v.txt", time.Now().Format("2006-01-02-15-04-05"))
	newFile, err := os.Create(newFileName)
	if err != nil {
		errStr := fmt.Sprintf("unable to create a file with error: %v.", err.Error())
		_, _ = os.Stderr.WriteString(errStr)
		os.Exit(1)
	}
	for _, s := range sortObj.Data {
		_, err := newFile.WriteString(s + "\n")
		if err != nil {
			errStr := fmt.Sprintf("unable to write to file with error: %v.", err.Error())
			_, _ = os.Stderr.WriteString(errStr)
			os.Exit(1)
		}
	}
}
