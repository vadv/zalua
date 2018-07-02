package settings

import (
	"fmt"
	"os"
	"time"
)

var (
	socketPath = "/tmp/zalua-mon.sock"
	initPath   = "/etc/zalua/init.lua"

	logPath                 = "/var/log/zabbix/zalua.log"
	storagePath             = "/tmp/zalua-storage.json"
	maxSizeRequest          = 64 * 1024
	defaultReadTimeoutInMs  = 100
	defaultWriteTimeoutInMs = 100
)

// путь до сокета
func SocketPath() string {
	return socketPath
}

// путь до плагинов
func InitPath() string {
	if path := os.Getenv("INIT_FILE"); path != `` {
		initPath = path
	}
	return initPath
}

// путь до файла с логами
func LogPath() string {
	if path := os.Getenv("LOG_PATH"); path != `` {
		logPath = path
	}
	return logPath
}

// путь до файла с временным стораджем
func StoragePath() string {
	return storagePath
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

var random = time.Now().UnixNano()

func ServerId() string {
	return fmt.Sprintf("server-id-%d", random)
}
