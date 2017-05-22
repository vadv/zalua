package socket

import (
	"log"
	"net"
	"os"

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
