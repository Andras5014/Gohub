package ginx

import (
	"fmt"
	"github.com/Andras5014/webook/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

func WrapBody[T any](l logx.Logger, fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusOK, InvalidParam())
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			l.WithCtx(ctx).Error("handle http error: ",
				logx.Error(err),
				logx.Any("path", ctx.Request.URL.Path),
				logx.Any("route", fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())),
			)
		}
		ctx.JSON(http.StatusOK, res)
	}
}
func Wrap(l logx.Logger, fn func(ctx *gin.Context) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := fn(ctx)
		if err != nil {
			l.WithCtx(ctx).Error("handle http error: ",
				logx.Error(err),
				logx.Any("path", ctx.Request.URL.Path),
				logx.Any("route", fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())),
			)
		}
		ctx.JSON(http.StatusOK, res)
	}
}

var vector prometheus.CounterVec

func InitCounter(opts prometheus.CounterOpts) {
	vector = *prometheus.NewCounterVec(opts, []string{"code"})
	prometheus.MustRegister(&vector)
}
