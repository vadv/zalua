package dsl

import (
	"fmt"

	prometheus "github.com/prometheus/client_golang/prometheus"
	lua "github.com/yuin/gopher-lua"
)

var regPrometheusCounter = make(map[string]*lua.LUserData, 0)

type dslPrometheusCounter struct {
	prometheus.Counter
}

func (c *dslConfig) dslNewPrometheusCounter(L *lua.LState) int {
	config := L.CheckTable(1)

	namespace := ""
	if config.RawGetString(`namespace`).Type() != lua.LTNil {
		namespace = config.RawGetString(`namespace`).String()
	}

	subsystem := ""
	if config.RawGetString(`subsystem`).Type() != lua.LTNil {
		subsystem = config.RawGetString(`subsystem`).String()
	}

	name := config.RawGetString(`name`).String()

	fullName := fmt.Sprintf("%s_%s_%s", namespace, subsystem, name)

	if ud, ok := regPrometheusCounter[fullName]; ok {
		L.Push(ud)
		return 1
	}

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      config.RawGetString(`help`).String(),
	})
	if err := prometheus.Register(counter); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = &dslPrometheusCounter{counter}
	L.SetMetatable(ud, L.GetTypeMetatable("prometheus_counter"))
	L.Push(ud)
	regPrometheusCounter[fullName] = ud
	return 1
}

func checkPrometheusCounter(L *lua.LState) *dslPrometheusCounter {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*dslPrometheusCounter); ok {
		return v
	}
	L.ArgError(1, "prometheus_counter expected")
	return nil
}

func (c *dslConfig) dslPrometheusCounterInc(L *lua.LState) int {
	counter := checkPrometheusCounter(L)
	counter.Inc()
	return 0
}

func (c *dslConfig) dslPrometheusCounterAdd(L *lua.LState) int {
	counter := checkPrometheusCounter(L)
	value := L.CheckNumber(2)
	counter.Add(float64(value))
	return 0
}
