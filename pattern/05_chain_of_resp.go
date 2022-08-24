package pattern

import (
	"log"
)

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

//Уровни логирования
const (
	DebugLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
)

//Logger -  интерфейс логгеров-обработчиков
type Logger interface {
	log(msg string, level int)
	setNext(logger Logger)
}

//Info - логгер информационных сообщений
type Info struct {
	next Logger
}

//добавляем следующий
func (d *Info) setNext(l Logger) {
	d.next = l
}

//логируем или передаем дальше
func (d *Info) log(msg string, level int) {
	if level == InfoLevel {
		log.Println("INFO:", msg)
		return
	}
	d.next.log(msg, level)
}

//Debug - логгер дебажной информации
type Debug struct {
	next Logger
}

//логируем или передаем дальше
func (d *Debug) log(msg string, level int) {
	if level == DebugLevel {
		log.Println("DEBUG:", msg)
		return
	}
	d.next.log(msg, level)
}

//добавляем следующий
func (d *Debug) setNext(l Logger) {
	d.next = l
}

//Warning - логгер предупреждений
type Warning struct {
	next Logger
}

//логируем или передаем дальше
func (d *Warning) log(msg string, level int) {
	if level == WarningLevel {
		log.Println("WARNING:", msg)
		return
	}
	d.next.log(msg, level)
}

//добавляем следующий
func (d *Warning) setNext(l Logger) {
	d.next = l
}

//Error - логгер ошибок
type Error struct {
	next Logger
}

//логируем
func (d *Error) log(msg string, level int) {
	log.Println("ERROR:", msg)
	return
}

//добавляем следующий
func (d *Error) setNext(l Logger) {
	d.next = l
}

//TestChain - демонстрация паттерна
func TestChain() {
	var errorLog Error
	var warLog Warning
	warLog.setNext(&errorLog)
	var infoLog Info
	infoLog.setNext(&warLog)
	var debugLog Debug
	debugLog.setNext(&infoLog)
	var aLogger Logger = &debugLog
	aLogger.log("hello, world", InfoLevel)
	aLogger.log("hello, world", ErrorLevel)
}
