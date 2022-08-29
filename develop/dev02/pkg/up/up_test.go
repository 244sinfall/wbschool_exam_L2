package up

import "testing"

func TestUnpack(t *testing.T) {
	output, err := Unpack(`a4bc2d5e`)
	if output != "aaaabccddddde" || err != nil {
		t.Logf("a4bc2d5e => %v", output)
		t.Fail()
	}
	output, err = Unpack(`abcd`)
	if output != "abcd" || err != nil {
		t.Logf("abcd => %v", output)
		t.Fail()
	}
	output, err = Unpack(`45`)
	if output != "" || err == nil {
		t.Logf("45 => %v", output)
		t.Fail()
	}
	output, err = Unpack(``)
	if output != "" || err == nil {
		t.Logf(" => %v", output)
		t.Fail()
	}
	output, err = Unpack(`qwe\4\5`)
	if output != "qwe45" || err != nil {
		t.Logf(`qwe\4\5 => %v`, output)
		t.Fail()
	}
	output, err = Unpack(`qwe\45`)
	if output != "qwe44444" || err != nil {
		t.Logf(`qwe\45 => %v`, output)
		t.Fail()
	}
	output, err = Unpack(`qwe\\5`)
	if output != `qwe\\\\\` || err != nil {
		t.Logf(`qwe\\5 => %v`, output)
		t.Fail()
	}
}
