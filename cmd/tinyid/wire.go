// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	"tinyid/internal/biz"
	"tinyid/internal/conf"
	"tinyid/internal/data"
	"tinyid/internal/server"
	"tinyid/internal/service"
)

// initApp init kratos application.
func initApp(context.Context, *conf.Server, *conf.Data, log.Logger, registry.Registrar) (*kratos.App, error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
