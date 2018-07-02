package dsl

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const stopPluginMessage = `plugin was stopped`

type plugin struct {
	sync.Mutex
	state       *lua.LState
	filename    string
	startedAt   int64
	completedAt int64
	running     bool
	lastErr     error
	cancelFunc  context.CancelFunc
}

type plugins struct {
	sync.Mutex
	list map[string]*plugin
}

func (l *plugins) insertPlugin(p *plugin) {
	l.Lock()
	defer l.Unlock()
	l.list[p.getFilename()] = p
}

func (l *plugins) allPlugins() []*plugin {
	l.Lock()
	defer l.Unlock()
	result := make([]*plugin, 0)
	for _, p := range l.list {
		result = append(result, p)
	}
	return result
}

var allPlugins = &plugins{list: make(map[string]*plugin)}

// получение плагина из lua-state
func checkPlugin(L *lua.LState) *plugin {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*plugin); ok {
		return v
	}
	L.ArgError(1, "plugin expected")
	return nil
}

// получение последней ошибки
func (p *plugin) getLastError() error {
	p.Lock()
	defer p.Unlock()
	if err := p.lastErr; err != nil {
		if strings.Contains(err.Error(), `context canceled`) {
			return fmt.Errorf(stopPluginMessage)
		}
	}
	return p.lastErr
}

// получение статуса об остановке
func (p *plugin) wasStopped() bool {
	if err := p.getLastError(); err != nil && err.Error() == stopPluginMessage {
		return true
	}
	return false
}

// получение статуса - запущено или нет
func (p *plugin) getIsRunning() bool {
	p.Lock()
	defer p.Unlock()
	return p.running
}

// получение файла - запущено или нет
func (p *plugin) getFilename() string {
	p.Lock()
	defer p.Unlock()
	return p.filename
}

// запуск плагина
func (p *plugin) start() {
	p.Lock()
	state := lua.NewState()
	Register(NewConfig(), state)
	p.state = state
	p.lastErr = nil
	p.startedAt = time.Now().Unix()
	p.running = true
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.state.SetContext(ctx)
	p.cancelFunc = cancelFunc
	p.Unlock()

	p.lastErr = p.state.DoFile(p.getFilename())
	p.running = false
	p.completedAt = time.Now().Unix()
}

// создание плагина
func (c *dslConfig) dslNewPlugin(L *lua.LState) int {
	p := &plugin{filename: L.CheckString(1)}
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable("plugin"))
	L.Push(ud)
	log.Printf("[INFO] Start plugin: `%s`\n", p.filename)
	allPlugins.insertPlugin(p)
	return 1
}

// получение file name плагина
func (c *dslConfig) dslPluginFilename(L *lua.LState) int {
	p := checkPlugin(L)
	L.Push(lua.LString(p.getFilename()))
	return 1
}

// запуск плагина в отдельном стейте
func (c *dslConfig) dslPluginRun(L *lua.LState) int {
	p := checkPlugin(L)
	go p.start()

	return 0
}

// получение ошибки
func (c *dslConfig) dslPluginError(L *lua.LState) int {
	p := checkPlugin(L)
	err := p.getLastError()
	if err == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(err.Error()))
	}
	return 1
}

// получение статуса запущен или нет
func (c *dslConfig) dslPluginIsRunning(L *lua.LState) int {
	p := checkPlugin(L)
	L.Push(lua.LBool(p.getIsRunning()))
	return 1
}

// остановка плагина
func (c *dslConfig) dslPluginStop(L *lua.LState) int {
	p := checkPlugin(L)
	log.Printf("[INFO] Stop plugin: `%s`\n", p.filename)
	p.cancelFunc()
	return 0
}

// плагин был остановлен
func (c *dslConfig) dslPluginWasStopped(L *lua.LState) int {
	p := checkPlugin(L)
	L.Push(lua.LBool(p.wasStopped()))
	return 1
}

// список все плагинов
func ListOfPlugins() []string {
	result := []string{}
	for _, p := range allPlugins.allPlugins() {
		err := p.getLastError()
		errStr := "<no error>"
		if err != nil {
			errStr = fmt.Sprintf("%s", err.Error())
			errStr = strings.TrimSpace(errStr)
		}
		result = append(result, fmt.Sprintf("%s\t\t%t\t\t%v", p.getFilename(), p.getIsRunning(), errStr))
	}
	return result
}
