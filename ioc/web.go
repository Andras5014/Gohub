package ioc

import (
	"context"
	"fmt"
	"github.com/Andras5014/webook/internal/web/handler/article"
	"github.com/Andras5014/webook/internal/web/handler/oauth2"
	"github.com/Andras5014/webook/internal/web/handler/user"
	ijwt "github.com/Andras5014/webook/internal/web/jwt"
	"github.com/Andras5014/webook/internal/web/middleware"
	"github.com/Andras5014/webook/pkg/ginx/middlewares/logger"
	"github.com/Andras5014/webook/pkg/ginx/middlewares/metric"
	"github.com/Andras5014/webook/pkg/ginx/middlewares/ratelimit"
	zapLogger "github.com/Andras5014/webook/pkg/logx"
	ratelimit2 "github.com/Andras5014/webook/pkg/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *user.Handler, oauth2Hdl *oauth2.WeChatHandler, articleHdl *article.Handler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	oauth2Hdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	return server

}

func InitLimiter(redisClient redis.Cmdable) ratelimit2.Limiter {
	return ratelimit2.NewRedisSlideWindowLimiter(redisClient, time.Second, 10)
}
func InitMiddlewares(limiter ratelimit2.Limiter, jwtHdl ijwt.Handler, l zapLogger.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		metric.NewBuilder("andras", "webook", "gin_http", "统计 gin 的 http 接口").Build(),
		logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("HTTP Request", zapLogger.Field{
				Key:   "al",
				Value: al,
			})
		}).AllowReqBody().AllowRespBody().Build(),
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
