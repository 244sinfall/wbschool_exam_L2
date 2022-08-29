package sorter

import (
	"os"
	"testing"
)

var testingFile, _ = os.Open("go_test_file.txt")
var testObj = FileSorter{
	File: testingFile,
	Data: nil,
	Options: SortingOptions{
		SelectedColumn:     "",
		IntValueSort:       false,
		IntSuffixValueSort: false,
		ReversedSort:       false,
		NoRepeatedString:   false,
		MonthNameSort:      false,
		IgnoreTailSpaces:   false,
		CheckOnly:          false,
	},
}

func TestDefaultSortAndReverse(t *testing.T) {
	testObj.Data = []string{"C", "D", "A"}
	testObj.defaultSort()
	if testObj.Data[0] != "A" {
		t.Fail()
	}
	testObj.Options.ReversedSort = true
	testObj.defaultSort()
	if testObj.Data[0] != "D" {
		t.Fail()
	}
	testObj.Options.ReversedSort = false
}

func TestColumnSort(t *testing.T) {
	testObj.Data = []string{"A B C", "1 2 3", "4 1 6"}
	testObj.columnSort(1)
	if testObj.Data[1] != "4 1 6" {
		t.Fail()
	}
}

func TestMonthSort(t *testing.T) {
	testObj.Data = []string{"JAN", "DEC", "NOV"}
	testObj.monthSort()
	if testObj.Data[1] != "NOV" {
		t.Fail()
	}
}

func TestNumericSort(t *testing.T) {
	testObj.Data = []string{"15", "3", "56"}
	testObj.numericSort()
	if testObj.Data[0] != "3" {
		t.Fail()
	}
}

func TestNumericAndSuffix(t *testing.T) {
	testObj.Data = []string{"999M", "1K"}
	testObj.numericSortWithSuffixes()
	if testObj.Data[0] != "1K" {
		t.Fail()
	}
}
