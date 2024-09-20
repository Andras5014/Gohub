//go:build wireinject

package startup

import (
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/repository/article"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	article2 "github.com/Andras5014/webook/internal/repository/dao/article"
	"github.com/Andras5014/webook/internal/service"
	article3 "github.com/Andras5014/webook/internal/web/handler/article"
	"github.com/Andras5014/webook/internal/web/handler/oauth2"
	"github.com/Andras5014/webook/internal/web/handler/user"
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
	article2.NewArticleDAO,
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

		user.NewUserHandler,
		article3.NewArticleHandler,
		oauth2.NewOAuth2WeChatHandler,

		ijwt.NewRedisJWTHandler,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

func InitArticleHandler() *article3.Handler {
	wire.Build(
		thirdPartySet,
		articleSvcProvider,
		article3.NewArticleHandler,
	)
	return new(article3.Handler)
}
func InitArticleHandlerV1(dao article2.ArticleDAO) *article3.Handler {
	wire.Build(

		InitLogger,

		article.NewArticleRepository,
		service.NewArticleService,
		article3.NewArticleHandler,
	)
	return new(article3.Handler)
}
