package dsl

import (
	"strconv"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"

	"zalua/storage"
)

func (c *dslConfig) dslStorageGet(L *lua.LState) int {
	key := L.CheckString(1)
	if value, ok := storage.Box.Get(key); ok {
		L.Push(lua.LString(value))
	} else {
		L.Push(lua.LNil)
	}
	return 1
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
	ttl := int64(300)
	if L.GetTop() == 3 {
		ttl = L.CheckInt64(3)
	}
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

	if lastValue, ok := speedCacheLastValue[metric]; ok {
		if lastTime, ok := speedCacheLastTime[metric]; ok {
			now := time.Now().UnixNano()
			diff := float64(val - lastValue)
			if counter && diff < 0 {
				// должны пропустить если counter и счетчик провернулся
			} else {
				value := float64(time.Second) * diff / float64(now-lastTime)
				storage.Box.Set(metric, strconv.FormatFloat(value, 'f', 2, 64), ttl)
			}
		}
	}
	speedCacheLastValue[metric] = val
	speedCacheLastTime[metric] = time.Now().UnixNano()

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
			val = strconv.FormatFloat(float64(luaFloat), 'f', 2, 64)
		} else {
			L.RaiseError("argument #2 must be string or number")
		}
	}
	ttl := int64(300)
	if L.GetTop() == 3 {
		ttl = L.CheckInt64(3)
	}
	storage.Box.Set(key, val, ttl)
	return 0
}

func (c *dslConfig) dslStorageList(L *lua.LState) int {
	list := storage.Box.List()
	result := L.CreateTable(len(list), 0)
	for _, key := range list {
		result.Append(lua.LString(key))
	}
	L.Push(result)
	return 1
}

func (c *dslConfig) dslStorageDelete(L *lua.LState) int {
	return 0
}
