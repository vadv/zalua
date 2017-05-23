package socket

import (
	"log"
	"net"
	"os"
	"time"

	"zalua/protocol"
	"zalua/settings"
)

type server struct {
	fd net.Listener
}

func ListenLoop(clientHandle func(net.Conn)) error {
	if exists() {
		if err := os.Remove(settings.SocketPath()); err != nil {
			return err
		}
	}
	fd, err := net.Listen("unix", settings.SocketPath())
	if err != nil {
		return err
	}
	result := &server{fd: fd}
	go result.healthCheck()
	result.run(clientHandle)
	return nil
}

func (s *server) run(clientHandle func(net.Conn)) {
	log.Printf("[INFO] Start listen %s\n", settings.SocketPath())
	for {
		clientFd, err := s.fd.Accept()
		if err != nil {
			log.Printf("[ERROR] handle client: %s\n", err.Error())
		}
		go clientHandle(clientFd)
	}

}

func (s *server) healthCheck() {

	validServerId := func() bool {
		client, err := GetClient()
		if err != nil {
			return false
		}
		defer client.Close()

		if response, err := client.SendMessage(protocol.COMMAND_SERVER_ID); err != nil {
			return false
		} else {
			return response == settings.ServerId()
		}
	}

	healTickChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-healTickChan:
			if !validServerId() {
				log.Printf("[ERROR] Socket %s is not alive or run with another server, exiting now...\n", settings.SocketPath())
				os.Exit(0)
			}
		}
	}
}
