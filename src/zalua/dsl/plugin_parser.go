package dsl

import (
	"log"
	goplugin "plugin"

	lua "github.com/yuin/gopher-lua"
)

type PluginParserInterface interface {
	ProcessData(string) (map[string]string, error)
}

type pluginParser struct {
	filename string
	parser   PluginParserInterface
}

// получение plugins из lua-state
func checkPluginParser(L *lua.LState) *pluginParser {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*pluginParser); ok {
		return v
	}
	L.ArgError(1, "plugin parser expected")
	return nil
}

// загрузка парсера
func (c *dslConfig) dslNewPluginParser(L *lua.LState) int {
	filename := L.CheckString(1)
	symbolName := `NewParser`
	if L.GetTop() > 1 {
		symbolName = L.CheckString(2)
	}
	p, err := goplugin.Open(filename)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	sym, err := p.Lookup(symbolName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	newPluginParser := sym.(PluginParserInterface)
	ud := L.NewUserData()
	ud.Value = &pluginParser{parser: newPluginParser, filename: filename}
	L.SetMetatable(ud, L.GetTypeMetatable("parser"))
	L.Push(ud)
	log.Printf("[INFO] Loaded parser plugin `%s` from `%s`\n", symbolName, filename)
	return 1
}

// выполнение парсинга
func (c *dslConfig) dslPluginParserParse(L *lua.LState) int {
	p := checkPluginParser(L)
	data := L.CheckString(2)
	t, err := p.parser.ProcessData(data)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	if t == nil {
		L.Push(lua.LNil)
		return 1
	}
	luaRow := L.CreateTable(0, len(t))
	for key, value := range t {
		luaRow.RawSetString(key, lua.LString(value))
	}
	L.Push(luaRow)
	return 1
}
