//go:build wireinject

package main

import (
	"github.com/Andras5014/webook/interactive/events"
	repository2 "github.com/Andras5014/webook/interactive/repository"
	cache2 "github.com/Andras5014/webook/interactive/repository/cache"
	dao2 "github.com/Andras5014/webook/interactive/repository/dao"
	service2 "github.com/Andras5014/webook/interactive/service"
	articleEvent "github.com/Andras5014/webook/internal/events/article"
	"github.com/Andras5014/webook/internal/repository"
	articleRepo "github.com/Andras5014/webook/internal/repository/article"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	articleDao "github.com/Andras5014/webook/internal/repository/dao/article"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web/handler/article"
	"github.com/Andras5014/webook/internal/web/handler/oauth2"
	"github.com/Andras5014/webook/internal/web/handler/user"
	ijwt "github.com/Andras5014/webook/internal/web/jwt"
	"github.com/Andras5014/webook/ioc"
	"github.com/google/wire"
)

var rankingSvcSet = wire.NewSet(
	cache.NewRedisRankingCache,
	repository.NewRankingRepository,
	service.NewRankingService,
)
var interactiveSvcSet = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewInteractiveRepository,
	cache2.NewInteractiveCache,
	dao2.NewInteractiveDAO,
)

var userSvcSet = wire.NewSet(
	service.NewUserService,
	repository.NewUserRepository,
	cache.NewUserCache,
	dao.NewUserDAO,
)

var articleSvcSet = wire.NewSet(
	service.NewArticleService,
	articleRepo.NewArticleRepository,
	articleDao.NewArticleDAO,
	cache.NewRedisArticleCache,
)
var codeSvcProvider = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	service.NewCodeService,
)

var thirdPartySet = wire.NewSet(
	ioc.InitConfig,

	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitRedisUniversalClient,
	ioc.InitRedSync,
	ioc.InitSmsService,
	ioc.InitKafka,
	ioc.InitSyncProducer,
	ioc.InitConsumers,
)

func InitApp() *App {
	wire.Build(
		//event
		articleEvent.NewSaramaSyncProducer,
		events.NewInteractiveReadEventBatchConsumer,

		user.NewUserHandler,
		userSvcSet,

		codeSvcProvider,

		article.NewArticleHandler,
		articleSvcSet,
		interactiveSvcSet,
		ioc.InitInteractiveGrpcClient,
		thirdPartySet,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.InitLimiter,

		// job
		rankingSvcSet,
		ioc.InitRankingJob,
		ioc.InitJobs,

		oauth2.NewOAuth2WeChatHandler,
		ioc.InitOAuth2WeChatService,
		ijwt.NewRedisJWTHandler,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
