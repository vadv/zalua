package dsl

import (
	"bytes"

	lua "github.com/yuin/gopher-lua"
	xmlpath "gopkg.in/xmlpath.v2"
)

func (d *dslConfig) dslXmlParse(L *lua.LState) int {

	xmlData := L.CheckString(1)
	xmlPath := L.CheckString(2)

	r := bytes.NewReader([]byte(xmlData))
	node, err := xmlpath.Parse(r)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	path, err := xmlpath.Compile(xmlPath)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	i := 1
	it, result := path.Iter(node), L.NewTable()
	for it.Next() {
		L.RawSetInt(result, i, lua.LString(it.Node().String()))
		i++
	}
	L.Push(result)
	return 1
}
