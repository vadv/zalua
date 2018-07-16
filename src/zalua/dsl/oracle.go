package dsl

import (
	"database/sql"
	"log"

	_ "github.com/go-goracle/goracle"
	lua "github.com/yuin/gopher-lua"
)

type oracleConn struct {
	connString string
	db         *sql.DB
}

func (o *oracleConn) connect() error {
	if o.db == nil {
		db, err := sql.Open("oracle", o.connString)
		if err != nil {
			return err
		}
		o.db = db
	}
	if err := o.db.Ping(); err != nil {
		o.db = nil
		return err
	}
	return nil
}

// создание коннекта
func (c *dslConfig) dslNewOracleConn(L *lua.LState) int {
	connString := L.CheckString(1)
	conn := &oracleConn{connString: connString}
	if err := conn.connect(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = conn
	L.SetMetatable(ud, L.GetTypeMetatable("oracle"))
	L.Push(ud)
	log.Printf("[INFO] New oracle connection `%s`\n", conn.connString)
	return 1
}

// получение connection из lua-state
func checkOracleConn(L *lua.LState) *oracleConn {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*oracleConn); ok {
		return v
	}
	L.ArgError(1, "oracle expected")
	return nil
}

// выполнение запроса
func (c *dslConfig) dslOracleQuery(L *lua.LState) int {
	conn := checkOracleConn(L)
	query := L.CheckString(2)
	sqlRows, err := conn.db.Query(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer sqlRows.Close()
	rows, err, column_count, row_count := parseRows(sqlRows, L)
	L.Push(rows)
	if err == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(err.Error()))
	}
	L.Push(column_count)
	L.Push(row_count)
	return 4
}

// закрытие соединения
func (c *dslConfig) dslOracleClose(L *lua.LState) int {
	conn := checkOracleConn(L)
	if conn.db != nil {
		conn.db.Close()
	}
	return 0
}
