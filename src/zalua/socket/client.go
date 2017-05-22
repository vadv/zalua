package socket

import (
	"net"
	"time"

	"zalua/settings"
)

type client struct {
	conn net.Conn
}

func GetClient() (*client, error) {
	conn, err := net.DialTimeout("unix", settings.SocketPath(), 200*time.Millisecond)
	if err != nil {
		return nil, err
	}
	result := &client{conn: conn}
	return result, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) write(msg string) (err error) {
	c.conn.SetWriteDeadline(time.Now().Add(settings.TimeoutWrite()))
	_, err = c.conn.Write([]byte(msg))
	return
}

func (c *client) read() (string, error) {
	buf := make([]byte, settings.MaxSizeRequest())
	c.conn.SetReadDeadline(time.Now().Add(time.Duration(settings.TimeoutWrite())))
	n, err := c.conn.Read(buf[:])
	if err != nil {
		return "", err
	}
	result := string(buf[0:n])
	return result, nil
}

func (c *client) SendMessage(msg string) (string, error) {
	if err := c.write(msg); err != nil {
		return "", err
	}
	return c.read()
}
