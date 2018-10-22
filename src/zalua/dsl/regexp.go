package dsl

import (
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

func checkRegexp(L *lua.LState) *regexp.Regexp {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*regexp.Regexp); ok {
		return v
	}
	return nil
}

// создание regexp
func (c *dslConfig) dslRegexpCompile(L *lua.LState) int {
	expr := L.CheckString(1)
	reg, err := regexp.Compile(expr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = reg
	L.SetMetatable(ud, L.GetTypeMetatable(`regexp`))
	L.Push(ud)
	return 1
}

func (c *dslConfig) dslRegexpIsMatch(L *lua.LState) int {
	str := L.CheckString(1)
	expr := L.CheckString(2)
	reg, err := regexp.Compile(expr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(reg.MatchString(str)))
	return 1
}

// match
func (c *dslConfig) dslRegexpMatch(L *lua.LState) int {
	reg := checkRegexp(L)
	str := L.CheckString(2)
	if reg == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("regexp is nil"))
		return 2
	}
	L.Push(lua.LBool(reg.MatchString(str)))
	return 1
}

// FindAllString
func (c *dslConfig) dslRegexpFindAllString(L *lua.LState) int {
	reg := checkRegexp(L)
	if reg == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("regexp is nil"))
		return 2
	}
	str := L.CheckString(2)
	count := -1
	if L.GetTop() > 2 {
		count = int(L.CheckNumber(3))
	}
	luaList := L.NewTable()
	for _, str := range reg.FindAllString(str, count) {
		luaList.Append(lua.LString(str))
	}
	L.Push(luaList)
	return 1
}

// FindString
func (c *dslConfig) dslRegexpFindString(L *lua.LState) int {
	reg := checkRegexp(L)
	str := L.CheckString(2)
	if reg == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("regexp is nil"))
		return 2
	}
	L.Push(lua.LString(reg.FindString(str)))
	return 1
}
