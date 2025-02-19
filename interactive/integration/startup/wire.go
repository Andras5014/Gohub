//go:build wireinject

package startup

import (
	"github.com/Andras5014/gohub/interactive/grpc"
	"github.com/Andras5014/gohub/interactive/repository"
	"github.com/Andras5014/gohub/interactive/repository/cache"
	"github.com/Andras5014/gohub/interactive/repository/dao"
	"github.com/Andras5014/gohub/interactive/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	InitRedis,
	InitDB,
	InitConfig,
	InitLogger,
)
var interactiveSvcSet = wire.NewSet(
	service.NewInteractiveService,
	repository.NewInteractiveRepository,
	cache.NewInteractiveCache,
	dao.NewInteractiveDAO,
)

func InitInteractiveSvc() service.InteractiveService {
	wire.Build(thirdPartySet, interactiveSvcSet)
	return service.NewInteractiveService(nil)
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(thirdPartySet, interactiveSvcSet, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
