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

func (d *dslConfig) dslStringsHasPrefix(L *lua.LState) int {
	str1 := L.CheckString(1)
	str2 := L.CheckString(2)
	result := strings.HasPrefix(str1, str2)
	L.Push(lua.LBool(result))
	return 1
}

func (d *dslConfig) dslStringsHasSuffix(L *lua.LState) int {
	str1 := L.CheckString(1)
	str2 := L.CheckString(2)
	result := strings.HasSuffix(str1, str2)
	L.Push(lua.LBool(result))
	return 1
}

func (d *dslConfig) dslStringsTrim(L *lua.LState) int {
	str1 := L.CheckString(1)
	str2 := L.CheckString(2)
	result := strings.Trim(str1, str2)
	L.Push(lua.LString(result))
	return 1
}
