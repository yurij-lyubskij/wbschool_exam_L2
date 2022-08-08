package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

//Unpack - Функция распаковки
func Unpack(s string) (string, error) {
	runes := []rune(s)
	length := len(runes)

	// пустую строку просто вернем
	if length == 0 {
		return s, nil
	}

	var result strings.Builder
	var number strings.Builder
	var prev rune

	// вернем ошибку, если строка начинается с числа
	if unicode.IsDigit(runes[0]) {
		return "", errors.New("некорректная строка")
	}

	for i := 0; i < length; i++ {
		switch {
		case unicode.IsDigit(runes[i]):

			//запишем число повторений в строку
			number.WriteRune(runes[i])

			// пока следующая руна - цифра, пишем в number
			for i < length-1 && unicode.IsDigit(runes[i+1]) {
				i++
				number.WriteRune(runes[i])
			}

			//переводим в число
			n, err := strconv.Atoi(number.String())
			if err != nil {
				return "", err
			}

			if n == 0 {
				return "", errors.New("0 повторений в строке - ошибка")
			}
			number.Reset()

			// если число повторений = 1, просто идем дальше
			if n == 1 {
				break
			}

			// добавляем нужное число повторений
			for j := 0; j < n-1; j++ {
				result.WriteRune(prev)
			}

		case runes[i] == '\\':

			// в конце строки '\' некорректен
			if i == length-1 {
				return "", errors.New("в конце строки '\\' некорректен")
			}

			// обрабатывается руну после escape как символ
			i++
			result.WriteRune(runes[i])
			prev = runes[i]
		default:
			result.WriteRune(runes[i])
			prev = runes[i]
		}
	}

	return result.String(), nil
}

func main() {
	str := `qwe\4\5`
	unpacked, err := Unpack(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(unpacked)
}
