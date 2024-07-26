package main

import (
	"fmt"
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/repository/cache"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/service/sms/memory"
	"github.com/Andras5014/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {

	server := InitWebServer()
	server.Run(":8080")
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	userDao := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(rdb)
	codeCache := cache.NewCodeCache(rdb)

	userRepo := repository.NewUserRepository(userDao, userCache)
	codeRepo := repository.NewCodeRepository(codeCache)
	userSvc := service.NewUserService(userRepo)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	userHandler := web.NewUserHandler(userSvc, codeSvc)
	return userHandler
}
func initWebServer() *gin.Engine {
	server := gin.Default()
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: "127.0.0.1:16379",
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			fmt.Println("origin", origin)
			//if strings.HasPrefix(origin, "http://127.0.0.1") {
			//	return true
			//}
			//return strings.Contains(origin, "andras.icu")
			return true
		},
		MaxAge:        12 * time.Hour,
		ExposeHeaders: []string{"x-jwt-token"},
	}))

	//session
	//store, err := redis.NewStore(16, "tcp", "127.0.0.1:16379", "", []byte("secret"), []byte("secret"))
	//if err != nil {
	//	panic(err)
	//}
	//server.Use(sessions.Sessions("mysession", store))
	return server
}
