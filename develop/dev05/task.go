package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var usageStr = `
Usage: grep [options] [pattern] [file]
file - stdin по умолчанию
Options:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки
-e - "pattern list", список паттернов
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
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

//CheckMatch - тип функуии проверки строки на соответствие
type CheckMatch func(s string) bool

//RegExprFuncGet - функция, составляющая регулярное выражение
//согласно переданным аргументам и возвращающая функцию
//проверки строки на соответствие и ошибку
func RegExprFuncGet(args []string) (CheckMatch, error) {
	var result CheckMatch
	regExprStr := strings.Builder{}
	regExprStr.WriteRune('(')
	if ignore {
		regExprStr.WriteString("?i)(")
	}
	for i, str := range args {
		if i != 0 {
			regExprStr.WriteRune('|')
		}
		regExprStr.WriteRune('(')
		regExprStr.WriteString(str)
		regExprStr.WriteRune(')')
	}
	regExprStr.WriteRune(')')
	regex := regExprStr.String()
	reg, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	result = func(s string) bool {
		res := reg.MatchString(s)
		if invert {
			res = !res
		}
		return res
	}
	return result, nil
}

//RegExprFuncGet - функция, возвращающая функцию
//проверки строки на соответствие и ошибку
//она либо составляет регулярное выражение
//с помощью RegExprFuncGet, либо проверяет вхождение
//подстроки в случае fixed
func MatchFuncGet(args []string) (CheckMatch, error) {
	var result CheckMatch
	if fixed {
		result = func(s string) bool {
			res := false
			for _, str := range args {
				res = res || (strings.Contains(s, str))
			}
			if invert {
				res = !res
			}
			return res
		}
		return result, nil
	}
	result, err := RegExprFuncGet(args)
	return result, err
}

//printCount - функция, читающая файл,
//и выводящая число подходящих строк
func printCount(sc *bufio.Scanner, ch CheckMatch, s io.Writer) error {
	var i int
	for sc.Scan() {
		if ch(sc.Text()) {
			i++
		}
	}
	fmt.Fprintln(s, i)
	// handle first encountered error while reading
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

//printFound - функция, читающая файл,
//и выводящая подходящие строки в
//соответствии с форматированием
//преимуществом явяется то,
//что функция читает файл построчно
//и хранит в памяти только буфер
//строк контекста
func printFound(sc *bufio.Scanner, ch CheckMatch, f io.Writer) error {
	//получаем контекст
	a := after
	b := before
	if a < context {
		a = context
	}
	if b < context {
		b = context
	}
	//буфер строк до и после
	var befBuff []string
	var afterBuff []string
	//номер строки
	i := 0
	//сколько еще строк считать в буфер
	//обновляется, если снова нашли соответствие
	var afterCount int
	//сканируем файл построчно
	for sc.Scan() {
		i++
		cur := sc.Text()
		out := cur
		//номер строки
		if line {
			out = fmt.Sprint(i, " ") + cur
		}
		//проверка на соответствие
		if ch(cur) {
			afterCount = a + 1
		}
		//по умолчанию пишем в буфер before
		if afterCount == 0 {
			befBuff = append(befBuff, out)
			//держим не более b строк
			if len(befBuff) > b {
				var newbuff []string
				newbuff = append(newbuff, befBuff[len(befBuff)-before:]...)
				befBuff = newbuff
			}
		}
		//если встретилось соответствие
		//записываем эту строку и строки After
		if afterCount > 0 {
			afterBuff = append(afterBuff, out)
			afterCount--
			//когда собрали нужные строки, выводим все
			if afterCount == 0 {
				for _, s := range befBuff {
					fmt.Fprintln(f, s)
				}
				for _, s := range afterBuff {
					fmt.Fprintln(f, s)
				}
				befBuff = []string{}
				afterBuff = []string{}
			}
		}
	}
	//если файл закончился, а есть невыведенные данные
	if len(afterBuff) > 0 {
		for _, s := range befBuff {
			fmt.Fprintln(f, s)
		}
		for _, s := range afterBuff {
			fmt.Fprintln(f, s)
		}
	}
	// обработка первой ошибки
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

//настройки утилиты
var (
	after    int
	before   int
	context  int
	count    bool
	ignore   bool
	invert   bool
	fixed    bool
	line     bool
	patterns bool
)

//grep - функция, открывающая файл
//читающая его и выводящая из него
//нужные строки или число строк
//в передаваемый writer
//также с обработкой ошибок
func grep(args []string, s io.Writer) {
	if len(args) < 1 {
		//мало аргументов
		usage()
		return
	}
	//как и настоящая утилита, ошибка, если эти флаги вместе
	if invert && (context > 0 || before > 0 || after > 0) {
		log.Fatalf("invalid context length argument")
	}

	//входной файл
	finName := args[len(args)-1]
	//получаем ридер на входной файл
	sc, err := ReaderGet(finName)
	if err != nil {
		log.Fatalf("Ошибка открытия файла: %s", err)
	}
	// остальные аргументы - регулярные выражения
	format := args[:len(args)-1]
	// без нужного флага нельзя больше 1 регулярного выражения
	if !patterns && (len(format) > 1) {
		log.Fatalf("передано больше одного выражения")
	}
	//получаем функцию проверки на соответствие
	ch, err := MatchFuncGet(format)
	if err != nil {
		log.Fatalf("Ошибка функции сравнения: %s", err)
	}
	//если нужно только число строк
	if count {
		err = printCount(sc, ch, s)
		if err != nil {
			log.Fatalf("Ошибка чтения: %s", err)
		}
		return
	}
	//выводим нужные строки
	err = printFound(sc, ch, s)
	if err != nil {
		log.Fatalf("Ошибка чтения: %s", err)
	}
	return
}

func main() {
	//настройки утилиты - возможные флаги
	flag.IntVar(&after, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&before, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&context, "C", 0, "(A+B) печатать ±N строк вокруг совпадения")
	flag.BoolVar(&count, "с", false, "количество строк")
	flag.BoolVar(&ignore, "i", false, "игнорировать регистр")
	flag.BoolVar(&invert, "v", false, "вместо совпадения, исключать")
	flag.BoolVar(&fixed, "F", false, "точное совпадение со строкой, не паттерн")
	flag.BoolVar(&line, "n", false, "печатать номер строки")
	flag.BoolVar(&patterns, "e", false, "список паттернов")
	flag.Usage = usage
	//парсим флаги
	flag.Parse()
	args := flag.Args()
	grep(args, os.Stdout)
}
