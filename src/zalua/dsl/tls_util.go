package dsl

import (
	"crypto/tls"
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslTLSUtilCertGetNotAfter(L *lua.LState) int {
	serverName, address := L.CheckString(1), ""
	if L.GetTop() > 1 {
		address = L.CheckString(2)
	} else {
		address = fmt.Sprintf("%s:443", serverName)
	}
	conn, err := tls.Dial(`tcp`, address, &tls.Config{ServerName: serverName})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	var minNotAfter time.Time
	for _, chain := range conn.ConnectionState().VerifiedChains {
		for _, cert := range chain {
			if minNotAfter.Unix() > cert.NotAfter.Unix() {
				minNotAfter = cert.NotAfter
			}
		}
	}
	L.Push(lua.LNumber(minNotAfter.Unix()))
	return 1
}
