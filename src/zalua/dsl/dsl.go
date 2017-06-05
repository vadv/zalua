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
		"filename":   config.dslPluginFilename,
		"run":        config.dslPluginRun,
		"stop":       config.dslPluginStop,
		"error":      config.dslPluginError,
		"is_running": config.dslPluginIsRunning,
	}))

	postgres := L.NewTypeMetatable("postgres")
	L.SetGlobal("postgres", postgres)
	L.SetField(postgres, "open", L.NewFunction(config.dslNewPgsqlConn))
	L.SetField(postgres, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"query": config.dslPgsqlQuery,
	}))

	storage := L.NewTypeMetatable("metrics")
	L.SetGlobal("metrics", storage)
	L.SetField(storage, "get", L.NewFunction(config.dslStorageGet))
	L.SetField(storage, "set", L.NewFunction(config.dslStorageSet))
	L.SetField(storage, "set_speed", L.NewFunction(config.dslStorageSetSpeed))
	L.SetField(storage, "set_counter_speed", L.NewFunction(config.dslStorageSetCounterSpeed))
	L.SetField(storage, "list", L.NewFunction(config.dslStorageList))
	L.SetField(storage, "delete", L.NewFunction(config.dslStorageDelete))

	ioutil := L.NewTypeMetatable("ioutil")
	L.SetGlobal("ioutil", ioutil)
	L.SetField(ioutil, "readfile", L.NewFunction(config.dslIoutilReadFile))

	filepath := L.NewTypeMetatable("filepath")
	L.SetGlobal("filepath", filepath)
	L.SetField(filepath, "base", L.NewFunction(config.dslFilepathBasename))
	L.SetField(filepath, "dir", L.NewFunction(config.dslFilepathDir))
	L.SetField(filepath, "ext", L.NewFunction(config.dslFilepathExt))
	L.SetField(filepath, "glob", L.NewFunction(config.dslFilepathGlob))

	os := L.NewTypeMetatable("os")
	L.SetGlobal("os", os)
	L.SetField(os, "stat", L.NewFunction(config.dslOsStat))
	L.SetField(os, "pagesize", L.NewFunction(config.dslOsPagesize))

	time := L.NewTypeMetatable("time")
	L.SetGlobal("time", time)
	L.SetField(time, "unix", L.NewFunction(config.dslTimeUnix))
	L.SetField(time, "unix_nano", L.NewFunction(config.dslTimeUnixNano))
	L.SetField(time, "sleep", L.NewFunction(config.dslTimeSleep))

	strings := L.NewTypeMetatable("strings")
	L.SetGlobal("strings", strings)
	L.SetField(strings, "split", L.NewFunction(config.dslStringsSplit))

	log := L.NewTypeMetatable("log")
	L.SetGlobal("log", log)
	L.SetField(log, "error", L.NewFunction(config.dslLogError))
	L.SetField(log, "info", L.NewFunction(config.dslLogInfo))

	json := L.NewTypeMetatable("json")
	L.SetGlobal("json", json)
	L.SetField(json, "decode", L.NewFunction(config.dslJsonDecode))
	L.SetField(json, "encode", L.NewFunction(config.dslJsonEncode))
}
