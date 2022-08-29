package main

import (
	"flag"
	"os"
	"regexp"
	"strings"
	"task5/pkg/grep"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	options := grep.Options{
		After:           0,
		Before:          0,
		Context:         0,
		LinesOnly:       false,
		IgnoreCase:      false,
		Invert:          false,
		Fixed:           false,
		PrintLineNumber: false,
	}
	flag.IntVar(&options.After, "A", 0, "print N strings after match")
	flag.IntVar(&options.Before, "B", 0, "print N strings before match")
	flag.IntVar(&options.Context, "C", 0, "print N strings before and after match")                     // Если есть A и B - прибавляем их
	flag.BoolVar(&options.LinesOnly, "c", false, "display how many lines match search")                 // Только количество совпадений
	flag.BoolVar(&options.IgnoreCase, "i", false, "ignore match / input case")                          // будем приводить все к нижнему
	flag.BoolVar(&options.Invert, "v", false, "print anything except matching line(s)")                 // простой свич на цикл
	flag.BoolVar(&options.Fixed, "F", false, "search only exact match, not pattern")                    // сравнение срок или сравнение мэтча регекс со всей строкой
	flag.BoolVar(&options.PrintLineNumber, "n", false, "also print number of line that found match on") // в output достаем lineNumber из match
	flag.Parse()
	args := flag.Args() // первый аргумент выражение/строка, второй - файл
	if len(args) != 2 {
		_, _ = os.Stderr.WriteString("Not enough arguments. Use regex and filepath\n")
		os.Exit(1)
	}
	f, err := os.Open(args[1])
	if err != nil {
		_, _ = os.Stderr.WriteString("Error opening file: " + err.Error() + "\n")
		os.Exit(1)
	}
	g := grep.Grep{
		File:    f,
		Options: options,
		Exp:     args[0],
		Regex:   nil,
		Matches: make([]grep.Match, 0, 1+options.After+options.Before+options.Context),
	}
	r, err := regexp.Compile(g.Exp)
	if err == nil {
		g.Regex = r
	} else if g.Options.IgnoreCase {
		g.Exp = strings.ToLower(g.Exp)
	}
	g.Execute()     // Первичный поиск совпадений и перенос файла в слайса
	g.PrintResult() // Напечатать результаты
	os.Exit(0)
}
