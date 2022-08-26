package main

import "testing"

func TestByteSumm(t *testing.T) {
	word1 := "привет"
	word2 := "тевирп"
	if byteSumFor(word1) != byteSumFor(word2) {
		t.Fail()
	}
	word3 := "прывет"
	if byteSumFor(word1) == byteSumFor(word3) {
		t.Fail()
	}
}

func TestAnagrams(t *testing.T) {
	var testArr = []string{"пятак", "пятка", "тяпка", "Капят", // Все ок
		"монетка", "катеном", // Два элемента
		"СЫР", // 1 элемент
		"листок", "слиток", "столик"}
	ang := showAnagrams(testArr)
	for k := range ang {
		byteSumFor(k)
		for _, v := range ang[k] {
			if byteSumFor(v) != byteSumFor(k) {
				t.Fail()
			}
		}
	}
	if _, ok := ang["сыр"]; ok {
		t.Fail()
	}
}
