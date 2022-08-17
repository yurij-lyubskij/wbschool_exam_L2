package main

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов, реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

//очевидное решение на select
//с рефлексией, чтобы определять количество кейсов select
//пишем из нескольких каналов в 1
//когда хотя бы 1 канал из исходных закроется,
//закроется и канал or
func or(channels ...<-chan interface{}) <-chan interface{} {
	orchan := make(chan interface{})
	//определяем количество кейсов select
	allCases := make([]reflect.SelectCase, len(channels))
	//собираем кейсы select
	for i, ch := range channels {
		allCases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	go func() {
		//при выходе закрываем канал
		defer close(orchan)
		ok := true
		var value interface{}
		//слушаем в бесконечном цикле,
		//пока один из каналов не закроется
		for ok {
			_, value, ok = reflect.Select(allCases)
			if ok {
				//пишем сообщение в общий канал
				orchan <- value
			}
		}
	}()
	return orchan
}

//решение через горутины
//и завершение горутин через контекст
func orParallel(channels ...<-chan interface{}) <-chan interface{} {
	orchan := make(chan interface{})
	//получаем контексn и функцию cancel
	ctx, cancel := context.WithCancel(context.Background())
	//используем WaitGroup, чтобы ждать завершения горутин
	wg := &sync.WaitGroup{}
	//для каждого канала запускаем горутину
	for _, ch := range channels {
		//увеличиваем счетчик
		wg.Add(1)
		go func(channel <-chan interface{}) {
			//после выполнения уменьшаем счетчик
			defer wg.Done()
			var value interface{}
			var ok bool
			//в бесконечном цикле слушаем канал
			//и ждем сообщения о завершении
			for {
				select {
				case value, ok = <-channel:
					//если получили сообщение
					//передаем его в канал or
					if ok {
						orchan <- value
						//если канал закрыли
						//отменяем контекст
					} else {
						cancel()
					}
				//если контекст отменен
				//завершаем выполнение
				case _ = <-ctx.Done():
					return
				}
			}
		}(ch)
	}
	go func() {
		//ждем завершения всех горутин после отмены контекста
		wg.Wait()
		//закрываем канал
		close(orchan)
	}()

	return orchan
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
			//c <- 1
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Println()
	fmt.Printf("done after %v", time.Since(start))
	runtime.GOMAXPROCS(0)
	start = time.Now()
	<-orParallel(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Microsecond),
		sig(1*time.Microsecond),
		sig(2*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Println()
	fmt.Printf("done after %v", time.Since(start))
}
