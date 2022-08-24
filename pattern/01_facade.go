package pattern

import "fmt"

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/

//CPU компьютера
type CPU struct {
}

//Freeze CPU
func (c *CPU) Freeze() {
	fmt.Println("Freeze")
}

//Jump to position
func (c *CPU) Jump(position int64) {
	fmt.Println("Jump")
}

//Execute command
func (c *CPU) Execute() {
	fmt.Println("Execute")
}

//HardDrive - жесткий диск компьютера
type HardDrive struct {
}

//Read - читать с диска
func (c *HardDrive) Read(lba int64, size int) (data *uint8) {
	fmt.Println("Read")
	return
}

//Memory - оперативная память
type Memory struct {
}

//Load - загрузить в память
func (c *Memory) Load(position int64, data *uint8) {
	fmt.Println("Load")
	return
}

//адрес, сектор, размер сектора загрузки
//т.к. пример абстрактный, зануляем
const bootAddress = 0
const bootSector = 0
const sectorSize = 0

//Facade - простой интерфейс
type Facade interface {
	Start()
}

//ComputerFacade - фасад компьютера
type ComputerFacade struct {
	cpu       CPU
	memory    Memory
	hardDrive HardDrive
}

//Start - запуск компьютера. Под капотом
//выполняет набор действий
func (c ComputerFacade) Start() {
	c.cpu.Freeze()
	c.memory.Load(bootAddress, c.hardDrive.Read(bootSector, sectorSize))
	c.cpu.Jump(bootAddress)
	c.cpu.Execute()
	return
}

//TestFacade - демонстрация шаблона
func TestFacade() {
	var computer Facade = &ComputerFacade{}
	computer.Start()
}
