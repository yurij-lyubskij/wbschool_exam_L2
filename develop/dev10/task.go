package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

var usageStr = `
Usage: go-telnet [] host port
Options:
--timeout (-t) - "timeout" - таймаут в секундах
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

var timeout string

//parseArgs - функция, в которой мы парсим аргументы
//считываем и приводим к нужному типу таймаут,
//считываем хост и порт
func parseArgs() (host string, port string, timeOut time.Duration, err error) {
	flag.Usage = usage
	flag.StringVar(&timeout, "timeout", "10s", "таймаут в секундах")
	flag.StringVar(&timeout, "t", "10s", "таймаут в секундах")
	flag.Parse()
	//проверяем, что таймаут - число, единицы измерения - секунды
	regExpr := regexp.MustCompile("^\\d+s")
	if !regExpr.MatchString(timeout) {
		return "", "", timeOut, errors.New("incorrect timeout value")
	}
	//приводим к целому типу
	timeNum, err := strconv.Atoi(timeout[:len(timeout)-1])
	if err != nil {
		return "", "", timeOut, errors.New("incorrect timeout value")
	}
	//получаем таймаут в секундах
	timeOut = time.Second * time.Duration(timeNum)
	args := flag.Args()
	//проверяем количество аргументов
	if len(args) < 2 {
		flag.Usage()
	}
	//получаем хост и порт
	host = args[0]
	port = args[1]
	err = nil
	return
}

//Reader - функция, считывающая информацию из
//сокета и отправляющая ее в стандартный вывод
//закрывается, если закрылся канал done
//или закрылся сокет
func Reader(conn net.Conn, done chan struct{}) {
	sc := bufio.NewScanner(conn)
	for sc.Scan() {
		select {
		case <-done:
			return
		default:
			str := sc.Text()
			fmt.Println(str)
		}
	}
	close(done)
}

//Writer - функция, считывающая информацию из
//стандартного ввода и отправляющая ее в сокет
//закрывается, если закрылся канал done
//или закончился ввод (ctrl + D)
func Writer(conn net.Conn, done chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-done:
			return
		default:
			str := scanner.Text()
			_, err := fmt.Fprintln(conn, str)
			if err != nil {
				close(done)
				return
			}
		}
	}
	return
}

func main() {
	//получаем хост, порт, таймаут
	host, port, timeOut, err := parseArgs()
	if err != nil {
		log.Fatalf("error while parsing arguments: %s", err)
	}
	//составляем адрес
	address := host + ":" + port
	//соединяемся с таймаутом
	conn, err := net.DialTimeout("tcp", address, timeOut)
	if err != nil {
		log.Fatalf("error happened: %s", err.Error())
	}
	//закрываем соединение
	defer conn.Close()

	log.Printf("client connected to %s:%s", host, port)
	//создаем канал для сигнала о конце работы
	done := make(chan struct{})
	// Запускаем чтение и запись
	go Reader(conn, done)
	Writer(conn, done)
	log.Printf("graceful shutdown")
}
