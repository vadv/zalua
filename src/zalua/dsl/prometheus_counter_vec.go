package dsl

import (
	"fmt"
	"strings"

	prometheus "github.com/prometheus/client_golang/prometheus"
	lua "github.com/yuin/gopher-lua"
)

var regPrometheusCounterVecUserData = make(map[string]*lua.LUserData, 0)

type dslPrometheusCounterVec struct {
	counter *prometheus.CounterVec
	vectors []string
}

func (c *dslConfig) dslNewPrometheusCounterVec(L *lua.LState) int {

	config := L.CheckTable(1)
	vec := config.RawGetString(`vec`)
	vectors, err := luaTblToSlice(vec)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

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

	if ud, ok := regPrometheusCounterVecUserData[fullName]; ok {
		gauge := ud.Value.(*dslPrometheusCounterVec)
		if strings.Join(gauge.vectors, "_") != strings.Join(vectors, "_") {
			L.Push(lua.LNil)
			L.Push(lua.LString("can't change vectors online"))
			return 2
		}
		L.Push(ud)
		return 1
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      config.RawGetString(`help`).String(),
	}, vectors)
	if err := prometheus.Register(counter); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = &dslPrometheusCounterVec{counter: counter, vectors: vectors}
	L.SetMetatable(ud, L.GetTypeMetatable("prometheus_counter_vec"))
	L.Push(ud)
	regPrometheusCounterVecUserData[fullName] = ud
	return 1
}

func checkPrometheusCounterVec(L *lua.LState) *dslPrometheusCounterVec {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*dslPrometheusCounterVec); ok {
		return v
	}
	L.ArgError(1, "prometheus_counter_vec expected")
	return nil
}

func (c *dslConfig) dslPrometheusCounterVecAdd(L *lua.LState) int {
	counter := checkPrometheusCounterVec(L)
	luaLabels := L.CheckTable(2)
	value := L.CheckNumber(3)
	labels, err := luaTblToPrometheusLabels(luaLabels)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}
	counter.counter.With(labels).Add(float64(value))
	return 0
}

func (c *dslConfig) dslPrometheusCounterVecInc(L *lua.LState) int {
	counter := checkPrometheusCounterVec(L)
	luaLabels := L.CheckTable(2)
	labels, err := luaTblToPrometheusLabels(luaLabels)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}
	counter.counter.With(labels).Inc()
	return 0
}
