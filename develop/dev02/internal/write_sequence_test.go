package internal

import (
	"strings"
	"testing"
)

func TestWriteSequence(t *testing.T) {
	output := strings.Builder{}
	WriteSequence('r', 10, &output)
	if output.String() != "rrrrrrrrrr" {
		t.Fail()
	}
	output.Reset()
	WriteSequence('a', 1, &output)
	if output.String() != "a" {
		t.Fail()
	}
	output.Reset()
	WriteSequence('z', 0, &output)
	if output.String() != "" {
		t.Fail()
	}
	output.Reset()
	WriteSequence('z', -1, &output)
	if output.String() != "" {
		t.Fail()
	}
	output.Reset()
}
