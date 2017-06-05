package dsl

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslTimeUnix(L *lua.LState) int {
	L.Push(lua.LNumber(time.Now().Unix()))
	return 1
}

func (d *dslConfig) dslTimeUnixNano(L *lua.LState) int {
	L.Push(lua.LNumber(time.Now().UnixNano()))
	return 1
}

func (d *dslConfig) dslTimeSleep(L *lua.LState) int {
	val := L.CheckInt64(1)
	time.Sleep(time.Duration(val) * time.Second)
	return 0
}
