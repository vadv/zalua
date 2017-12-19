package dsl

import (
	"errors"
	"strconv"

	lua "github.com/yuin/gopher-lua"
	yaml "gopkg.in/yaml.v2"
)

var (
	errYamlFunction = errors.New("cannot encode function to YAML")
	errYamlChannel  = errors.New("cannot encode channel to YAML")
	errYamlState    = errors.New("cannot encode state to YAML")
	errYamlUserData = errors.New("cannot encode userdata to YAML")
	errYamlNested   = errors.New("cannot encode recursively nested tables to YAML")
)

type yamlValue struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (c *dslConfig) dslYamlDecode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := yaml.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(fromYAML(L, value))
	return 1
}

func (c *dslConfig) dslYamlEncode(L *lua.LState) int {
	value := L.CheckAny(1)

	visited := make(map[*lua.LTable]bool)
	data, err := toYAML(value, visited)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}

func (j yamlValue) MarshalJSON() ([]byte, error) {
	return toYAML(j.LValue, j.visited)
}

func toYAML(value lua.LValue, visited map[*lua.LTable]bool) (data []byte, err error) {
	switch converted := value.(type) {
	case lua.LBool:
		data, err = yaml.Marshal(converted)
	case lua.LChannel:
		err = errYamlChannel
	case lua.LNumber:
		data, err = yaml.Marshal(converted)
	case *lua.LFunction:
		err = errYamlFunction
	case *lua.LNilType:
		data, err = yaml.Marshal(converted)
	case *lua.LState:
		err = errYamlState
	case lua.LString:
		data, err = yaml.Marshal(converted)
	case *lua.LTable:
		var arr []yamlValue
		var obj map[string]yamlValue

		if visited[converted] {
			panic(errYamlNested)
		}
		visited[converted] = true

		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			i, numberKey := k.(lua.LNumber)
			if numberKey && obj == nil {
				index := int(i) - 1
				if index != len(arr) {
					// map out of order; convert to map
					obj = make(map[string]yamlValue)
					for i, value := range arr {
						obj[strconv.Itoa(i+1)] = value
					}
					obj[strconv.Itoa(index+1)] = yamlValue{v, visited}
					return
				}
				arr = append(arr, yamlValue{v, visited})
				return
			}
			if obj == nil {
				obj = make(map[string]yamlValue)
				for i, value := range arr {
					obj[strconv.Itoa(i+1)] = value
				}
			}
			obj[k.String()] = yamlValue{v, visited}
		})
		if obj != nil {
			data, err = yaml.Marshal(obj)
		} else {
			data, err = yaml.Marshal(arr)
		}
	case *lua.LUserData:
		// TODO: call metatable __tostring?
		err = errYamlUserData
	}
	return
}

func fromYAML(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case int:
		return lua.LNumber(converted)
	case int64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromYAML(L, item))
		}
		return arr
	case map[interface{}]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(fromYAML(L, key), fromYAML(L, item))
		}
		return tbl
	case interface{}:
		if v, ok := converted.(bool); ok {
			return lua.LBool(v)
		}
		if v, ok := converted.(float64); ok {
			return lua.LNumber(v)
		}
		if v, ok := converted.(string); ok {
			return lua.LString(v)
		}
	}
	return lua.LNil
}
