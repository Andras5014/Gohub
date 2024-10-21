package main

import (
	"github.com/Andras5014/webook/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	Server    *gin.Engine
	Consumers []events.Consumer
}
