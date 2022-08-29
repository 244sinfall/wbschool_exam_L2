package cut

import (
	"os"
	"testing"
)

func TestCut_getFields(t *testing.T) {
	c := Cut{
		Source: nil,
		F:      Fields{"NAME", "LASTNAME"},
		D:      "\t",
		S:      false,
	}
	x, err := c.getFields("AGE\tNAME\tLASTNAME")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if len(x) != 2 {
		t.Log(x)
		t.Fail()
	}
	if x[0] != 1 && x[1] != 2 {
		t.Log(x)
		t.Fail()
	}
}

func TestFields_Set(t *testing.T) {
	var f Fields
	err := f.Set("1,2,3,4")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if len(f) != 4 {
		t.Log(f)
		t.Fail()
	}
	if f[2] != "3" {
		t.Log(f)
		t.Fail()
	}
}

func TestFields_String(t *testing.T) {
	f := Fields{"A", "B", "C", "D"}
	str := f.String()
	if str != "A,B,C,D" {
		t.Log(str)
		t.Fail()
	}
}

func TestCut_Write(t *testing.T) {
	f, err := os.Open("test_file1.txt")
	if err != nil {
		t.Log(err.Error())
		t.Fail()

	}
	c := Cut{
		Source: f,
		F:      Fields{"NAME", "LASTNAME"},
		D:      "\t",
		S:      false,
	}
	strs, err := c.Write()
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	if len(strs) != 3 {
		t.Log(strs)
		t.Fail()
	}
	if strs[0] != "NAME\tLASTNAME" {
		t.Log(strs[0])
		t.Fail()
	}
}
