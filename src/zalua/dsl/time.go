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

func (d *dslConfig) dslTimeParse(L *lua.LState) int {
	layout, value := L.CheckString(1), L.CheckString(2)
	result, err := time.Parse(layout, value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LNumber(result.UTC().Unix()))
	return 1
}
