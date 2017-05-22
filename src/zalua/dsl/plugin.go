package dsl

import (
	"log"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var pluginsStates = make(map[string]*lua.LState, 0)
var pluginsErrors = make(map[string]error, 0)
var pluginsLock = &sync.Mutex{}

// список запущенных плагинов
func ListOfPlugins() []string {
	pluginsLock.Lock()
	defer pluginsLock.Unlock()

	result := make([]string, 0)
	for file, _ := range pluginsStates {
		result = append(result, file)
	}
	return result
}

type plugin struct {
	filename string
}

func checkPlugin(L *lua.LState) *plugin {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*plugin); ok {
		return v
	}
	L.ArgError(1, "plugin expected")
	return nil
}

// создание плагина
func (c *dslConfig) dslNewPlugin(L *lua.LState) int {
	p := &plugin{filename: L.CheckString(1)}
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable("plugin"))
	L.Push(ud)
	log.Printf("[INFO] Load new plugin `%s`\n", p.filename)
	return 1
}

// получение file name плагина
func (c *dslConfig) dslPluginFilename(L *lua.LState) int {
	p := checkPlugin(L)
	L.Push(lua.LString(p.filename))
	return 1
}

// запуск плагина в отдельном стейте
func (c *dslConfig) dslPluginRun(L *lua.LState) int {
	p := checkPlugin(L)
	go pluginStart(p.filename)

	return 0
}

// получение последней ошибки
func (c *dslConfig) dslPluginCheck(L *lua.LState) int {
	pluginsLock.Lock()
	defer pluginsLock.Unlock()

	p := checkPlugin(L)
	filename := p.filename

	_, found := pluginsStates[filename]
	if !found {
		L.RaiseError("plugin '%s' is not started", filename)
		return 0
	}
	err, ok := pluginsErrors[filename]
	if !ok {
		L.RaiseError("plugin '%s' is not activated", filename)
		return 0
	}
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	L.Push(lua.LNil)
	return 1
}

// остановка плагина
func (c *dslConfig) dslPluginStop(L *lua.LState) int {
	pluginsLock.Lock()
	defer pluginsLock.Unlock()

	p := checkPlugin(L)
	filename := p.filename

	state, found := pluginsStates[filename]
	if !found {
		L.RaiseError("plugin '%s' is not started", filename)
		return 0
	}

	state.RaiseError("stop")
	defer state.Close()

	delete(pluginsStates, filename)
	delete(pluginsErrors, filename)
	return 0
}

func pluginStart(filename string) {
	pluginsLock.Lock()
	state := lua.NewState()
	Register(NewConfig(), state)
	pluginsStates[filename] = state
	pluginsErrors[filename] = nil
	pluginsLock.Unlock()
	pluginsErrors[filename] = state.DoFile(filename)
}
