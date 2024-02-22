package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"project0/internal/events"
)

type App struct {
	server  *gin.Engine
	consumers []events.Consumer
	cron *cron.Cron
}
