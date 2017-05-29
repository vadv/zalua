package dsl

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslStringsSplit(L *lua.LState) int {
	str := L.CheckString(1)
	delim := L.CheckString(2)
	strSlice := strings.Split(str, delim)
	result := L.CreateTable(len(strSlice), 0)
	for _, str := range strSlice {
		result.Append(lua.LString(str))
	}
	L.Push(result)
	return 1
}
