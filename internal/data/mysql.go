package data

import (
	"database/sql"
	"github.com/pkg/errors"
	"tinyid/internal/conf"

	// database driver
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db   *sql.DB
	conf *conf.Data_Database
}

func NewMySQL(c *conf.Data_Database) *DB {
	if c.QueryTimeout.AsDuration() == 0 ||
		c.ExecTimeout.AsDuration() == 0 ||
		c.TranTimeout.AsDuration() == 0 {
		panic("mysql must be set query/execute/transaction timeout")
	}

	db, err := sql.Open(c.Driver, c.Source)
	if err != nil {
		panic(errors.Wrap(err, "open mysql error"))
	}
	db.SetMaxIdleConns(int(c.Idle))
	db.SetMaxOpenConns(int(c.Active))
	db.SetConnMaxIdleTime(c.IdleTimeout.AsDuration())

	return &DB{
		db:   db,
		conf: c,
	}
}
