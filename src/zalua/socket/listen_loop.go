package socket

import (
	"log"
	"net"
	"os"
	"time"

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
	healTickChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-healTickChan:
			if !Alive() {
				log.Printf("[ERROR] Socket %s is not alive, exiting now...\n", settings.SocketPath())
				os.Exit(0)
			}
		}
	}
}
