package main

import (
	"strings"
	"testing"
)

func TestGetAmount(t *testing.T) {
	str := strings.Builder{}
	str.WriteString("1234")
	amount := getAmount(&str)
	if amount != 1234 {
		t.Logf("1234 => %v", amount)
		t.Fail()
	}
	str.Reset()
	str.WriteString("")
	amount = getAmount(&str)
	if amount != 1 {
		t.Fail()
		t.Logf(" => %v", amount)
	}
	str.Reset()
	str.WriteString("-24")
	amount = getAmount(&str)
	if amount != -24 {
		t.Logf("-24 => %v", amount)
		t.Fail()
	}
	str.Reset()
	str.WriteString("11134")
	amount = getAmount(&str)
	if amount != 11134 {
		t.Logf("11134 => %v", amount)
		t.Fail()
	}
	str.Reset()
	str.WriteString("qwe")
	amount = getAmount(&str)
	if amount != 0 {
		t.Logf("qwe => %v", amount)
		t.Fail()
	}
	str.Reset()
}

func TestWriteSequence(t *testing.T) {
	output := strings.Builder{}
	writeSequence('r', 10, &output)
	if output.String() != "rrrrrrrrrr" {
		t.Fail()
	}
	output.Reset()
	writeSequence('a', 1, &output)
	if output.String() != "a" {
		t.Fail()
	}
	output.Reset()
	writeSequence('z', 0, &output)
	if output.String() != "" {
		t.Fail()
	}
	output.Reset()
	writeSequence('z', -1, &output)
	if output.String() != "" {
		t.Fail()
	}
	output.Reset()
}

func TestUnpack(t *testing.T) {
	output, err := unpack(`a4bc2d5e`)
	if output != "aaaabccddddde" || err != nil {
		t.Logf("a4bc2d5e => %v", output)
		t.Fail()
	}
	output, err = unpack(`abcd`)
	if output != "abcd" || err != nil {
		t.Logf("abcd => %v", output)
		t.Fail()
	}
	output, err = unpack(`45`)
	if output != "" || err == nil {
		t.Logf("45 => %v", output)
		t.Fail()
	}
	output, err = unpack(``)
	if output != "" || err == nil {
		t.Logf(" => %v", output)
		t.Fail()
	}
	output, err = unpack(`qwe\4\5`)
	if output != "qwe45" || err != nil {
		t.Logf(`qwe\4\5 => %v`, output)
		t.Fail()
	}
	output, err = unpack(`qwe\45`)
	if output != "qwe44444" || err != nil {
		t.Logf(`qwe\45 => %v`, output)
		t.Fail()
	}
	output, err = unpack(`qwe\\5`)
	if output != `qwe\\\\\` || err != nil {
		t.Logf(`qwe\\5 => %v`, output)
		t.Fail()
	}
}
