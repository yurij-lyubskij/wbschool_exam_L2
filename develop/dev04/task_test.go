package main

import (
	"testing"
)

func TestOne(t *testing.T) {
	words := &[]string{"пион", "цветок"}
	anagrams := FindAnagrams(words)
	if len(*anagrams) > 0 {
		t.Fatalf("не должно быть множеств из 1 элемента")
	}
}

func TestUnique(t *testing.T) {
	words := &[]string{"пион", "Пион", "онип"}
	anagrams := FindAnagrams(words)
	if len(*anagrams) != 1 {
		t.Fatalf("анаграмы найдены неверно")
	}
	for _, value := range *anagrams {
		if len(*value) != 2 {
			t.Fatalf("не должно быть повторов")
		}
	}
}

func TestFull(t *testing.T) {
	words := &[]string{"тяпка", "катяп", "пятак", "ПЯТКА", "пятка", "пятка", "отдельно", "листок", "слиток", "столик"}
	anagrams := FindAnagrams(words)
	expectedKeys := []string{"тяпка", "листок"}
	for key, value := range *anagrams {
		if key != expectedKeys[0] && key != expectedKeys[1] {
			t.Fatalf("ключ выбран неправильно")
		}
		strSorted := sortRunes(key)
		for _, word := range *value {
			if sortRunes(word) != strSorted {
				t.Fatalf("не анаграмма")
			}
		}
	}
}
