package dsl

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"

	"zalua/logger"
)

func dslLog(level string, L *lua.LState) {
	logTime := time.Now().Format("02/01/2006 15:04:05")
	fd, err := logger.GetLogFD()
	if err != nil {
		L.RaiseError("internal error: %s", err.Error())
	}
	fmt.Fprintf(fd, "%s [%s] %s\n", logTime, level, L.CheckString(1))
}

func (d *dslConfig) dslLogError(L *lua.LState) int {
	dslLog("ERROR", L)
	return 0
}

func (d *dslConfig) dslLogInfo(L *lua.LState) int {
	dslLog("INFO", L)
	return 0
}
