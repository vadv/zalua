package dsl

import (
	"bytes"

	lua "github.com/yuin/gopher-lua"
	xmlpath "gopkg.in/xmlpath.v2"
)

func (d *dslConfig) dslXmlLoad(L *lua.LState) int {
	xmlStr := L.CheckString(1)
	r := bytes.NewReader([]byte(xmlStr))
	node, err := xmlpath.ParseHTML(r)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(newXmlNode(L, node))
	return 1
}

func (d *dslConfig) dslXmlCompile(L *lua.LState) int {
	xpathStr := L.CheckString(1)
	path, err := xmlpath.Compile(xpathStr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(newXmlPath(L, path))
	return 1
}

type typeXmlNode struct {
	base *xmlpath.Node
}
type typeXmlPath struct {
	base *xmlpath.Path
}

type typeXmlIter struct {
	base *xmlpath.Iter
}

const luaNodeTypeName = "xmlpath.node"
const luaPathTypeName = "xmlpath.path"
const luaIterTypeName = "xmlpath.iter"

func registerXmlType(L *lua.LState, module *lua.LTable) {

	//reg node
	nodemt := L.NewTypeMetatable(luaNodeTypeName)
	L.SetField(module, "node", nodemt)
	L.SetField(nodemt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"string": nodeXmlString,
	}))

	//reg path
	pathmt := L.NewTypeMetatable(luaPathTypeName)
	L.SetField(module, "path", pathmt)
	L.SetField(pathmt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"iter": xmlIter,
	}))

	//reg iter
	itermt := L.NewTypeMetatable(luaIterTypeName)
	L.SetField(module, "iter", itermt)
	L.SetField(itermt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"next": xmlNext,
		"node": xmlNode,
	}))

}

func newXmlNode(L *lua.LState, n *xmlpath.Node) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &typeXmlNode{
		n,
	}
	L.SetMetatable(ud, L.GetTypeMetatable(luaNodeTypeName))
	return ud
}

func checkXmlNode(L *lua.LState) *typeXmlNode {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*typeXmlNode); ok {
		return v
	}
	L.ArgError(1, "node expected")
	return nil
}

func newXmlPath(L *lua.LState, p *xmlpath.Path) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &typeXmlPath{
		p,
	}
	L.SetMetatable(ud, L.GetTypeMetatable(luaPathTypeName))
	return ud
}

func checkXmlPath(L *lua.LState) *typeXmlPath {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*typeXmlPath); ok {
		return v
	}
	L.ArgError(1, "path expected")
	return nil
}

func newXmlIter(L *lua.LState, i *xmlpath.Iter) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = &typeXmlIter{
		i,
	}
	L.SetMetatable(ud, L.GetTypeMetatable(luaIterTypeName))
	return ud
}

func checkXmlIter(L *lua.LState) *typeXmlIter {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*typeXmlIter); ok {
		return v
	}
	L.ArgError(1, "iter expected")
	return nil
}

//iter := path.iter(node)
func xmlIter(L *lua.LState) int {
	path := checkXmlPath(L)
	if L.GetTop() == 2 {
		ut := L.CheckUserData(2)
		if node, ok := ut.Value.(*typeXmlNode); ok {
			it := path.base.Iter(node.base)
			ltab := L.NewTable()
			i := 1
			for it.Next() {
				L.RawSetInt(ltab, i, newXmlNode(L, it.Node()))
				i++
			}
			L.Push(ltab)
			//L.Push(newXmlIter(L, it))
			return 1
		}
	}
	L.ArgError(1, "node expected")
	return 0
}

//support lua standard iterator
//hasNext := iter.next()
func xmlNext(L *lua.LState) int {
	iter := checkXmlIter(L)
	L.Push(lua.LBool(iter.base.Next()))
	return 1
}

//node := iter.node()
func xmlNode(L *lua.LState) int {
	iter := checkXmlIter(L)
	L.Push(newXmlNode(L, iter.base.Node()))
	return 1
}

//string := node.string()
func nodeXmlString(L *lua.LState) int {
	node := checkXmlNode(L)
	L.Push(lua.LString(node.base.String()))
	return 1
}
