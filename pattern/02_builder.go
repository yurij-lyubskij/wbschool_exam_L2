package pattern

import "fmt"

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

//Pizza - наш продукт
type Pizza struct {
	Sauce   string
	Topping string
}

//SetSauce - часть создания пиццы
func (p *Pizza) SetSauce(s string) {
	p.Sauce = s
}

//SetTopping - часть создания пиццы
func (p *Pizza) SetTopping(t string) {
	p.Topping = t
}

//ShowPizza  - показываем, какой объект создался
func (p *Pizza) ShowPizza() {
	fmt.Println("sauce = ", p.Sauce)
	fmt.Println("topping = ", p.Topping)
}

//AbstractBuilder - интерфейс для создания пиццы
type AbstractBuilder interface {
	GetPizza() Pizza
	BuildSauce()
	BuildTopping()
}

//HawaiianPizzaBuilder - конкретный builder
type HawaiianPizzaBuilder struct {
	pizza Pizza
}

//BuildSauce - реализация метода интерфейса
func (p *HawaiianPizzaBuilder) BuildSauce() {
	p.pizza.SetSauce("mild")
}

//BuildTopping - реализация метода интерфейса
func (p *HawaiianPizzaBuilder) BuildTopping() {
	p.pizza.SetTopping("ham and pineapple")
}

//GetPizza - получаем готовый продукт
func (p *HawaiianPizzaBuilder) GetPizza() Pizza {
	return p.pizza
}

//Waiter - класс-директор. Управляет созданием
//пиццы с помощью builder
type Waiter struct {
	pizzaBuilder AbstractBuilder
}

//SetBuilder - задает конкретный builder
func (p *Waiter) SetBuilder(builder AbstractBuilder) {
	p.pizzaBuilder = builder
}

//MakePizza - делает пиццу с помощью
//builder по шагам
func (p *Waiter) MakePizza() {
	p.pizzaBuilder.BuildTopping()
	p.pizzaBuilder.BuildSauce()
}

//GetPizza - получаем готовую пиццу
func (p *Waiter) GetPizza() Pizza {
	return p.pizzaBuilder.GetPizza()
}

//TestBuilder - показывает builder
func TestBuilder() {
	var waiter Waiter
	var builder HawaiianPizzaBuilder
	waiter.SetBuilder(&builder)
	waiter.MakePizza()
	pizza := waiter.GetPizza()
	pizza.ShowPizza()
}
