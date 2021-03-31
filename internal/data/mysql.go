package data

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"tinyid/internal/conf"

	// database driver
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(c *conf.Data_Database) *sqlx.DB {
	if c.QueryTimeout.AsDuration() == 0 || c.ExecTimeout.AsDuration() == 0 {
		panic("mysql must be set query/execute timeout")
	}

	db, err := sqlx.Connect(c.Driver, c.Source)
	if err != nil {
		panic(errors.Wrap(err, "open mysql error"))
	}

	db.SetMaxIdleConns(int(c.Idle))
	db.SetMaxOpenConns(int(c.Active))
	db.SetConnMaxIdleTime(c.IdleTimeout.AsDuration())

	return db
}
