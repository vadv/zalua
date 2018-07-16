package dsl

import (
	"database/sql"
	"log"

	lua "github.com/yuin/gopher-lua"
)

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
