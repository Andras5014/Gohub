//go:build wireinject

package main

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

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,

		dao.NewUserDAO,
		article2.NewArticleDAO,

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		article.NewArticleRepository,

		service.NewCodeService,
		service.NewUserService,
		service.NewArticleService,
		ioc.InitSmsService,
		ioc.InitOAuth2WeChatService,
		ioc.InitConfig,
		ioc.InitLogger,
		ijwt.NewRedisJWTHandler,

		user.NewUserHandler,
		oauth2.NewOAuth2WeChatHandler,
		article3.NewArticleHandler,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.InitLimiter,
	)
	return new(gin.Engine)
}
