package settings

import "time"

var (
	socketPath       = "/tmp/zalua-mon.sock"
	pluginConfigPath = "/etc/zalua/config.lua"

	logPath                 = "/var/log/zabbix/zalua.log"
	maxSizeRequest          = 1024
	defaultReadTimeoutInMs  = 100
	defaultWriteTimeoutInMs = 100
)

// путь до сокета
func SocketPath() string {
	return socketPath
}

// путь до плагинов
func PluginConfigPath() string {
	return pluginConfigPath
}

// путь до файла с логами
func LogPath() string {
	return logPath
}

// чтение из сокета
func TimeoutRead() time.Duration {
	return time.Duration(defaultReadTimeoutInMs) * time.Millisecond
}

// запись в сокет
func TimeoutWrite() time.Duration {
	return time.Duration(defaultWriteTimeoutInMs) * time.Millisecond
}

// масксимальное сообщение в сокете
func MaxSizeRequest() int64 {
	return int64(maxSizeRequest)
}
