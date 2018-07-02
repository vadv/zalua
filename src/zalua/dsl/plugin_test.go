package dsl

import (
	"os"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestPlugin(t *testing.T) {

	testLua1 := `
p = plugin.new("plugin_test.lua")
p:run()
time.sleep(1)
err = p:error(); if err then error(err) end

p:stop() -- нужно проверить что плагин реально остановился
time.sleep(1)
err = p:error(); if err then error(err) end
`

	testLuas2 := `
state = p:was_stopped()
if not state then error("plugin must be stopped") end
`

	if err := os.Setenv(`LOG_PATH`, os.DevNull); err != nil {
		t.Fatalf("set log: %s\n", err)
	}

	state := lua.NewState()
	Register(NewConfig(), state)
	if err := state.DoString(testLua1); err != nil {
		if !strings.Contains(err.Error(), stopPluginMessage) {
			t.Fatalf("execute lua error: %s\n", err.Error())
		}
	} else {
		t.Fatalf("must be stop error message\n")
	}

	if err := state.DoString(testLuas2); err != nil {
		t.Fatalf("execute lua error: %s\n", err.Error())
	}
}
