package ioc

import (
	"fmt"
	"github.com/Andras5014/webook/internal/web"
	"github.com/Andras5014/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRouters(server)
	return server

}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cors.New(cors.Config{
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
		})
	}
}
