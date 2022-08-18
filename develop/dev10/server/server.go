package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const host = "localhost"
const port = "8080"

//Writer - функция, считывающая информацию из
//сокета и отправляющая ее в сокет
func Writer(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		str := scanner.Text()
		fmt.Fprintln(conn, "echo:"+str)
	}
	return
}

//Примитивный эхо-сервер
func main() {
	log.Println("server up")
	//адрес составляем из 2 констант
	address := host + ":" + port
	//слушаем на порту
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	//получаем сканер для получения команд со стндартного ввода
	sc := bufio.NewScanner(os.Stdin)
	//принимам соединение
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	//начинаем читать/писать из/в сокета
	go Writer(conn)
	//блокируемся до получения команды на завершение
	for sc.Scan() {
		if sc.Text() == "quit" {
			break
		}
	}
	log.Println("server down")
}
