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
