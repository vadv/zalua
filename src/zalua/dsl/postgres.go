package dsl

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	lua "github.com/yuin/gopher-lua"
)

type pgsqlConn struct {
	host     string
	database string
	user     string
	passwd   string
	port     int
	db       *sql.DB
}

func (p *pgsqlConn) connectionString() string {
	return fmt.Sprintf("host='%s' port='%d' user='%s' dbname='%s' password='%s' sslmode='disable' fallback_application_name='zalua' connect_timeout='5'",
		p.host, p.port, p.user, p.database, p.passwd)
}

func (p *pgsqlConn) connect() error {
	if p.db == nil {
		db, err := sql.Open("postgres", p.connectionString())
		if err != nil {
			return err
		}
		p.db = db
	}
	if err := p.db.Ping(); err != nil {
		p.db = nil
		return err
	}
	return nil
}

// получение connection из lua-state
func checkPgsqlConn(L *lua.LState) *pgsqlConn {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*pgsqlConn); ok {
		return v
	}
	L.ArgError(1, "postgres expected")
	return nil
}

// парсит sql.rows и  `rows, err, column_count, row_count = db:query()` выполнить запрос
func parseRows(sqlRows *sql.Rows, L *lua.LState) (luaRows *lua.LTable, resultErr error, luaColumnCount lua.LNumber, luaRowCount lua.LNumber) {
	cols, err := sqlRows.Columns()
	if err != nil {
		resultErr = err
		return
	}
	// пробегаем по строкам
	luaRows = L.CreateTable(0, len(cols))
	rowCount := 1
	for sqlRows.Next() {
		columns := make([]interface{}, len(cols))
		pointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			pointers[i] = &columns[i]
		}
		err := sqlRows.Scan(pointers...)
		if err != nil {
			resultErr = err
			return
		}
		luaRow := L.CreateTable(0, len(cols))
		for i, _ := range cols {
			valueP := pointers[i].(*interface{})
			value := *valueP
			switch converted := value.(type) {
			case bool:
				luaRow.RawSetInt(i+1, lua.LBool(converted))
			case float64:
				luaRow.RawSetInt(i+1, lua.LNumber(converted))
			case int64:
				luaRow.RawSetInt(i+1, lua.LNumber(converted))
			case string:
				luaRow.RawSetInt(i+1, lua.LString(converted))
			case []byte:
				luaRow.RawSetInt(i+1, lua.LString(string(converted)))
			case nil:
				luaRow.RawSetInt(i+1, lua.LNil)
			default:
				log.Printf("[ERROR] postgresql unknown type (value: `%#v`, converted: `%#v`)\n", value, converted)
				luaRow.RawSetInt(i+1, lua.LNil) // на самом деле ничего не значит
			}
		}
		luaRows.RawSet(lua.LNumber(rowCount), luaRow)
		rowCount++
	}
	luaColumnCount = lua.LNumber(len(cols) + 1)
	luaRowCount = lua.LNumber(rowCount)
	return
}

// создание коннекта
func (c *dslConfig) dslNewPgsqlConn(L *lua.LState) int {

	setStringValue := func(c *pgsqlConn, t *lua.LTable, key string) {
		luaVal := t.RawGetString(key)
		if val, ok := luaVal.(lua.LString); ok {
			switch key {
			case "host":
				c.host = string(val)
			case "user", "username":
				c.user = string(val)
			case "passwd", "password":
				c.passwd = string(val)
			case "db", "database":
				c.database = string(val)
			default:
				L.RaiseError("unknown option key: %s", key)
			}
		}
	}

	conn := &pgsqlConn{
		host:     "127.0.0.1",
		database: "postgres",
		user:     "postgres",
		passwd:   "",
		port:     5432,
	}
	tbl := L.CheckTable(1)
	setStringValue(conn, tbl, "host")
	setStringValue(conn, tbl, "database")
	setStringValue(conn, tbl, "user")
	setStringValue(conn, tbl, "passwd")
	luaPort := tbl.RawGetString("port")
	if port, ok := luaPort.(lua.LNumber); ok {
		conn.port = int(port)
	}

	if err := conn.connect(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = conn
	L.SetMetatable(ud, L.GetTypeMetatable("postgres"))
	L.Push(ud)
	log.Printf("[INFO] New postgres connection `%s`\n", conn.connectionString())
	return 1
}

// выполнение запроса
func (c *dslConfig) dslPgsqlQuery(L *lua.LState) int {
	conn := checkPgsqlConn(L)
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
func (c *dslConfig) dslPgsqlClose(L *lua.LState) int {
	conn := checkPgsqlConn(L)
	if conn.db != nil {
		conn.db.Close()
	}
	return 0
}
