package dsl

import (
	"fmt"

	prometheus "github.com/prometheus/client_golang/prometheus"
	lua "github.com/yuin/gopher-lua"
)

var regPrometheusGaugeUserData = make(map[string]*lua.LUserData, 0)

type dslPrometheusGauge struct {
	prometheus.Gauge
}

func (c *dslConfig) dslNewPrometheusGauge(L *lua.LState) int {

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

	if ud, ok := regPrometheusGaugeUserData[fullName]; ok {
		L.Push(ud)
		return 1
	}

	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      config.RawGetString(`help`).String(),
	})
	if err := prometheus.Register(gauge); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = &dslPrometheusGauge{gauge}
	L.SetMetatable(ud, L.GetTypeMetatable("prometheus_gauge"))
	L.Push(ud)
	regPrometheusGaugeUserData[fullName] = ud
	return 1
}

func checkPrometheusGauge(L *lua.LState) *dslPrometheusGauge {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*dslPrometheusGauge); ok {
		return v
	}
	L.ArgError(1, "prometheus_gauge expected")
	return nil
}

func (c *dslConfig) dslPrometheusGaugeAdd(L *lua.LState) int {
	gauge := checkPrometheusGauge(L)
	value := L.CheckNumber(2)
	gauge.Add(float64(value))
	return 0
}

func (c *dslConfig) dslPrometheusGaugeSet(L *lua.LState) int {
	gauge := checkPrometheusGauge(L)
	value := L.CheckNumber(2)
	gauge.Set(float64(value))
	return 0
}
