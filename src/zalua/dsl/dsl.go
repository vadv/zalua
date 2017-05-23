package dsl

import lua "github.com/yuin/gopher-lua"

type dslConfig struct{}

func NewConfig() *dslConfig {
	return &dslConfig{}
}

func Register(config *dslConfig, L *lua.LState) {

	plugin := L.NewTypeMetatable("plugin")
	L.SetGlobal("plugin", plugin)
	L.SetField(plugin, "new", L.NewFunction(config.dslNewPlugin))
	L.SetField(plugin, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"filename": config.dslPluginFilename,
		"run":      config.dslPluginRun,
		"stop":     config.dslPluginStop,
		"check":    config.dslPluginCheck,
	}))

	storage := L.NewTypeMetatable("metrics")
	L.SetGlobal("metrics", storage)
	L.SetField(storage, "get", L.NewFunction(config.dslStorageGet))
	L.SetField(storage, "set", L.NewFunction(config.dslStorageSet))
	L.SetField(storage, "set_speed", L.NewFunction(config.dslStorageSetSpeed))
	L.SetField(storage, "set_counter_speed", L.NewFunction(config.dslStorageSetCounterSpeed))
	L.SetField(storage, "list", L.NewFunction(config.dslStorageList))
	L.SetField(storage, "delete", L.NewFunction(config.dslStorageDelete))

	utils := L.NewTypeMetatable("utils")
	L.SetGlobal("utils", utils)
	L.SetField(utils, "sleep", L.NewFunction(config.dslSleep))

	filepath := L.NewTypeMetatable("filepath")
	L.SetGlobal("filepath", filepath)
	L.SetField(filepath, "base", L.NewFunction(config.dslFilepathBasename))
	L.SetField(filepath, "dir", L.NewFunction(config.dslFilepathDir))
	L.SetField(filepath, "ext", L.NewFunction(config.dslFilepathExt))
	L.SetField(filepath, "glob", L.NewFunction(config.dslFilepathGlob))

	log := L.NewTypeMetatable("log")
	L.SetGlobal("log", log)
	L.SetField(log, "error", L.NewFunction(config.dslLogError))
	L.SetField(log, "info", L.NewFunction(config.dslLogInfo))

	json := L.NewTypeMetatable("json")
	L.SetGlobal("json", json)
	L.SetField(json, "decode", L.NewFunction(config.dslJsonDecode))
	L.SetField(json, "encode", L.NewFunction(config.dslJsonEncode))
}
