package main

import (
	"github.com/Andras5014/gohub/internal/events"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type App struct {
	Server    *gin.Engine
	Consumers []events.Consumer
	Cron      *cron.Cron
}
