package dsl

import (
	"database/sql"

	_ "github.com/go-goracle/goracle"
)

func (o *oracleConn) connect() error {
	if o.db == nil {
		db, err := sql.Open("goracle", o.connString)
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
