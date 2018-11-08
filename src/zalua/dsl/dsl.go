package dsl

import (
	"runtime"

	lua "github.com/yuin/gopher-lua"
)

type dslConfig struct{}

func NewConfig() *dslConfig {
	return &dslConfig{}
}

func Register(config *dslConfig, L *lua.LState) {

	plugin := L.NewTypeMetatable("plugin")
	L.SetGlobal("plugin", plugin)
	L.SetField(plugin, "new", L.NewFunction(config.dslNewPlugin))
	L.SetField(plugin, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"filename":    config.dslPluginFilename,
		"run":         config.dslPluginRun,
		"stop":        config.dslPluginStop,
		"error":       config.dslPluginError,
		"was_stopped": config.dslPluginWasStopped,
		"is_running":  config.dslPluginIsRunning,
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
	L.SetField(ioutil, "read_file", L.NewFunction(config.dslIoutilReadFile))

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
	L.SetField(time, "parse", L.NewFunction(config.dslTimeParse))

	http := L.NewTypeMetatable("http")
	L.SetGlobal("http", http)
	L.SetField(http, "get", L.NewFunction(config.dslHttpGet))
	L.SetField(http, "post", L.NewFunction(config.dslHttpPost))
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

	crypto := L.NewTypeMetatable("crypto")
	L.SetGlobal("crypto", crypto)
	L.SetField(crypto, "md5", L.NewFunction(config.dslCryptoMD5))

	json := L.NewTypeMetatable("json")
	L.SetGlobal("json", json)
	L.SetField(json, "decode", L.NewFunction(config.dslJsonDecode))
	L.SetField(json, "encode", L.NewFunction(config.dslJsonEncode))

	yaml := L.NewTypeMetatable("yaml")
	L.SetGlobal("yaml", yaml)
	L.SetField(yaml, "decode", L.NewFunction(config.dslYamlDecode))

	xmlPath := L.NewTypeMetatable("xmlpath")
	L.SetGlobal("xmlpath", xmlPath)
	L.SetField(xmlPath, "parse", L.NewFunction(config.dslXmlParse))

	cmd := L.NewTypeMetatable("cmd")
	L.SetGlobal("cmd", cmd)
	L.SetField(cmd, "exec", L.NewFunction(config.dslCmdExec))

	tlsUtil := L.NewTypeMetatable("tls_util")
	L.SetGlobal("tls_util", tlsUtil)
	L.SetField(tlsUtil, "cert_not_after", L.NewFunction(config.dslTLSUtilCertGetNotAfter))

	human := L.NewTypeMetatable("human")
	L.SetGlobal("human", human)
	L.SetField(human, "time", L.NewFunction(config.dslHumanizeTime))

	goruntime := L.NewTypeMetatable("goruntime")
	L.SetGlobal("goruntime", goruntime)
	L.SetField(goruntime, "goarch", lua.LString(runtime.GOARCH))
	L.SetField(goruntime, "goos", lua.LString(runtime.GOOS))

	regexp := L.NewTypeMetatable("regexp")
	L.SetGlobal("regexp", regexp)
	L.SetField(regexp, "compile", L.NewFunction(config.dslRegexpCompile))
	L.SetField(regexp, "match", L.NewFunction(config.dslRegexpIsMatch))
	L.SetField(regexp, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"match":           config.dslRegexpMatch,
		"find_all_string": config.dslRegexpFindAllString,
		"find_all":        config.dslRegexpFindAllString,
		"find_string":     config.dslRegexpFindString,
		"find":            config.dslRegexpFindString,
	}))

	prometheus := L.NewTypeMetatable("prometheus")
	L.SetGlobal("prometheus", prometheus)
	L.SetField(prometheus, "listen", L.NewFunction(config.dslPrometheusListen))

	prometheus_counter := L.NewTypeMetatable("prometheus_counter")
	L.SetGlobal("prometheus_counter", prometheus_counter)
	L.SetField(prometheus_counter, "new", L.NewFunction(config.dslNewPrometheusCounter))
	L.SetField(prometheus_counter, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"inc": config.dslPrometheusCounterInc,
		"add": config.dslPrometheusCounterAdd,
	}))

	prometheus_counter_vec := L.NewTypeMetatable("prometheus_counter_vec")
	L.SetGlobal("prometheus_counter_vec", prometheus_counter_vec)
	L.SetField(prometheus_counter_vec, "new", L.NewFunction(config.dslNewPrometheusCounterVec))
	L.SetField(prometheus_counter_vec, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"inc": config.dslPrometheusCounterVecInc,
		"add": config.dslPrometheusCounterVecAdd,
	}))

	prometheus_gauge := L.NewTypeMetatable("prometheus_gauge")
	L.SetGlobal("prometheus_gauge", prometheus_gauge)
	L.SetField(prometheus_gauge, "new", L.NewFunction(config.dslNewPrometheusGauge))
	L.SetField(prometheus_gauge, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"add": config.dslPrometheusGaugeAdd,
		"set": config.dslPrometheusGaugeSet,
	}))

	prometheus_gauge_vec := L.NewTypeMetatable("prometheus_gauge_vec")
	L.SetGlobal("prometheus_gauge_vec", prometheus_gauge_vec)
	L.SetField(prometheus_gauge_vec, "new", L.NewFunction(config.dslNewPrometheusGaugeVec))
	L.SetField(prometheus_gauge_vec, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"add": config.dslPrometheusGaugeVecAdd,
		"set": config.dslPrometheusGaugeVecSet,
	}))
}
