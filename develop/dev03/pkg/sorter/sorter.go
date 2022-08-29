package sorter

import (
	"bufio"
	"fmt"
	"man_sort/internal"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type SortingOptions struct {
	SelectedColumn     string
	IntValueSort       bool
	IntSuffixValueSort bool
	ReversedSort       bool
	NoRepeatedString   bool
	MonthNameSort      bool
	IgnoreTailSpaces   bool
	CheckOnly          bool
}

type FileSorter struct {
	File    *os.File
	Data    []string
	Options SortingOptions
}

func (s *FileSorter) columnSort(columnId int) {
	sort.Slice(s.Data[1:], func(i, j int) bool {
		iColumn := strings.Split(s.Data[i], " ")[columnId]
		jColumn := strings.Split(s.Data[j], " ")[columnId]
		return s.Options.sortingOrder(iColumn < jColumn)
	})
}

func (s *FileSorter) numericSort() {
	numberOnly, err := regexp.Compile("[0-9]+")
	if err != nil {
		_, _ = os.Stderr.WriteString("error on regex compile. Doing default sort\n")
		s.defaultSort()
	}
	sort.Slice(s.Data, func(i, j int) bool {
		in, _ := strconv.Atoi(numberOnly.FindString(s.Data[i]))
		jn, _ := strconv.Atoi(numberOnly.FindString(s.Data[j]))
		return s.Options.sortingOrder(in < jn)
	})
}

func (s *FileSorter) defaultSort() {
	sort.Slice(s.Data, func(i, j int) bool {
		return s.Options.sortingOrder(s.Data[i] < s.Data[j])
	})
}

func (o SortingOptions) sortingOrder(in bool) bool {
	if o.ReversedSort {
		return !in
	}
	return in
}

func (s *FileSorter) checkQuit(sorted bool) {
	if s.Options.CheckOnly {
		_ = s.File.Close()
		if sorted {
			fmt.Println("File is sorted.")
		} else {
			fmt.Println("File is not sorted.")
		}
		os.Exit(0)
	}
}

func (s *FileSorter) ScanFile() {
	scanner := bufio.NewScanner(s.File)
	s.Data = make([]string, 0, 8)
mainLoop:
	for scanner.Scan() {
		text := scanner.Text()
		if s.Options.IgnoreTailSpaces {
			trimmedText := strings.TrimSpace(text)
			if text != trimmedText {
				s.checkQuit(false)
				text = trimmedText
			}
		}
		if s.Options.CheckOnly {
			for _, str := range s.Data {
				if str == text {
					s.checkQuit(false)
					continue mainLoop
				}
			}
		}
		s.Data = append(s.Data, text)
	}
}

func (s *FileSorter) Sort() {
	var oldData = make([]string, len(s.Data))
	copy(oldData, s.Data)
	switch {
	case s.Options.IntSuffixValueSort:
		s.numericSortWithSuffixes()
	case s.Options.IntValueSort:
		s.numericSort()
	case s.Options.MonthNameSort:
		s.monthSort()
	case s.Options.SelectedColumn != "":
		columnId := internal.FindColumnIndex(s.Data[0], s.Options.SelectedColumn)
		s.columnSort(columnId)
	default:
		s.defaultSort()
	}
	if s.Options.CheckOnly {
		for i, v := range oldData {
			if v != s.Data[i] {
				s.checkQuit(false)
			}
		}
		s.checkQuit(true)
	}
}

var months = map[string]int{"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4,
	"MAY": 5, "JUN": 6, "JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12}

func (s *FileSorter) monthSort() {
	sort.Slice(s.Data, func(i, j int) bool {
		var vi, vj int
		for k, v := range months {
			if strings.Contains(s.Data[i], k) {
				vi = v
			}
			if strings.Contains(s.Data[j], k) {
				vj = v
			}
		}
		return s.Options.sortingOrder(vi < vj)
	})
}

func (s *FileSorter) numericSortWithSuffixes() {
	numberWithSuffix, err := regexp.Compile("[0-9]+([KMGT])?")
	if err != nil {
		_, _ = os.Stderr.WriteString("error on regex compile. Doing default sort\n")
		s.defaultSort()
	}
	sort.Slice(s.Data, func(i, j int) bool {
		ins := numberWithSuffix.FindString(s.Data[i])
		jns := numberWithSuffix.FindString(s.Data[j])
		return s.Options.sortingOrder(internal.GetSuffixValue(ins) < internal.GetSuffixValue(jns))
	})
}
