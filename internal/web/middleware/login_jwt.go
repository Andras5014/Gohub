package middleware

import (
	ijwt "github.com/Andras5014/webook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}

}
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" ||
			path == "/users/refresh_token" {
			// 不需要登录校验
			return
		}
		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		//如果redis崩溃不至于全部用户都过不校验 可以选择直接return
		// if redis 崩溃 return
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// redis问题 或者退出登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//续约或者长短token
		//claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
		//tokenStr, err = token.SignedString([]byte("secret"))
		//if err != nil {
		//	log.Println("续约失败")
		//}
		//ctx.Header("x-jwt-token", tokenStr)
		ctx.Set("claims", claims)
		ctx.Set("userId", claims.Uid)
	}

}
