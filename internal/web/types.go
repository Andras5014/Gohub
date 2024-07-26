package web

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterRouters(engine *gin.Engine)
}
