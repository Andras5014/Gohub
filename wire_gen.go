// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web"
	"github.com/Andras5014/webook/internal/web/jwt"
	"github.com/Andras5014/webook/ioc"
	"github.com/gin-gonic/gin"
)

import (
	_ "gorm.io/driver/mysql"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := ioc.InitRedis()
	limiter := ioc.InitLimiter(cmdable)
	handler := jwt.NewRedisJWTHandler(cmdable)
	v := ioc.InitMiddlewares(limiter, handler)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSmsService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, handler)
	oauth2Service := ioc.InitOAuth2WeChatService()
	oAuth2WeChatHandler := web.NewOAuth2WeChatHandler(oauth2Service, userService, handler)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WeChatHandler)
	return engine
}
