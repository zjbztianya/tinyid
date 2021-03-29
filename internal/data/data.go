package data

import (
	"database/sql"
	"github.com/google/wire"
	clientv3 "go.etcd.io/etcd/client/v3"
	"tinyid/internal/conf"
	// database driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewIdgenRepo)

// Data .
type Data struct {
	// mysql client
	db *sql.DB
	// etcd client
	etcd *clientv3.Client
}

// NewData .
func NewData(c *conf.Data) (*Data, error) {

	return &Data{}, nil
}
