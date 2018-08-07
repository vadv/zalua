package dsl

import (
	"crypto/tls"
	"fmt"
	"net"
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
	conn, err := net.DialTimeout(`tcp`, address, 5*time.Second)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	client := tls.Client(conn, &tls.Config{ServerName: serverName})
	handshakeErr := client.Handshake()

	var minNotAfter *time.Time
	for _, chain := range client.ConnectionState().VerifiedChains {
		for _, cert := range chain {
			if minNotAfter == nil || minNotAfter.Unix() > cert.NotAfter.Unix() {
				minNotAfter = &cert.NotAfter
			}
		}
	}

	L.Push(lua.LNumber(minNotAfter.Unix()))
	if handshakeErr == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(handshakeErr.Error()))
	}
	return 2
}
