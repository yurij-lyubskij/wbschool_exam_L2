package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type myrunes struct {
	Runes []rune
}

func (m *myrunes) less(i, j int) bool {
	return m.Runes[i] < m.Runes[j]
}

//sortRunes - функция, которая сортирует
//руны в строке
func sortRunes(s string) string {
	Myrunes := myrunes{Runes: []rune(s)}
	sort.Slice(Myrunes.Runes, Myrunes.less)
	return string(Myrunes.Runes)
}

type mystrings struct {
	Strings []string
}

func (m *mystrings) less(i, j int) bool {
	return m.Strings[i] < m.Strings[j]
}

//sortStrings - функция, которая сортирует
//строки в слайсе
func sortStrings(s []string) *[]string {
	MyStrings := mystrings{Strings: s}
	sort.Slice(MyStrings.Strings, MyStrings.less)
	return &MyStrings.Strings
}

//unique - функция, которая убирает
//повторяющиеся строки
func unique(input []string) []string {
	StringSet := make(map[string]struct{})
	for _, str := range input {
		_, ok := StringSet[str]
		if !ok {
			StringSet[str] = struct{}{}
		}
	}
	uniq := make([]string, len(StringSet))
	i := 0
	for str := range StringSet {
		uniq[i] = str
		i++
	}
	return uniq
}

//FindInitial - функция, которая находит анаграммы
//и записывает их в мапу в порядке нахождения слов
//слова приводятся к нижнему регистру
//ключом данной мапы является строка,
//в которой руны отсортированы по возрастанию
func FindInitial(words *[]string) map[string][]string {
	result := make(map[string][]string)
	for _, word := range *words {
		lowerWord := strings.ToLower(word)
		key := sortRunes(lowerWord)
		result[key] = append(result[key], lowerWord)
	}
	return result
}

//FindAnagrams - функция, которая находит
//и обрабатывает анаграммы в словаре
//на выходе - мапа, ключом которой
//является первое встреченное слово из множества.
//слова в множестве не повторяются
//и упорядочены по алфавиту
func FindAnagrams(words *[]string) *map[string]*[]string {
	initialMap := FindInitial(words)
	result := make(map[string]*[]string)
	for _, value := range initialMap {
		uniqAnagrams := unique(value)
		if len(uniqAnagrams) == 1 {
			continue
		}
		key := uniqAnagrams[0]
		result[key] = sortStrings(uniqAnagrams)
	}
	return &result
}

func main() {
	words := &[]string{"тяпка", "катяп", "пятак", "ПЯТКА", "пятка", "отдельно", "листок", "слиток", "столик"}
	fmt.Println("Исходный словарь:")
	fmt.Printf("%v\n", *words)
	anagrams := FindAnagrams(words)
	fmt.Println("Найденные анаграммы:")
	for key, value := range *anagrams {
		fmt.Printf(key)
		fmt.Printf("%v\n", *value)
	}
}
