package dsl

import (
	"log"
	"net/http"

	promhttp "github.com/prometheus/client_golang/prometheus/promhttp"
	lua "github.com/yuin/gopher-lua"
)

type prometheusSrv struct {
	listenAddr string
}

func (p *prometheusSrv) start(L *lua.LState) {
	log.Printf("[INFO] start prometheus listener at %s\n", p.listenAddr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(p.listenAddr, nil); err != nil {
		L.RaiseError("listen prometheus: %s", err.Error())
	}
}

func (c *dslConfig) dslPrometheusListen(L *lua.LState) int {
	listenAddr := L.CheckString(1)
	p := &prometheusSrv{listenAddr: listenAddr}
	go p.start(L)
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable("prometheus"))
	L.Push(ud)
	return 1
}

func checkPrometheus(L *lua.LState) *prometheusSrv {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*prometheusSrv); ok {
		return v
	}
	L.ArgError(1, "prometheus expected")
	return nil
}
