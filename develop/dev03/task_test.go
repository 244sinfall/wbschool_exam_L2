package main

import (
	"os"
	"testing"
)

var testingFile, _ = os.Open("go_test_file.txt")
var testObj = fileSorter{
	file: testingFile,
	data: nil,
	options: sortingOptions{
		selectedColumn:     "",
		intValueSort:       false,
		intSuffixValueSort: false,
		reversedSort:       false,
		noRepeatedString:   false,
		monthNameSort:      false,
		ignoreTailSpaces:   false,
		checkOnly:          false,
	},
}

func TestDefaultSortAndReverse(t *testing.T) {
	testObj.data = []string{"C", "D", "A"}
	testObj.defaultSort()
	if testObj.data[0] != "A" {
		t.Fail()
	}
	testObj.options.reversedSort = true
	testObj.defaultSort()
	if testObj.data[0] != "D" {
		t.Fail()
	}
	testObj.options.reversedSort = false
}

func TestColumnSort(t *testing.T) {
	testObj.data = []string{"A B C", "1 2 3", "4 1 6"}
	testObj.columnSort(1)
	if testObj.data[1] != "4 1 6" {
		t.Fail()
	}
}

func TestMonthSort(t *testing.T) {
	testObj.data = []string{"JAN", "DEC", "NOV"}
	testObj.monthSort()
	if testObj.data[1] != "NOV" {
		t.Fail()
	}
}

func TestNumericSort(t *testing.T) {
	testObj.data = []string{"15", "3", "56"}
	testObj.numericSort()
	if testObj.data[0] != "3" {
		t.Fail()
	}
}

func TestNumericAndSuffix(t *testing.T) {
	testObj.data = []string{"999M", "1K"}
	testObj.numericSortWithSuffixes()
	if testObj.data[0] != "1K" {
		t.Fail()
	}
}
