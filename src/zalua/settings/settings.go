package settings

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	socketPath = "/tmp/zalua-mon.sock"
	initPath   = "/etc/zalua/config.lua"

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
	return initPath
}

// путь до файла с логами
func LogPath() string {
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

var (
	random   = rand.New(rand.NewSource(time.Now().Unix()))
	serverId = random.Int63n(10000)
)

func ServerId() string {
	return fmt.Sprintf("server-id-%d", serverId)
}
