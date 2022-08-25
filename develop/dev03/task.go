package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
type sortingOptions struct {
	selectedColumn     string
	intValueSort       bool
	intSuffixValueSort bool
	reversedSort       bool
	noRepeatedString   bool
	monthNameSort      bool
	ignoreTailSpaces   bool
	checkOnly          bool
}

type fileSorter struct {
	file    *os.File
	data    []string
	options sortingOptions
}

func (s *fileSorter) checkQuit(sorted bool) {
	if s.options.checkOnly {
		_ = s.file.Close()
		if sorted {
			fmt.Println("File is sorted.")
		} else {
			fmt.Println("File is not sorted.")
		}
		os.Exit(0)
	}
}

func (s *fileSorter) scanFile() {
	scanner := bufio.NewScanner(s.file)
	s.data = make([]string, 0, 8)
mainLoop:
	for scanner.Scan() {
		text := scanner.Text()
		if s.options.ignoreTailSpaces {
			trimmedText := strings.TrimSpace(text)
			if text != trimmedText {
				s.checkQuit(false)
				text = trimmedText
			}
		}
		if s.options.checkOnly {
			for _, str := range s.data {
				if str == text {
					s.checkQuit(false)
					continue mainLoop
				}
			}
		}
		s.data = append(s.data, text)
	}
}

func (s *fileSorter) sort() {
	var oldData = make([]string, len(s.data))
	copy(oldData, s.data)
	switch {
	case s.options.intSuffixValueSort:
		s.numericSortWithSuffixes()
	case s.options.intValueSort:
		s.numericSort()
	case s.options.monthNameSort:
		s.monthSort()
	case s.options.selectedColumn != "":
		columnId := findColumnIndex(s.data[0], s.options.selectedColumn)
		s.columnSort(columnId)
	default:
		s.defaultSort()
	}
	if s.options.checkOnly {
		for i, v := range oldData {
			if v != s.data[i] {
				s.checkQuit(false)
			}
		}
		s.checkQuit(true)
	}
}

func findColumnIndex(mainColumns string, requestedColumn string) int {
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

var months = map[string]int{"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4,
	"MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12}

func (s *fileSorter) monthSort() {
	sort.Slice(s.data, func(i, j int) bool {
		var vi, vj int
		for k, v := range months {
			if strings.Contains(s.data[i], k) {
				vi = v
			}
			if strings.Contains(s.data[j], k) {
				vj = v
			}
		}
		return s.options.sortingOrder(vi < vj)
	})
}

var suffixes = map[string]int{"K": 1000, "M": 1_000_000, "G": 1_000_000_000, "T": 1_000_000_000_000}

func (s *fileSorter) numericSortWithSuffixes() {
	numberWithSuffix, err := regexp.Compile("[0-9]+([KMGT])?")
	if err != nil {
		_, _ = os.Stderr.WriteString("error on regex compile. Doing default sort\n")
		s.defaultSort()
	}
	sort.Slice(s.data, func(i, j int) bool {
		ins := numberWithSuffix.FindString(s.data[i])
		jns := numberWithSuffix.FindString(s.data[j])
		return s.options.sortingOrder(getSuffixValue(ins) < getSuffixValue(jns))
	})
}

func getSuffixValue(number string) int {
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

func (s *fileSorter) columnSort(columnId int) {
	sort.Slice(s.data[1:], func(i, j int) bool {
		iColumn := strings.Split(s.data[i], " ")[columnId]
		jColumn := strings.Split(s.data[j], " ")[columnId]
		return s.options.sortingOrder(iColumn < jColumn)
	})
}

func (s *fileSorter) numericSort() {
	numberOnly, err := regexp.Compile("[0-9]+")
	if err != nil {
		_, _ = os.Stderr.WriteString("error on regex compile. Doing default sort\n")
		s.defaultSort()
	}
	sort.Slice(s.data, func(i, j int) bool {
		in, _ := strconv.Atoi(numberOnly.FindString(s.data[i]))
		jn, _ := strconv.Atoi(numberOnly.FindString(s.data[j]))
		return s.options.sortingOrder(in < jn)
	})
}

func (s *fileSorter) defaultSort() {
	sort.Slice(s.data, func(i, j int) bool {
		return s.options.sortingOrder(s.data[i] < s.data[j])
	})
}

func (o sortingOptions) sortingOrder(in bool) bool {
	if o.reversedSort {
		return !in
	}
	return in
}

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
	options := sortingOptions{
		selectedColumn:     *column,
		intValueSort:       *integerValueSort,
		intSuffixValueSort: *integerAndSuffixValueSort,
		reversedSort:       *reverseSort,
		noRepeatedString:   *noRepeatedStrings,
		monthNameSort:      *monthNameSort,
		ignoreTailSpaces:   *ignoreTailSpaces,
		checkOnly:          *checkIfSorted,
	}
	fileName := os.Args[len(os.Args)-1]
	file, err := os.Open(fileName)
	if err != nil {
		errStr := fmt.Sprintf("unable to open a file with error: %v. Provided path: %v\n", err.Error(), fileName)
		_, _ = os.Stderr.WriteString(errStr)
		os.Exit(1)
	}

	sorter := fileSorter{
		file:    file,
		options: options,
	}
	sorter.scanFile()
	_ = sorter.file.Close()
	if sorter.data == nil || len(sorter.data) == 0 {
		_, _ = io.WriteString(os.Stderr, "file is empty\n")
		os.Exit(1)
	}
	sorter.sort()
	newFileName := fmt.Sprintf("output-%v.txt", time.Now().Format("2006-01-02-15-04-05"))
	newFile, err := os.Create(newFileName)
	if err != nil {
		errStr := fmt.Sprintf("unable to create a file with error: %v.", err.Error())
		_, _ = os.Stderr.WriteString(errStr)
		os.Exit(1)
	}
	for _, s := range sorter.data {
		_, err := newFile.WriteString(s + "\n")
		if err != nil {
			errStr := fmt.Sprintf("unable to write to file with error: %v.", err.Error())
			_, _ = os.Stderr.WriteString(errStr)
			os.Exit(1)
		}
	}
}
