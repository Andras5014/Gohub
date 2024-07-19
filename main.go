package main

import (
	"github.com/Andras5014/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	uh := web.UserHandler{}
	uh.RegisterRouters(engine)
	engine.Run(":8080")
}
