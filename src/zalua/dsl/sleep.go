package dsl

import (
	"time"

	"github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslSleep(L *lua.LState) int {
	val := L.CheckInt64(1)
	time.Sleep(time.Duration(val) * time.Second)
	return 0
}
