package main

import (
	"github.com/beevik/ntp"
	"io"
	"os"
	"time"
)

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

func getTime() (time.Time, error) {
	receivedTime, err := ntp.Time("0.ru.pool.ntp.org")
	return receivedTime, err
}

func main() {
	receivedTime, err := getTime()
	if err != nil {
		_, _ = io.WriteString(os.Stderr, "Unable to get time from NTP server: "+err.Error())
		os.Exit(1)
	}
	_, _ = io.WriteString(os.Stdout, receivedTime.Format(time.RFC850)+" UTC")
	os.Exit(0)
}
