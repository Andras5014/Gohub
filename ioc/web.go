package ioc

import (
	"fmt"
	"github.com/Andras5014/webook/internal/web"
	ijwt "github.com/Andras5014/webook/internal/web/jwt"
	"github.com/Andras5014/webook/internal/web/middleware"
	"github.com/Andras5014/webook/pkg/ginx/middlewares/ratelimit"
	ratelimit2 "github.com/Andras5014/webook/pkg/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler, oauth2 *web.OAuth2WeChatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	oauth2.RegisterRoutes(server)
	return server

}

func InitLimiter(redisClient redis.Cmdable) ratelimit2.Limiter {
	return ratelimit2.NewRedisSlideWindowLimiter(redisClient, time.Second, 10)
}
func InitMiddlewares(limiter ratelimit2.Limiter, jwtHdl ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		ratelimit.NewBuilder(limiter).Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).Build(),
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
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		})
	}
}
