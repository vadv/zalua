package dsl

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestExecCmd(t *testing.T) {

	testStr := `
local result, err = cmd.exec("unknown command")
if (err == nil) then error("run unknown "..tostring(result.code)) end

local result, err = cmd.exec("echo ok")
if not(err == nil) then error(err) end

local out = "ok\n"; if goruntime.goos == "windows" then out = "ok\r\n" end
if not(result.stdout == out) then error(result.stdout) end
`

	state := lua.NewState()
	Register(NewConfig(), state)
	if err := state.DoString(testStr); err != nil {
		t.Fatalf("execute lua error: %s\n", err.Error())
	}

}
