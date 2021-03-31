package main

import (
	"context"
	"flag"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"

	etcd "github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gopkg.in/yaml.v2"
	"tinyid/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

func newApp(ctx context.Context, logger log.Logger, hs *http.Server, gs *grpc.Server, rg registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rg),
		kratos.Context(ctx),
	)
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	//service register
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   bc.Data.Etcd.Endpoints,
		DialTimeout: bc.Data.Etcd.DialTimeout.AsDuration()})
	if err != nil {
		panic(err)
	}
	reg := etcd.New(cli,
		etcd.Namespace("idgen"),
		etcd.RegisterTTL(bc.Data.Etcd.TTL.AsDuration()),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := initApp(ctx, bc.Server, bc.Data, logger, reg)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
