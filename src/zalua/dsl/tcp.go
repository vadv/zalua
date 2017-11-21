package dsl

import (
	"fmt"
	"log"
	"net"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type tcpConn struct {
	address string
	tcp     net.Conn
}

func (c *tcpConn) connect() error {
	conn, err := net.DialTimeout("tcp", c.address, 5*time.Second)
	if err != nil {
		return err
	}
	c.tcp = conn
	return nil
}

// получение connection из lua-state
func checkTCPConn(L *lua.LState) *tcpConn {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*tcpConn); ok {
		return v
	}
	L.ArgError(1, "tcp connection expected")
	return nil
}

// создание коннекта
func (c *dslConfig) dslNewTCPConn(L *lua.LState) int {
	addr := L.CheckString(1)
	t := &tcpConn{address: addr}
	if err := t.connect(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = t
	L.SetMetatable(ud, L.GetTypeMetatable("tcp"))
	L.Push(ud)
	log.Printf("[INFO] New tcp connection to `%s`\n", t.address)
	return 1
}

// выполнение записи
func (c *dslConfig) dslTCPWrite(L *lua.LState) int {
	conn := checkTCPConn(L)
	data := L.CheckString(2)
	count, err := conn.tcp.Write([]byte(data))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("write to `%s`: %s", conn.address, err.Error())))
		return 1
	}
	if count != len(data) {
		L.Push(lua.LString(fmt.Sprintf("write to `%s` get: %d except: %d", conn.address, count, len(data))))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

// закрытие соединения
func (c *dslConfig) dslTCPClose(L *lua.LState) int {
	conn := checkTCPConn(L)
	if conn.tcp != nil {
		log.Printf("[INFO] Close tcp connection to `%s`\n", conn.address)
		conn.tcp.Close()
	}
	return 0
}
