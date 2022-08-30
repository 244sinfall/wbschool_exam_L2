package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"shell/pkg/controls"
	"strings"
)

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/
func main() {
	homeDir, err := os.UserHomeDir()
	err = os.Chdir(homeDir)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error starting the shell: " + err.Error() + "\n")
		os.Exit(1)
	}
	for {
		wd, _ := os.Getwd()
		fmt.Print("DimaShell | " + wd + ">")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		commands := controls.SeparateCommands(input)
		if len(commands) == 1 {
			command := strings.Split(commands[0], " ")
			controls.ExecCommand(command)
		}
		if len(commands) > 1 {
			var buf bytes.Buffer
			for i, cmd := range commands {
				old := os.Stdout
				r, w, _ := os.Pipe()
				if i != len(commands)-1 {
					os.Stdout = w
				}
				command := strings.Split(cmd, " ")
				if buf.Len() != 0 {
					out := strings.TrimSpace(buf.String())
					command = append(command, out)
				}
				controls.ExecCommand(command)
				buf.Reset()
				if i != len(commands)-1 {
					_ = w.Close()
					os.Stdout = old
					_, _ = io.Copy(&buf, r)
				}
			}
		}
	}
}
