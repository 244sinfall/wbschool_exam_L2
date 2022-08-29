package grep

import (
	"os"
	"testing"
)

var gO = Options{
	After:           0,
	Before:          0,
	Context:         0,
	LinesOnly:       false,
	IgnoreCase:      false,
	Invert:          false,
	Fixed:           false,
	PrintLineNumber: false,
}

var g = Grep{
	File:    nil,
	Options: gO,
	Exp:     "donec",
	Regex:   nil,
	Matches: make([]Match, 0, 1+gO.After+gO.Before+gO.Context),
}

func TestGrep_Execute(t *testing.T) {
	g.File, _ = os.Open("test_file.txt")
	g.Options.IgnoreCase = true
	g.Execute()
	t.Log(g.Matches)
	if len(g.Matches) != 4 {
		t.Fail()
	}
}

func TestGrep_Execute2(t *testing.T) {
	g.Options.IgnoreCase = false
	g.Matches = []Match{}
	g.Execute()
	t.Log(g.Matches)
	if len(g.Matches) != 0 {
		t.Fail()
	}
}
