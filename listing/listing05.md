Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
Вывод:
error

Интерфейс хранит в себе тип интерфейса и тип самого значения.
В Go интерфейсный тип выглядит вот так:
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

Где tab — это указатель на itable — структуру, которая хранит некоторые метаданные о типе
и список методов, используемых для удовлетворения интерфейса.
data — указывает на фактическую переменную с конкретным (статическим) типом.

Значение любого интерфейса, является nil в случае, когда и значение, и тип являются nil.
В нашем случае значение (фактическая переменная) - nil, a тип - кастомной ошибки, 
реализующей интерфейс ошибки, информация о типе хранится в tab,
поэтому интерфейс не равен nil. 

```
