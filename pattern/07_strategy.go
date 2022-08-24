package pattern

import (
	"fmt"
)

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

//Strategy - интерфейс стратегии
type Strategy interface {
	execute()
}

//StrategyA - конкретная стратегия
type StrategyA struct {
}

//реализация интерфейса
//здесь реализован сам алгоритм
func (s *StrategyA) execute() {
	fmt.Println("executing strategy A")
}

//StrategyB - конкретная стратегия
type StrategyB struct {
}

//реализация интерфейса
//здесь реализован другой алгоритм
func (s *StrategyB) execute() {
	fmt.Println("executing strategy B")
}

//Context  - содержит
//интерфейс для работы
//со стратегией
type Context struct {
	strategy Strategy
}

//выполняет стратегию
func (c *Context) operate() {
	if c.strategy == nil {
		fmt.Println("no strategy for now")
		return
	}
	c.strategy.execute()
}

//SetStrategy - меняем стратегию в runtime
func (c *Context) SetStrategy(s Strategy) {
	c.strategy = s
}

//TestStrategy - демонстрация паттерна
func TestStrategy() {
	var context Context
	context.operate()
	var strA StrategyA
	var strB StrategyB
	context.SetStrategy(&strA)
	context.operate()
	context.SetStrategy(&strB)
	context.operate()
	context.SetStrategy(&strA)
	context.operate()
}
