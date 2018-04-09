package dsl

import (
	"crypto/md5"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslCryptoMD5(L *lua.LState) int {
	h := md5.New()
	h.Write([]byte(L.CheckString(1)))
	L.Push(lua.LString(fmt.Sprintf("%x", h.Sum(nil))))
	return 1
}
