package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var usageStr = `
Usage: cut [options] [file]
file - stdin по умолчанию
Options:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

//настройки утилиты
var (
	fields    string
	delimiter string
	separated bool
)

//начало и конец диапазона
type myrange struct {
	start  int
	finish int
}

//проверка вхождения в диапазон
func (m *myrange) in(i int) bool {
	return (m.start <= i) && (i <= m.finish)
}

//AllRanges - масссив диапазонов
type AllRanges struct {
	ranges []myrange
}

//проверка вхождения во все диапазоны
func (m *AllRanges) in(i int) bool {
	for _, a := range m.ranges {
		if a.in(i) {
			return true
		}
	}
	return false
}

//NewRanges - функция, где переводим диапазоны из строк в наш тип
func NewRanges(fields []string) (AllRanges, error) {
	var regExp [4]*regexp.Regexp
	var regStr = [4]string{`^\d+$`, `^\d+-\d+$`, `^\d+-$`, `^\-\d+$`}
	res := AllRanges{}
	var err error
	for i := 0; i < 4; i++ {
		regExp[i], err = regexp.Compile(regStr[i])
		if err != nil {
			return res, err
		}
	}
	for _, field := range fields {
		var matched bool
		for i := 0; i < 4; i++ {
			if regExp[i].MatchString(field) {
				matched = true
				oneRange := strings.Split(field, "-")
				var r myrange
				if i < 3 {
					r.start, err = strconv.Atoi(oneRange[0])
					r.start--
					if err != nil {
						return res, err
					}
				}
				if i%2 == 1 {
					r.finish, err = strconv.Atoi(oneRange[len(oneRange)-1])
					r.finish--
					if err != nil {
						return res, err
					}
				}
				if i == 2 {
					r.finish = math.MaxInt
				}
				if i == 0 {
					r.finish = r.start
				}
				if r.start <= r.finish {
					res.ranges = append(res.ranges, r)
				} else {
					fmt.Println(r.start, r.finish)
					return res, errors.New("invalid decreasing range")
				}
			}
		}
		if !matched {
			return res, errors.New("invalid range")
		}
	}
	return res, nil
}

//SplitStrings - функция, которая построчно читает файл,
//разбивает строку на колонки и выводит нужные
func SplitStrings(sc *bufio.Scanner, r AllRanges, f io.Writer) error {
	for sc.Scan() {
		str := sc.Text()
		//если нет сепаратора, выводим или нет в зависимости
		//от флага
		if !(strings.Contains(str, delimiter)) {
			if !separated {
				fmt.Fprintln(f, str)
			}
			continue
		}
		//разделяем строку на колонки
		splitStr := strings.Split(str, delimiter)
		var joinStr []string
		//ищем подходящие колонки
		for i, s := range splitStr {
			if r.in(i) {
				joinStr = append(joinStr, s)
			}
		}
		//склеиваем подходящие колонки
		res := strings.Join(joinStr, delimiter)
		//выводим подходящие колонки
		fmt.Fprintln(f, res)
	}
	// обработка первой ошибки
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

//ReaderGet - читает файл по передаваемому имени,
//возвращает сканнер и ошибку
//при пустом имени читает из stdin
func ReaderGet(fileName string) (*bufio.Scanner, error) {
	var sc *bufio.Scanner
	if len(fileName) == 0 {
		sc = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		sc = bufio.NewScanner(f)
	}
	return sc, nil
}

//cut - проверяем число аргументов,
//открываем файл, ищем диапазоны,
//выводим нужные колонки
func cut(args []string, writer io.Writer) {
	if len(args) > 1 {
		usage()
		return
	}
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	}
	sc, err := ReaderGet(fName)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %s", err)
	}
	if fields == "" {
		log.Fatalf("you must specify a list of bytes, characters, or fields")
	}

	fieldsList := strings.Split(fields, ",")
	ranges, err := NewRanges(fieldsList)
	if err != nil {
		log.Fatalf("Неправильно указано поле: %s", err)
	}
	err = SplitStrings(sc, ranges, writer)
	return
}

func main() {
	//настройки утилиты - возможные флаги
	flag.StringVar(&fields, "f", "", "выбрать поля (колонки)")
	flag.StringVar(&delimiter, "d", "\t", "использовать другой разделитель")
	flag.BoolVar(&separated, "s", false, "только строки с разделителем")
	flag.Usage = usage
	//парсим флаги
	flag.Parse()
	args := flag.Args()
	//основная работа утилиты
	cut(args, os.Stdout)
}
