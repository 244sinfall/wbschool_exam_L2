package internal

import (
	"strings"
	"testing"
)

func TestGetAmount(t *testing.T) {
	str := strings.Builder{}
	str.WriteString("1234")
	amount := GetAmount(&str)
	if amount != 1234 {
		t.Logf("1234 => %v", amount)
		t.Fail()
	}
	str.Reset()
	str.WriteString("")
	amount = GetAmount(&str)
	if amount != 1 {
		t.Fail()
		t.Logf(" => %v", amount)
	}
	str.Reset()
	str.WriteString("-24")
	amount = GetAmount(&str)
	if amount != -24 {
		t.Logf("-24 => %v", amount)
		t.Fail()
	}
	str.Reset()
	str.WriteString("11134")
	amount = GetAmount(&str)
	if amount != 11134 {
		t.Logf("11134 => %v", amount)
		t.Fail()
	}
	str.Reset()
	str.WriteString("qwe")
	amount = GetAmount(&str)
	if amount != 0 {
		t.Logf("qwe => %v", amount)
		t.Fail()
	}
	str.Reset()
}
