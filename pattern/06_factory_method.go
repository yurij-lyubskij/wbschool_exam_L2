package pattern

import "fmt"

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

//Product  - интерфейс продукта
type Product interface {
	getName() string
}

//ConcreteProductA  - конкретный продукт
type ConcreteProductA struct {
}

//реализация интерфейса
func (ConcreteProductA) getName() string {
	return "ConcreteProductA"
}

//ConcreteProductB  - конкретный продукт
type ConcreteProductB struct {
}

//реализация интерфейса
func (ConcreteProductB) getName() string {
	return "ConcreteProductB"
}

//Creator - интерфейс для создания продукта
type Creator interface {
	factoryMethod() Product
}

//ConcreteCreatorA - создание продукта А
type ConcreteCreatorA struct {
}

//создание продукта А
func (ConcreteCreatorA) factoryMethod() Product {
	return ConcreteProductA{}
}

//ConcreteCreatorB - создание продукта B
type ConcreteCreatorB struct {
}

//создание продукта B
func (ConcreteCreatorB) factoryMethod() Product {
	return ConcreteProductB{}
}

//TestFactoryMeth - демонстрация паттерна
func TestFactoryMeth() {
	//создаем объекты, способные порождать продукт
	var CreatorA ConcreteCreatorA
	var CreatorB ConcreteCreatorB
	//сохраняем в массив интерфейсов
	creators := []Creator{&CreatorA, &CreatorB}
	//создаем продукт с помощью интерфейса
	//и сохраняем, как интерфейс
	for _, creator := range creators {
		var product = creator.factoryMethod()
		//смотрим, что за продукт
		fmt.Println(product.getName())
	}
}
