package dsl

import (
	"fmt"
	"syscall"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslStatFs(L *lua.LState) int {
	path := L.CheckString(1)
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("syscall statfs error: %s\n", err.Error())))
		return 2
	}
	result := L.NewTable()
	result.RawSetString("size", lua.LNumber(float64(fs.Blocks)*float64(fs.Bsize)))
	result.RawSetString("free", lua.LNumber(float64(fs.Bfree)*float64(fs.Bsize)))
	result.RawSetString("avail", lua.LNumber(float64(fs.Bavail)*float64(fs.Bsize)))
	result.RawSetString("files", lua.LNumber(float64(fs.Files)))
	result.RawSetString("files_free", lua.LNumber(float64(fs.Ffree)))
	L.Push(result)
	return 1
}
