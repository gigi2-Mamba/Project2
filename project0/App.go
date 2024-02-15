package main

import (
	"github.com/gin-gonic/gin"
	"project0/internal/events"
)

type App struct {
	server  *gin.Engine
	consumers []events.Consumer
}
