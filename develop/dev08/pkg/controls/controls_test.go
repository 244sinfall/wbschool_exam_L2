package controls

import (
	"bytes"
	"github.com/shirou/gopsutil/process"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestExecProcess(t *testing.T) {
	timeStart := time.Now().UnixMilli()
	ExecProcess([]string{"notepad"})
	pcs, _ := process.Processes()
	for _, p := range pcs {
		n, _ := p.Name()
		if strings.Contains(n, "Notepad") {
			ct, _ := p.CreateTime()
			if ct-timeStart > 5000 {
				t.Fail()
				t.Log(ct, timeStart)
			} else {
				t.Log("OK")
				return
			}
		}
	}
	t.Fail()
}

func TestKillProcess(t *testing.T) {
	ExecProcess([]string{"notepad"})
	time.Sleep(1 * time.Second)
	var pid int32
	pcs, _ := process.Processes()
	for _, p := range pcs {
		n, _ := p.Name()
		if strings.Contains(n, "Notepad") {
			pid = p.Pid
		}
	}
	if pid != 0 {
		KillProcess(int(pid))
	} else {
		t.Fail()
	}
	exist, _ := process.PidExists(pid)
	if exist {
		t.Log(process.PidExists(pid))
		t.Fail()
	}
}

func TestKillProcessByName(t *testing.T) {
	ExecProcess([]string{"notepad"})
	time.Sleep(1 * time.Second)
	var pid int32
	pcs, _ := process.Processes()
	for _, p := range pcs {
		n, _ := p.Name()
		if strings.Contains(n, "Notepad") {
			pid = p.Pid
			break
		}
	}
	if pid != 0 {
		KillProcessByName("Notepad.exe")
		time.Sleep(1 * time.Second)
		exist, _ := process.PidExists(pid)
		if exist {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestReadDirectory(t *testing.T) {
	wd, _ := os.Getwd()
	read, _ := os.ReadDir(wd)
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	//
	ReadDirectory(wd)
	//
	_ = w.Close()
	os.Stdout = old
	_, _ = io.Copy(&buf, r)
	if len(strings.Split(buf.String(), "\n"))-1 != len(read) { // -1 => last line has new line symbol too
		t.Log(len(read), len(strings.Split(buf.String(), "\n")))
		t.Fail()
	}
}

func TestSeparateCommands(t *testing.T) {
	cmd1 := "cmd1 | cmd2 | cmd3 param arg param | cmd4 | cmd55"
	cmds1 := SeparateCommands(cmd1)
	cmd2 := `echo "hello world" |`
	cmds2 := SeparateCommands(cmd2)
	cmd3 := `echo | pwd | cd 123`
	cmds3 := SeparateCommands(cmd3)
	if len(cmds1) != 5 || cmds1[2] != "cmd3 param arg param" {
		t.Fail()
		t.Log(cmds1)
	}
	if len(cmds2) != 1 || cmds2[0] != `echo "hello world"` {
		t.Fail()
		t.Log(cmds2)
	}
	if len(cmds3) != 3 || cmds3[2] != "cd 123" {
		t.Fail()
		t.Log(cmds3)
	}
}

func TestSetDirectory(t *testing.T) {
	_ = os.Mkdir("test_dir", 0666)
	SetDirectory("test_dir")
	wd, _ := os.Getwd()
	if !strings.Contains(wd, "test_dir") {
		t.Fail()
		t.Log(wd)
	}
	SetDirectory("C:\\")
	wd, _ = os.Getwd()
	if wd != "C:\\" {
		t.Fail()
		t.Log(wd)
	}
}
