package dsl

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// из LState забрать 3 или 4 аргумент: table tags && int64
func getTagsTtlFromState(L *lua.LState) (map[string]string, int64) {

	ttl := int64(300)
	tags := make(map[string]string, 0)

	// check 3, 4 args (ttl or tags)
	if L.GetTop() > 2 {

		any3 := L.CheckAny(3)
		var any4 lua.LValue
		if L.GetTop() > 3 {
			any4 = L.CheckAny(4)
		} else {
			any4 = lua.LNil
		}
		var ttlSet, tagsSet bool

		// check 3
		switch any3.Type() {
		case lua.LTNumber:
			ttl = L.CheckInt64(3)
			ttlSet = true
		case lua.LTTable:
			newTags, err := tblToStringMap(L.CheckTable(3))
			if err != nil {
				L.RaiseError(fmt.Sprintf("argument #3 to table: %s", err.Error()))
			}
			tagsSet = true
			tags = newTags
		default:
			L.RaiseError("argument #3 unknown type: %v", any3.Type())
		}

		// check 4
		if any4.Type() != lua.LTNil {
			switch any4.Type() {
			case lua.LTNumber:
				if ttlSet {
					L.RaiseError("argument #4 must be table (ttl already geted)")
				}
				ttl = L.CheckInt64(4)
			case lua.LTTable:
				if tagsSet {
					L.RaiseError("argument #4 must be int (tags already geted)")
				}
				newTags, err := tblToStringMap(L.CheckTable(4))
				if err != nil {
					L.RaiseError(fmt.Sprintf("argument #4 to table: %s", err.Error()))
				}
				tags = newTags
			default:
				L.RaiseError("argument #4 unknown type: %v", any3.Type())
			}
		}

	}

	return tags, ttl
}

func stringMapsToTable(L *lua.LState, tags map[string]string) lua.LValue {
	if tags == nil {
		return lua.LNil
	}
	result := L.CreateTable(0, len(tags))
	for key, value := range tags {
		result.RawSetString(key, lua.LString(value))
	}
	return result
}

func tblToStringMap(tbl *lua.LTable) (result map[string]string, err error) {
	if tbl == nil {
		return nil, fmt.Errorf("table empty")
	}
	result = make(map[string]string, 0)
	cb := func(k, v lua.LValue) {
		if kStr, ok := k.(lua.LString); ok {
			if vStr, ok := v.(lua.LString); ok {
				result[string(kStr)] = string(vStr)
				return
			}
			if vFloat, ok := v.(lua.LNumber); ok {
				result[string(kStr)] = vFloat.String()
				return
			}
		}
		if err == nil {
			err = fmt.Errorf("can't convert map")
		}
	}
	tbl.ForEach(cb)
	return result, err
}
