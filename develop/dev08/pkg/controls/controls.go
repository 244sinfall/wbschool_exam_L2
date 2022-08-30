package controls

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func SeparateCommands(input string) []string {
	commands := strings.Split(input, "|")
	for i, cmd := range commands {
		commands[i] = strings.TrimSpace(cmd)
	}
	return commands
}

func ExecCommand(command []string) {
	wd, _ := os.Getwd()
	switch command[0] {
	case "cd":
		newDir := wd
		if len(command) == 1 {
			newDir, _ = os.UserHomeDir()
		}
		if len(command) == 2 {
			newDir = command[1] // Пробуем для начала смотреть из текущей директории
		}
		SetDirectory(newDir)
	case "ls":
		ReadDirectory(wd)
	case "pwd":
		fmt.Println(wd)
	case "ps":
		ShowProcesses()
	case "kill":
		if len(command) == 2 {
			pid, err := strconv.Atoi(command[1])
			if err != nil {
				KillProcessByName(command[1])
				break
			} else {
				KillProcess(pid)
			}

		}
	case "exec":
		if len(command) > 1 {
			ExecProcess(command[1:])
		}
	case "echo":
		if len(command) > 1 {
			fmt.Println(strings.Join(command[1:], " "))
		}
	default:
		ExecProcess(command)
	case "\\quit":
		os.Exit(0)
	}
}

func ExecProcess(command []string) {
	wd, _ := os.Getwd()
	path := wd + "\\" + command[0]
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = command[0]
	}
	cmd := exec.Command(path)
	if len(command) > 1 {
		cmd.Args = command[1:]
	}
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func KillProcessByName(n string) {
	processes, _ := process.Processes()
	var removed bool
	for _, p := range processes {
		name, _ := p.Name()
		if n == name {
			err := p.Kill()
			removed = true
			if err != nil {
				fmt.Printf("Error: %v | Process: %v\n", err, p.Pid)
				continue
			}

		}
	}
	if !removed {
		fmt.Println("No processes with such name: ", n)
	}
}

func KillProcess(pid int) {
	processes, _ := process.Processes()
	for _, p := range processes {
		if int(p.Pid) == pid {
			err := p.Kill()
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	fmt.Println("No process with such pid.")
}

func ShowProcesses() {
	v, _ := mem.VirtualMemory()
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	processes, _ := process.Processes()
	for _, p := range processes {
		n, _ := p.Name()
		percent, _ := p.MemoryPercent()
		fmt.Printf("Process: %v\tPID: %v\tUsing memory: %v\n", n, p.Pid, percent)
	}
}

func ReadDirectory(wd string) {
	entries, err := os.ReadDir(wd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, e := range entries {
		fmt.Println(e.Type(), e.Name())
	}
}

func SetDirectory(dir string) {
	wd, _ := os.Getwd()
	var err error
	newDir := wd + "\\" + dir // Пробуем для начала смотреть из текущей директории
	err = os.Chdir(newDir)
	if err != nil {
		err = os.Chdir(dir)
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
