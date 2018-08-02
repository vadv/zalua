package dsl

import (
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"

	"zalua/storage"
)

func (c *dslConfig) dslStorageGet(L *lua.LState) int {
	key := L.CheckString(1)
	tags := make(map[string]string, 0)
	if L.GetTop() > 1 {
		tbl := L.CheckTable(2)
		newTags, err := tblToStringMap(tbl)
		if err != nil {
			L.RaiseError("argument #2 must be table string:string")
		}
		tags = newTags
	}

	if value, ok := storage.Box.Get(key, tags); ok {
		L.Push(lua.LString(value.GetValue()))
		L.Push(stringMapsToTable(L, value.GetTags()))
	} else {
		L.Push(lua.LNil)
		L.Push(lua.LNil)
	}
	return 2
}

var speedCacheLastValue = make(map[string]float64, 0)
var speedCacheLastTime = make(map[string]int64, 0)
var speedCacheLock = &sync.Mutex{}

func setSpeed(L *lua.LState, counter bool) int {
	speedCacheLock.Lock()
	defer speedCacheLock.Unlock()

	metric := L.CheckString(1)
	luaVal := L.CheckAny(2)
	val := float64(0)
	tags, ttl := getTagsTtlFromState(L)

	// парсим как строку
	if luaStr, ok := luaVal.(lua.LString); ok {
		if floatVal, err := strconv.ParseFloat(string(luaStr), 64); err == nil {
			val = floatVal
		}
	} else {
		// парсим как float
		if luaFloat, ok := luaVal.(lua.LNumber); ok {
			val = float64(luaFloat)
		} else {
			L.RaiseError("argument #2 must be string or number")
		}
	}

	metricKey := metric
	if len(tags) > 0 {
		data, err := json.Marshal(&tags)
		if err == nil {
			metricKey = metricKey + string(data)
		}
	}

	if lastValue, ok := speedCacheLastValue[metricKey]; ok {
		if lastTime, ok := speedCacheLastTime[metricKey]; ok {
			now := time.Now().UnixNano()
			diff := float64(val - lastValue)
			if counter && diff < 0 {
				// должны пропустить если counter и счетчик провернулся
			} else {
				value := float64(time.Second) * diff / float64(now-lastTime)
				valueStr := strconv.FormatFloat(value, 'f', 2, 64)
				if math.Abs(value) < 0.01 {
					valueStr = strconv.FormatFloat(value, 'f', 4, 64)
				}
				storage.Box.Set(metric, valueStr, tags, ttl)
			}
		}
	}
	speedCacheLastValue[metricKey] = val
	speedCacheLastTime[metricKey] = time.Now().UnixNano()

	return 0
}

func (c *dslConfig) dslStorageSetSpeed(L *lua.LState) int {
	return setSpeed(L, false)
}

func (c *dslConfig) dslStorageSetCounterSpeed(L *lua.LState) int {
	return setSpeed(L, true)
}

func (c *dslConfig) dslStorageSet(L *lua.LState) int {
	key := L.CheckString(1)
	luaVal := L.CheckAny(2)
	val := ""
	if luaStr, ok := luaVal.(lua.LString); ok {
		val = string(luaStr)
	} else {
		if luaFloat, ok := luaVal.(lua.LNumber); ok {
			value := float64(luaFloat)
			val = strconv.FormatFloat(value, 'f', 2, 64)
			if math.Abs(value) < 0.01 {
				val = strconv.FormatFloat(value, 'f', 6, 64)
			}
		} else {
			L.RaiseError("argument #2 must be string or number")
		}
	}
	tags, ttl := getTagsTtlFromState(L)
	storage.Box.Set(key, val, tags, ttl)
	return 0
}

func (c *dslConfig) dslStorageList(L *lua.LState) int {
	list := storage.Box.List()
	result := L.CreateTable(len(list), 0)
	for _, item := range list {
		t := L.CreateTable(4, 0)
		L.SetField(t, "metric", lua.LString(item.GetMetric()))
		L.SetField(t, "value", lua.LString(item.GetValue()))
		L.SetField(t, "tags", stringMapsToTable(L, item.GetTags()))
		L.SetField(t, "at", lua.LNumber(item.GetCreatedAt()))
		result.Append(t)
	}
	L.Push(result)
	return 1
}

func (c *dslConfig) dslStorageDelete(L *lua.LState) int {
	return 0
}
