package dsl

import (
	"bytes"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslCmdExec(L *lua.LState) int {
	command := L.CheckString(1)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", command)
	case "windows":
		cmd = exec.Command("cmd", "/c", command)
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(`unsupported os`))
		return 2
	}

	stdout, stderr := bytes.Buffer{}, bytes.Buffer{}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(10) * time.Second):
		L.Push(lua.LNil)
		L.Push(lua.LString(`execute timeout after 10 seconds`))
		return 2
	case err := <-done:
		result := L.CreateTable(0, 0)
		L.SetField(result, "stdout", lua.LString(stdout.String()))
		L.SetField(result, "stderr", lua.LString(stderr.String()))
		L.SetField(result, "code", lua.LNumber(-1))

		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					L.SetField(result, "code", lua.LNumber(int64(status.ExitStatus())))
				}
			}
			L.Push(result)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.SetField(result, "code", lua.LNumber(0))
		L.Push(result)
		return 1
	}
}
