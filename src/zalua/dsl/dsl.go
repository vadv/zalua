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

	tacScanner := L.NewTypeMetatable("tac")
	L.SetGlobal("tac", tacScanner)
	L.SetField(tacScanner, "open", L.NewFunction(config.dslTacOpen))
	L.SetField(tacScanner, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"line":  config.dslTacLine,
		"close": config.dslTacClose,
	}))

	postgres := L.NewTypeMetatable("postgres")
	L.SetGlobal("postgres", postgres)
	L.SetField(postgres, "open", L.NewFunction(config.dslNewPgsqlConn))
	L.SetField(postgres, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"close": config.dslPgsqlClose,
		"query": config.dslPgsqlQuery,
	}))

	tcp := L.NewTypeMetatable("tcp")
	L.SetGlobal("tcp", tcp)
	L.SetField(tcp, "open", L.NewFunction(config.dslNewTCPConn))
	L.SetField(tcp, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"close": config.dslTCPClose,
		"write": config.dslTCPWrite,
	}))

	dslPluginParser := L.NewTypeMetatable("parser")
	L.SetGlobal("parser", dslPluginParser)
	L.SetField(dslPluginParser, "load", L.NewFunction(config.dslNewPluginParser))
	L.SetField(dslPluginParser, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"parse": config.dslPluginParserParse,
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

	os := L.NewTypeMetatable("goos")
	L.SetGlobal("goos", os)
	L.SetField(os, "stat", L.NewFunction(config.dslOsStat))
	L.SetField(os, "pagesize", L.NewFunction(config.dslOsPagesize))

	time := L.NewTypeMetatable("time")
	L.SetGlobal("time", time)
	L.SetField(time, "unix", L.NewFunction(config.dslTimeUnix))
	L.SetField(time, "unix_nano", L.NewFunction(config.dslTimeUnixNano))
	L.SetField(time, "sleep", L.NewFunction(config.dslTimeSleep))

	http := L.NewTypeMetatable("http")
	L.SetGlobal("http", http)
	L.SetField(http, "get", L.NewFunction(config.dslHttpGet))
	L.SetField(http, "escape", L.NewFunction(config.dslHttpEscape))
	L.SetField(http, "unescape", L.NewFunction(config.dslHttpUnEscape))

	strings := L.NewTypeMetatable("strings")
	L.SetGlobal("strings", strings)
	L.SetField(strings, "split", L.NewFunction(config.dslStringsSplit))
	L.SetField(strings, "has_prefix", L.NewFunction(config.dslStringsHasPrefix))
	L.SetField(strings, "has_suffix", L.NewFunction(config.dslStringsHasSuffix))
	L.SetField(strings, "trim", L.NewFunction(config.dslStringsTrim))

	log := L.NewTypeMetatable("log")
	L.SetGlobal("log", log)
	L.SetField(log, "error", L.NewFunction(config.dslLogError))
	L.SetField(log, "info", L.NewFunction(config.dslLogInfo))

	json := L.NewTypeMetatable("json")
	L.SetGlobal("json", json)
	L.SetField(json, "decode", L.NewFunction(config.dslJsonDecode))
	L.SetField(json, "encode", L.NewFunction(config.dslJsonEncode))

	yaml := L.NewTypeMetatable("yaml")
	L.SetGlobal("yaml", yaml)
	L.SetField(yaml, "decode", L.NewFunction(config.dslYamlDecode))
	L.SetField(yaml, "encode", L.NewFunction(config.dslYamlEncode))

	xmlPath := L.NewTypeMetatable("xmlpath")
	L.SetGlobal("xmlpath", xmlPath)
	L.SetField(xmlPath, "parse", L.NewFunction(config.dslXmlParse))
}
