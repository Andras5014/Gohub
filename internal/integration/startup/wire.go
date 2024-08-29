//go:build wireinject

package startup

import (
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/repository/article"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web"
	ijwt "github.com/Andras5014/webook/internal/web/jwt"
	"github.com/Andras5014/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	InitConfig,

	InitLogger,
	InitDB,
	InitRedis,
	InitSmsService,
)

var codeSvcProvider = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	service.NewCodeService,
)

var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewUserRepository,
	service.NewUserService,
)

var articleSvcProvider = wire.NewSet(
	dao.NewArticleDAO,
	article.NewArticleRepository,
	service.NewArticleService,
)

var oauth2SvcProvider = wire.NewSet(
	InitOAuth2WeChatService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		codeSvcProvider,
		articleSvcProvider,
		oauth2SvcProvider,

		// handler
		ioc.InitMiddlewares,
		ioc.InitLimiter,

		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WeChatHandler,

		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		articleSvcProvider,
		web.NewArticleHandler,
	)
	return new(web.ArticleHandler)
}
