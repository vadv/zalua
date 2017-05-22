package socket

import (
	"os"

	"zalua/protocol"
	"zalua/settings"
)

func exists() bool {
	if _, err := os.Stat(settings.SocketPath()); err == nil {
		return true
	}
	return false
}

// проверка живой или не живой сокет
func Alive() bool {

	if !exists() {
		return false
	}

	client, err := GetClient()
	if err != nil {
		return false
	}
	defer client.Close()

	if response, err := client.SendMessage(protocol.PING); err != nil {
		return false
	} else {
		return response == protocol.PONG
	}

}
