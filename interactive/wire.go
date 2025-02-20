//go:build wireinject

package main

import (
	"github.com/Andras5014/gohub/interactive/events"
	"github.com/Andras5014/gohub/interactive/grpc"
	"github.com/Andras5014/gohub/interactive/ioc"
	"github.com/Andras5014/gohub/interactive/repository"
	"github.com/Andras5014/gohub/interactive/repository/cache"
	"github.com/Andras5014/gohub/interactive/repository/dao"
	"github.com/Andras5014/gohub/interactive/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitKafka,
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitConfig,
	ioc.InitLogger,
)
var interactiveSvcSet = wire.NewSet(
	service.NewInteractiveService,
	repository.NewInteractiveRepository,
	cache.NewInteractiveCache,
	dao.NewInteractiveDAO,
)

func InitApp() *App {
	wire.Build(
		interactiveSvcSet,
		thirdPartySet,
		grpc.NewInteractiveServiceServer,
		ioc.InitGRPCxServer,
		ioc.InitConsumers,
		events.NewInteractiveReadEventBatchConsumer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
