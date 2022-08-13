package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type sliceSort func(fSlice []string, rev bool) func(i int, j int) bool

func readFile(fileName string) ([]string, error) {
	var sc *bufio.Scanner
	if len(fileName) == 0 {
		sc = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		sc = bufio.NewScanner(f)
		defer f.Close()
	}

	var fSlice []string
	// Чтение файла с ридером
	for sc.Scan() {
		fSlice = append(fSlice, sc.Text())
	}

	// handle first encountered error while reading
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return fSlice, nil
}

func writeFile(fileName string, fSlice []string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	for _, s := range fSlice {
		_, err := fmt.Fprintln(f, s)
		if err != nil {
			return err
		}
	}
	f.Close()
	return nil
}

//Строки с цифрами размещаются выше других строк
//Строки, начинающиеся с букв нижнего регистра размещаются выше
//Сортировка выполняется в соответствии c алфавитом
//Строки сначала сортируются по алфавиту, а уже вторично по другим правилам.
func order(i, j string) bool {
	a := strings.ToLower(i)
	b := strings.ToLower(j)
	if a == b {
		return i > j
	}
	return a < b
}

func xor(a, b bool) bool {
	return (a && !b) || (!a && b)
}

func alphabetSort(fSlice []string, rev bool) func(i, j int) bool {
	return func(i, j int) (res bool) {
		defer func() {
			res = xor(res, rev)
		}()
		return order(fSlice[i], fSlice[j])
	}
}

func numSort(fSlice []string, rev bool) func(i, j int) bool {
	return func(i, j int) (res bool) {
		defer func() {
			res = xor(res, rev)
		}()
		a, err1 := strconv.Atoi(fSlice[i])
		b, err2 := strconv.Atoi(fSlice[j])
		if err1 != nil && err2 != nil {
			return order(fSlice[i], fSlice[j])
		}
		if err1 != nil {
			return true
		}
		if err2 != nil {
			return false
		}
		return a < b
	}
}

func getColumnMap(fSlice []string, sep string, k int) (StringMap map[string][]string, keys []string, unsorted []string) {
	StringMap = make(map[string][]string)
	for _, str := range fSlice {
		a := strings.Split(str, sep)
		if len(a)-1 < k {
			unsorted = append(unsorted, str)
			continue
		}
		_, ok := StringMap[a[k]]
		if !ok {
			keys = append(keys, a[k])
		}
		StringMap[a[k]] = append(StringMap[a[k]], str)
	}
	return
}

func Unique(input []string) []string {
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

func sortByColumn(fSlice []string, k int, mainSort, subSort sliceSort, reversal bool) []string {
	StringMap, keys, unsorted := getColumnMap(fSlice, " ", k)
	less := subSort(unsorted, reversal)
	var result []string
	sort.Slice(unsorted, less)
	result = append(result, unsorted...)
	less = mainSort(keys, reversal)
	sort.Slice(keys, less)
	for _, key := range keys {
		less = subSort(StringMap[key], reversal)
		sort.Slice(StringMap[key], less)
		result = append(result, StringMap[key]...)
	}
	return result
}

func sortSort(fSlice []string, k int, unique, numeric, reversal bool) []string {
	subSort := alphabetSort
	mainSort := alphabetSort
	if numeric {
		mainSort = numSort
	}

	if unique {
		fSlice = Unique(fSlice)
	}

	if k >= 0 {
		return sortByColumn(fSlice, k, mainSort, subSort, reversal)
	}

	less := mainSort(fSlice, reversal)
	sort.Slice(fSlice, less)
	return fSlice
}

func main() {
	fName := "test.txt"
	fSlice, err := readFile(fName)
	if err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	k := 6
	uniq := true
	numeric := false
	reverse := false

	fSlice = sortSort(fSlice, k, uniq, numeric, reverse)

	err = writeFile(fName, fSlice)

	if err != nil {
		log.Fatalf("Error while writing file: %s", err)
	}
}
