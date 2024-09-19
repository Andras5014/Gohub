package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WrapBody[T any](fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, Result{Code: 400, Msg: err.Error()})
			return
		}
		res, err := fn(ctx, req)
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
