package pattern

import (
	"fmt"
	"math"
)

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

//Shape - интерфейс
//геометрической фигуры
type Shape interface {
	accept(visitor Visitor)
}

//Circle - круг
type Circle struct {
	r float64
}

//Реализация интерфейса для круга
func (c *Circle) accept(visitor Visitor) {
	visitor.visitForCircle(c)
}

// Rectangle - прямоугольник
type Rectangle struct {
	a float64
	b float64
}

//Реализация интерфейса для прямоугольника
func (r *Rectangle) accept(visitor Visitor) {
	visitor.visitForRectangle(r)
}

//Visitor - интерфейс, который позволяет добавлять новый
//функционал к гемоетрическим фигурам
type Visitor interface {
	visitForCircle(c *Circle)
	visitForRectangle(r *Rectangle)
}

//СalcPerim - реализует интерфейс.
//методы считают периметр разных фигур
type СalcPerim struct {
}

//периметр круга
func (calc *СalcPerim) visitForCircle(c *Circle) {
	fmt.Println(2 * math.Pi * c.r)
}

//периметр прямоугольника
func (calc *СalcPerim) visitForRectangle(r *Rectangle) {
	fmt.Println(2 * (r.a + r.b))
}

//TestVisitor - демонстрация шаблона
func TestVisitor() {
	circle := Circle{1}
	rectangle := Rectangle{1, 2}
	perim := СalcPerim{}
	circle.accept(&perim)
	rectangle.accept(&perim)
}
