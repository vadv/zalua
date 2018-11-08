package dsl

import (
	"fmt"
	"sort"
	"strings"

	prometheus "github.com/prometheus/client_golang/prometheus"
	lua "github.com/yuin/gopher-lua"
)

var regPrometheusGaugeVecUserData = make(map[string]*lua.LUserData, 0)

type dslPrometheusGaugeVec struct {
	gauge   *prometheus.GaugeVec
	vectors []string
}

func luaTblToSlice(val lua.LValue) ([]string, error) {
	result := make([]string, 0)
	tbl, ok := val.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("bad value type: %s", val.Type().String())
	}
	tbl.ForEach(func(k lua.LValue, v lua.LValue) {
		result = append(result, v.String())
	})
	sort.Strings(result)
	return result, nil
}

func luaTblToPrometheusLabels(val lua.LValue) (prometheus.Labels, error) {
	result := make(map[string]string, 0)
	tbl, ok := val.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("bad value type: %s", val.Type().String())
	}
	tbl.ForEach(func(k lua.LValue, v lua.LValue) {
		result[k.String()] = v.String()
	})
	return result, nil
}

func (c *dslConfig) dslNewPrometheusGaugeVec(L *lua.LState) int {

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

	if ud, ok := regPrometheusGaugeVecUserData[fullName]; ok {
		gauge := ud.Value.(*dslPrometheusGaugeVec)
		if strings.Join(gauge.vectors, "_") != strings.Join(vectors, "_") {
			L.Push(lua.LNil)
			L.Push(lua.LString("can't change vectors online"))
			return 2
		}
		L.Push(ud)
		return 1
	}

	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		Help:      config.RawGetString(`help`).String(),
	}, vectors)
	if err := prometheus.Register(gauge); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = &dslPrometheusGaugeVec{gauge: gauge, vectors: vectors}
	L.SetMetatable(ud, L.GetTypeMetatable("prometheus_gauge_vec"))
	L.Push(ud)
	regPrometheusGaugeVecUserData[fullName] = ud
	return 1
}

func checkPrometheusGaugeVec(L *lua.LState) *dslPrometheusGaugeVec {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*dslPrometheusGaugeVec); ok {
		return v
	}
	L.ArgError(1, "prometheus_gauge_vec expected")
	return nil
}

func (c *dslConfig) dslPrometheusGaugeVecAdd(L *lua.LState) int {
	gauge := checkPrometheusGaugeVec(L)
	luaLabels := L.CheckTable(2)
	value := L.CheckNumber(3)
	labels, err := luaTblToPrometheusLabels(luaLabels)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}
	gauge.gauge.With(labels).Add(float64(value))
	return 0
}

func (c *dslConfig) dslPrometheusGaugeVecSet(L *lua.LState) int {
	gauge := checkPrometheusGaugeVec(L)
	luaLabels := L.CheckTable(2)
	value := L.CheckNumber(3)
	labels, err := luaTblToPrometheusLabels(luaLabels)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	}
	gauge.gauge.With(labels).Set(float64(value))
	return 0
}
