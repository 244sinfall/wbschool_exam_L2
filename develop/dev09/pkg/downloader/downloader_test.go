package downloader

import (
	"math/rand"
	"os"
	"testing"
)

func TestSeparateLinkToSlice(t *testing.T) {
	site1 := "https://www.google.com"
	site2 := "https://eu.shop.battle.net/ru-ru"
	site3 := "asdasd"
	sep, err := SeparateLinkToSlice(site1)
	if err != nil || sep[0] != "https" || sep[1] != "www.google.com" {
		t.Log(sep)
		t.Fail()
	}
	sep, err = SeparateLinkToSlice(site2)
	if err != nil || sep[0] != "https" || sep[1] != "eu.shop.battle.net" || sep[2] != "ru-ru" {
		t.Log(sep)
		t.Fail()
	}
	_, err = SeparateLinkToSlice(site3)
	if err == nil {
		t.Log(sep)
		t.Fail()
	}
}

func TestValidateLink(t *testing.T) {
	r := rand.Intn(100)
	link1 := []string{"testsite.com", "validate", string(rune(r))}
	ValidateLink(link1)
	f, err := os.Open("testsite.com\\validate\\" + string(rune(r)) + "\\index.html")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	_ = f.Close()
	_ = os.RemoveAll("testsite.com\\validate")
	link2 := []string{"testsite.com", "validate", "main.js"}
	ValidateLink(link2)
	f, err = os.Open("testsite.com\\validate\\main.js")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	_ = f.Close()
	err = os.RemoveAll("testsite.com\\validate")
	if err != nil {
		t.Log()
	}

}

func TestValidatePath(t *testing.T) {
	path1 := "testsite.com"
	path2 := "testsite.com\\index.html"
	path3 := "testsite.com\\main.js"
	err := ValidatePath(&path1)
	if err != nil || path1 != "testsite.com\\index.html" {
		t.Fail()
		t.Log(path1, err)
	}
	err = ValidatePath(&path2)
	if err != nil || path1 != "testsite.com\\index.html" {
		t.Fail()
		t.Log(path2, err)
	}
	err = ValidatePath(&path3)
	if err == nil {
		t.Fail()
		t.Log(path3)
	}

}
