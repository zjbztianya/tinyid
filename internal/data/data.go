package data

import (
	"context"
	"github.com/google/wire"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
	"tinyid/internal/conf"
	// database driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewIdgenRepo)

// Data .
type Data struct {
	// mysql client
	db *DB
	// etcd client
	etcd     *clientv3.Client
	kv       clientv3.KV
	ttl      time.Duration
	workerID int64
}

// NewData .
func NewData(ctx context.Context, c *conf.Data, s *conf.Server) (*Data, error) {
	data := &Data{
		db:   NewMySQL(c.Database),
		etcd: NewEtcd(c.Etcd),
		ttl:  c.Etcd.TTL.AsDuration(),
	}
	data.kv = clientv3.NewKV(data.etcd)
	if err := InitEtcd(ctx, data, s.Grpc.Addr); err != nil {
		return nil, err
	}
	return data, nil
}
