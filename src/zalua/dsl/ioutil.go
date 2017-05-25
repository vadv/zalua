package dsl

import (
	"io/ioutil"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslIoutilReadFile(L *lua.LState) int {
	filename := L.CheckString(1)
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		L.Push(lua.LString(string(data)))
	} else {
		L.Push(lua.LNil)
	}
	return 1
}
