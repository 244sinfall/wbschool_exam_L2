package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
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

type grep struct {
	file    *os.File
	data    []string
	options grepOptions
	exp     string
	regex   *regexp.Regexp
	matches []match
}

type match struct {
	line       string
	lineNumber int
}

type grepOptions struct {
	after           int
	before          int
	context         int
	linesOnly       bool
	ignoreCase      bool
	invert          bool
	fixed           bool
	printLineNumber bool
}

func (g *grep) addMatch(m match) {
	g.matches = append(g.matches, m)
}

func (g *grep) checkLine(s string, n int) {
	l := s
	if g.options.ignoreCase {
		l = strings.ToLower(s)
	}
	if g.regex != nil { // В регулярке бессмысленен флаг ignoreCase, поскольку такая возможность вложена во входящий паттерн
		m := g.regex.FindString(l)
		if m != "" && (!g.options.fixed || m == strings.TrimSpace(l)) { // Если не точное совпадение, или паттерн полностью совпадает со строкой
			g.addMatch(match{
				line:       "Match -> " + s,
				lineNumber: n,
			})
		}
	} else {
		if strings.Contains(l, g.exp) {
			if !g.options.fixed || g.exp == strings.TrimSpace(l) {
				g.addMatch(match{
					line:       "Match -> " + s,
					lineNumber: n,
				})
			}
		}
	}
}

func (g *grep) execute() {
	scanner := bufio.NewScanner(g.file)
	fileStats, err := g.file.Stat()
	if err != nil {
		_, _ = os.Stderr.WriteString("Error while analyzing the file: " + err.Error() + "\n")
		os.Exit(1)
	}
	fileSize := fileStats.Size() // Возвращает размер файла в байтах.
	// Предположим, что один символ строки будет в 2 байта, получим количество символов.
	// Предположим, что в строке 80 символов. Это будет размером слайса.
	sliceSize := 2 * fileSize / 80
	if sliceSize < 1 {
		sliceSize = 1
	}
	g.data = make([]string, 0, sliceSize)
	// Заполнение данных для дальнейшего поиска по индексу + ищем первичные входы
	var lineNumber int
	for scanner.Scan() {
		text := scanner.Text()
		g.checkLine(text, lineNumber)
		g.data = append(g.data, text)
		lineNumber++
	}
}

func (g *grep) printLine(i int) {
	if g.options.printLineNumber {
		prefix := fmt.Sprintf("Line %d: ", i+1)
		fmt.Println(prefix, g.data[i])
	} else {
		fmt.Println(g.data[i])
	}
}

func (g *grep) printMatch(m match) {
	if g.options.printLineNumber {
		prefix := fmt.Sprintf("Line %d: ", m.lineNumber+1)
		fmt.Println(prefix, m.line)
	} else {
		fmt.Println(m.line)
	}
}

func (g *grep) printMatches() {
	currentIndex := 0
	for _, m := range g.matches {
		leftOffset := g.options.context + g.options.before
		startIndex := m.lineNumber - leftOffset
		if startIndex < 0 {
			startIndex = 0
		}
		rightOffset := g.options.context + g.options.after
		endIndex := m.lineNumber + rightOffset
		if endIndex > len(g.data)-1 {
			endIndex = len(g.data) - 1
		}
		if !g.options.invert {
			for i := startIndex; i < m.lineNumber; i++ {
				g.printLine(i)
			}
			g.printMatch(m)
			for i := m.lineNumber + 1; i <= endIndex; i++ {
				g.printLine(i)
			}
		} else {
			for i := currentIndex; i < startIndex; i++ {
				g.printLine(i)
			}
			currentIndex = endIndex + 1
		}
	}
}

func (g *grep) printResult() {
	if len(g.matches) == 0 {
		fmt.Println("No matches!")
		return
	}
	if g.options.linesOnly {
		fmt.Printf("Found %d lines\n", len(g.matches))
		return
	}
	g.printMatches()
}

func main() {
	options := grepOptions{
		after:           0,
		before:          0,
		context:         0,
		linesOnly:       false,
		ignoreCase:      false,
		invert:          false,
		fixed:           false,
		printLineNumber: false,
	}
	flag.IntVar(&options.after, "A", 0, "print N strings after match")
	flag.IntVar(&options.before, "B", 0, "print N strings before match")
	flag.IntVar(&options.context, "C", 0, "print N strings before and after match")                     // Если есть A и B - прибавляем их
	flag.BoolVar(&options.linesOnly, "c", false, "display how many lines match search")                 // Только количество совпадений
	flag.BoolVar(&options.ignoreCase, "i", false, "ignore match / input case")                          // будем приводить все к нижнему
	flag.BoolVar(&options.invert, "v", false, "print anything except matching line(s)")                 // простой свич на цикл
	flag.BoolVar(&options.fixed, "F", false, "search only exact match, not pattern")                    // сравнение срок или сравнение мэтча регекс со всей строкой
	flag.BoolVar(&options.printLineNumber, "n", false, "also print number of line that found match on") // в output достаем lineNumber из match
	flag.Parse()
	args := flag.Args() // первый аргумент выражение/строка, второй - файл
	fmt.Println(args)
	if len(args) != 2 {
		_, _ = os.Stderr.WriteString("Not enough arguments. Use regex and filepath\n")
		os.Exit(1)
	}
	f, err := os.Open(args[1])
	if err != nil {
		_, _ = os.Stderr.WriteString("Error opening file: " + err.Error() + "\n")
		os.Exit(1)
	}
	g := grep{
		file:    f,
		options: options,
		exp:     args[0],
		regex:   nil,
		matches: make([]match, 0, 1+options.after+options.before+options.context),
	}
	r, err := regexp.Compile(g.exp)
	if err == nil {
		g.regex = r
	} else if g.options.ignoreCase {
		g.exp = strings.ToLower(g.exp)
	}
	g.execute()     // Первичный поиск совпадений и перенос файла в слайса
	g.printResult() // Напечатать результаты
	os.Exit(0)
}
