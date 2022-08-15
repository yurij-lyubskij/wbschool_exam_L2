package main

import (
	"bufio"
	"flag"
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

var usageStr = `
Usage: sort [options] <input file> <output file>
output file - stdout по умолчанию
Options:
-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

// функция, которая возвращает
//функцию less для slice.Sort для конкретного слайса
//сортирует в нужном порядке, при передаче флага - в обратном порядке
type sliceSort func(fSlice []string, rev bool) func(i int, j int) bool

//ReadFile - читает файл по передаваемому имени,
//возвращает слайс строк и ошибку
//если имя пустое, читает из stdin
func ReadFile(fileName string) ([]string, error) {
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

//WriteFile - пишет файл по передаваемому имени,
//возвращает ошибку при необходимости
//если имя пустое, пишет в Stdout
func WriteFile(fileName string, fSlice []string) error {
	var f *os.File
	var err error
	if len(fileName) == 0 {
		f = os.Stdout
	} else {
		f, err = os.Create(fileName)
		if err != nil {
			return err
		}
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

//alphabetSort - функция, которая возвращает
//функцию less для slice.Sort для конкретного слайса
//сортирует по алфавиту, при передаче флага - в обратном порядке
func alphabetSort(fSlice []string, rev bool) func(i, j int) bool {
	return func(i, j int) (res bool) {
		defer func() {
			if rev {
				res = !res
			}
		}()
		return order(fSlice[i], fSlice[j])
	}
}

// numSort- функция, которая возвращает
//функцию less для slice.Sort для конкретного слайса
//сортирует по величине числа (numeric),
//при передаче флага - в обратном порядке
func numSort(fSlice []string, rev bool) func(i, j int) bool {
	return func(i, j int) (res bool) {
		defer func() {
			if rev {
				res = !res
			}
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

//getColumnMap - функция, которая принимет
//разделитель и номер колонки
//разделяет строки на колонки,
//выделяет нужную колонку
//возвращает слайс с нужной колонкой
//мапу, где ключ - нужная колонка,
//значение - исходная строка
//и слайс строк, где нет нужной колонки
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

//sortByColumn - функция, которая сортирует
//строки по значению в нужной колонке
//значения, где нет требуемой колонки
//идут в начале, по значению в нужной
//колонке сортируем функцией mainSort
//при совпадении ключей или отсутствии колонки
//сортируем функцией subSort
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

//Sort - функция, которая принимает слайс
//и структуру с параметрами сортировки
//и возвращает отсортированный слайс
func Sort(fSlice []string, keys Keys) []string {
	subSort := alphabetSort
	mainSort := alphabetSort
	if keys.Numeric {
		mainSort = numSort
	}

	if keys.Unique {
		fSlice = unique(fSlice)
	}

	if keys.Column >= 0 {
		return sortByColumn(fSlice, keys.Column, mainSort, subSort, keys.Reverse)
	}

	less := mainSort(fSlice, keys.Reverse)
	sort.Slice(fSlice, less)
	return fSlice
}

//Keys - настройки сортировки
type Keys struct {
	Column  int
	Numeric bool
	Reverse bool
	Unique  bool
}

func main() {
	var keys Keys
	//настройки сортировки - возможные флаги
	flag.IntVar(&keys.Column, "k", -1, "column to sort by")
	flag.BoolVar(&keys.Numeric, "n", false, "numerical sort")
	flag.BoolVar(&keys.Reverse, "r", false, "reverse order")
	flag.BoolVar(&keys.Unique, "u", false, "unique lines")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
		return
	}
	//входной файл
	finName := args[0]
	//выходной файл
	var foutName string
	if len(args) > 1 {
		foutName = args[1]
	}
	//читаем входной файл
	fSlice, err := ReadFile(finName)
	if err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	if keys.Column != -1 {
		keys.Column--
	}
	//сортируем
	fSlice = Sort(fSlice, keys)
	//пишем выходной файл
	err = WriteFile(foutName, fSlice)
	if err != nil {
		log.Fatalf("Error while writing file: %s", err)
	}
}
