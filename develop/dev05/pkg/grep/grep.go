package grep

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Grep struct {
	File    *os.File
	data    []string
	Options Options
	Exp     string
	Regex   *regexp.Regexp
	Matches []Match
}

type Match struct {
	line       string
	lineNumber int
}

type Options struct {
	After           int
	Before          int
	Context         int
	LinesOnly       bool
	IgnoreCase      bool
	Invert          bool
	Fixed           bool
	PrintLineNumber bool
}

func (g *Grep) addMatch(m Match) {
	g.Matches = append(g.Matches, m)
}

func (g *Grep) checkLine(s string, n int) {
	l := s
	if g.Options.IgnoreCase {
		l = strings.ToLower(s)
	}
	if g.Regex != nil { // В регулярке бессмысленен флаг ignoreCase, поскольку такая возможность вложена во входящий паттерн
		m := g.Regex.FindString(l)
		if m != "" && (!g.Options.Fixed || m == strings.TrimSpace(l)) { // Если не точное совпадение, или паттерн полностью совпадает со строкой
			g.addMatch(Match{
				line:       "Match -> " + s,
				lineNumber: n,
			})
		}
	} else {
		if strings.Contains(l, g.Exp) {
			if !g.Options.Fixed || g.Exp == strings.TrimSpace(l) {
				g.addMatch(Match{
					line:       "Match -> " + s,
					lineNumber: n,
				})
			}
		}
	}
}

func (g *Grep) Execute() {
	scanner := bufio.NewScanner(g.File)
	fileStats, err := g.File.Stat()
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

func (g *Grep) printLine(i int) {
	if g.Options.PrintLineNumber {
		prefix := fmt.Sprintf("Line %d: ", i+1)
		fmt.Println(prefix, g.data[i])
	} else {
		fmt.Println(g.data[i])
	}
}

func (g *Grep) printMatch(m Match) {
	if g.Options.PrintLineNumber {
		prefix := fmt.Sprintf("Line %d: ", m.lineNumber+1)
		fmt.Println(prefix, m.line)
	} else {
		fmt.Println(m.line)
	}
}

func (g *Grep) printMatches() {
	currentIndex := 0
	for _, m := range g.Matches {
		leftOffset := g.Options.Context + g.Options.Before
		startIndex := m.lineNumber - leftOffset
		if startIndex < 0 {
			startIndex = 0
		}
		rightOffset := g.Options.Context + g.Options.After
		endIndex := m.lineNumber + rightOffset
		if endIndex > len(g.data)-1 {
			endIndex = len(g.data) - 1
		}
		if !g.Options.Invert {
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

func (g *Grep) PrintResult() {
	if len(g.Matches) == 0 {
		fmt.Println("No matches!")
		return
	}
	if g.Options.LinesOnly {
		fmt.Printf("Found %d lines\n", len(g.Matches))
		return
	}
	g.printMatches()
}
