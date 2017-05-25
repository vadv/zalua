package dsl

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslOsStat(L *lua.LState) int {
	path := L.CheckString(1)
	stat, err := os.Stat(path)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}
	result := L.NewTable()
	L.SetField(result, "size", lua.LNumber(stat.Size()))
	L.SetField(result, "is_dir", lua.LBool(stat.IsDir()))
	L.SetField(result, "mod_time", lua.LNumber(stat.ModTime().Unix()))
	L.Push(result)
	return 1
}

func (d *dslConfig) dslOsPagesize(L *lua.LState) int {
	L.Push(lua.LNumber(os.Getpagesize()))
	return 1
}
