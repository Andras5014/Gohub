//go:build wireinject

package main

import (
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web"
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

		cache.NewUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewCodeService,
		service.NewUserService,
		ioc.InitSmsService,
		ioc.InitOAuth2WeChatService,
		ioc.InitConfig,
		ioc.InitLogger,
		ijwt.NewRedisJWTHandler,

		web.NewUserHandler,
		web.NewOAuth2WeChatHandler,
		ioc.InitWebServer,
		ioc.InitMiddlewares,
		ioc.InitLimiter,
	)
	return new(gin.Engine)
}
