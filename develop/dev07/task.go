package main

import (
	"fmt"
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

func createOrChannel(channels ...<-chan interface{}) <-chan interface{} {

	output := make(chan interface{}) // Создаем канал, в который будем посылать сообщение ,как только один из каналов закроется
	for _, c := range channels {     // Итерируемся по всем каналам, которые дали в функцию
		go func(c <-chan interface{}) { // Запускаем на каждый канал горутину, которая будет принимать сигналы через value, ok
			_, ok := <-c
			if !ok { // Если ok - false, это значит, что канал c закрыт. Т.е один из каналов, которые мы объединяем
				output <- struct{}{} // Сразу же посылаем сообщение в output канал, который мы вернули уже после запуска горутины
			}
		}(c)
	}
	return output
}

func main() {
	start := time.Now()
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	<-createOrChannel(sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute))
	fmt.Printf("Done after %v", time.Since(start))
}
