package server

import (
	"log"
	"os"
	"time"

	lua "github.com/yuin/gopher-lua"

	"zalua/dsl"
	"zalua/settings"
)

func LoadPlugins() {
	go loadPlugins()
	time.Sleep(100 * time.Millisecond)
	log.Printf("[INFO] Plugins loaded\n")
}

// сама конструкция подразумевает что нам не нужен супервизор
func loadPlugins() {
	log.Printf("[INFO] Load settings file %s\n", settings.PluginConfigPath())
	state := lua.NewState()
	dsl.Register(dsl.NewConfig(), state)
	if err := state.DoFile(settings.PluginConfigPath()); err != nil {
		log.Printf("[FATAL] Settings file: %s\n", err.Error())
		os.Exit(20)
	}
	log.Printf("[INFO] Settings file loaded, nothing to do, exit now\n")
	os.Exit(0)
}
