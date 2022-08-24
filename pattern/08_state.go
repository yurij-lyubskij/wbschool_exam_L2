package pattern

import "fmt"

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

//AbstractContext - интерфейс контекста
type AbstractContext interface {
	operation()
	SetState(s State)
}

//StateContext  - содержит состояние
type StateContext struct {
	state State
}

//вызываем обработку и смену состояния
func (c *StateContext) operation() {
	c.state.handle()
}

//SetState - меняем состояние
func (c *StateContext) SetState(s State) {
	c.state = s
}

//State - интерфейс состояния
type State interface {
	handle()
	switchState(c AbstractContext)
}

//StateA - первое состояние
//содержит абстрактный контекст
//(указатель на контекст)
//для переключения состояний
type StateA struct {
	con AbstractContext
}

//переключаем состояние
func (s *StateA) switchState(con AbstractContext) {
	newstate := StateB{}
	newstate.con = con
	con.SetState(&newstate)
}

//обработка и переключение состояния
func (s *StateA) handle() {
	fmt.Println("from A to B")
	s.switchState(s.con)
}

//StateB - первое состояние
//содержит абстрактный контекст
//(указатель на контекст)
//для переключения состояний
type StateB struct {
	con AbstractContext
}

//обработка и переключение состояния
func (s *StateB) handle() {
	fmt.Println("from B to A")
	s.switchState(s.con)
}

//переключаем состояние
func (s *StateB) switchState(con AbstractContext) {
	newstate := StateA{}
	newstate.con = con
	con.SetState(&newstate)
}

//NewStateContext создает контекст с начальным состоянием
func NewStateContext() AbstractContext {
	var c StateContext
	var state = StateA{}
	state.con = &c
	c.state = &state
	return &c
}

//TestState - демонстрация паттерна
func TestState() {
	var context = NewStateContext()
	context.operation()
	context.operation()
	context.operation()
	context.operation()
}
