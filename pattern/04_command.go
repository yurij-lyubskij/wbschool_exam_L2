package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

//Command - интерфейс комманды
type Command interface {
	execute()
}

//LightSwitch - receiver
//выключатель с возможностью
//смены состояния
type LightSwitch struct {
	On bool
}

//TurnOn - включить свет
func (l *LightSwitch) TurnOn() {
	l.On = true
}

//TurnOff - выключить свет
func (l *LightSwitch) TurnOff() {
	l.On = false
}

//LightOn - команда на
//включение света
type LightOn struct {
	l *LightSwitch
}

//реализация интерфейса команды
func (l *LightOn) execute() {
	l.l.TurnOn()
}

//LightOff - команда на
//выключение света
type LightOff struct {
	l *LightSwitch
}

//реализация интерфейса команды
func (l *LightOff) execute() {
	l.l.TurnOff()
}

//Invoker - класс,
//вызывающий команды
//по интерфейсу
type Invoker struct {
	commands []Command
}

//Register  - регистрируем команду
//возвращаем ее номер
func (i *Invoker) Register(c Command) (num int) {
	i.commands = append(i.commands, c)
	return len(i.commands) - 1
}

//Execute  - выполняем команду
func (i *Invoker) Execute(num int) {
	i.commands[num].execute()
}

//Client - управляет всеми объектами
//the client decides which receiver objects
//it assigns to the command objects, and which
//commands it assigns to the invoker. The client
//decides which commands to execute at which points (с)
type Client struct {
	switcher LightSwitch
	invoker  Invoker
}

//Run - выполянем команды в нужном парядке
func (client *Client) Run() {
	on := LightOn{&client.switcher}
	off := LightOff{&client.switcher}
	NumOn := client.invoker.Register(&on)
	NumOff := client.invoker.Register(&off)
	client.invoker.Execute(NumOn)
	fmt.Println(client.switcher)
	client.invoker.Execute(NumOff)
	fmt.Println(client.switcher)
}

//TestCommand - демонстрация паттерна
func TestCommand() {
	var client Client
	client.Run()
}
