package dsl

import (
	"encoding/json"
	"errors"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

var (
	errJsonFunction = errors.New("cannot encode function to JSON")
	errJsonChannel  = errors.New("cannot encode channel to JSON")
	errJsonState    = errors.New("cannot encode state to JSON")
	errJsonUserData = errors.New("cannot encode userdata to JSON")
	errJsonNested   = errors.New("cannot encode recursively nested tables to JSON")
)

type jsonValue struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (c *dslConfig) dslJsonDecode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := json.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(fromJSON(L, value))
	return 1
}

func (c *dslConfig) dslJsonEncode(L *lua.LState) int {
	value := L.CheckAny(1)

	visited := make(map[*lua.LTable]bool)
	data, err := toJSON(value, visited)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}

func (j jsonValue) MarshalJSON() ([]byte, error) {
	return toJSON(j.LValue, j.visited)
}

func toJSON(value lua.LValue, visited map[*lua.LTable]bool) (data []byte, err error) {
	switch converted := value.(type) {
	case lua.LBool:
		data, err = json.Marshal(converted)
	case lua.LChannel:
		err = errJsonChannel
	case lua.LNumber:
		data, err = json.Marshal(converted)
	case *lua.LFunction:
		err = errJsonFunction
	case *lua.LNilType:
		data, err = json.Marshal(converted)
	case *lua.LState:
		err = errJsonState
	case lua.LString:
		data, err = json.Marshal(converted)
	case *lua.LTable:
		var arr []jsonValue
		var obj map[string]jsonValue

		if visited[converted] {
			panic(errJsonNested)
		}
		visited[converted] = true

		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			i, numberKey := k.(lua.LNumber)
			if numberKey && obj == nil {
				index := int(i) - 1
				if index != len(arr) {
					// map out of order; convert to map
					obj = make(map[string]jsonValue)
					for i, value := range arr {
						obj[strconv.Itoa(i+1)] = value
					}
					obj[strconv.Itoa(index+1)] = jsonValue{v, visited}
					return
				}
				arr = append(arr, jsonValue{v, visited})
				return
			}
			if obj == nil {
				obj = make(map[string]jsonValue)
				for i, value := range arr {
					obj[strconv.Itoa(i+1)] = value
				}
			}
			obj[k.String()] = jsonValue{v, visited}
		})
		if obj != nil {
			data, err = json.Marshal(obj)
		} else {
			data, err = json.Marshal(arr)
		}
	case *lua.LUserData:
		// TODO: call metatable __tostring?
		err = errJsonUserData
	}
	return
}

func fromJSON(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromJSON(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), fromJSON(L, item))
		}
		return tbl
	}
	return lua.LNil
}
