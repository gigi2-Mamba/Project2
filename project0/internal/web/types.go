package web

import "github.com/gin-gonic/gin"

// Created by Changer on 2024/2/6.
// Copyright 2024 programmer.

type Handler interface {
	RegisterRoutes(server *gin.Engine)
}

type Page struct {
	Limit  int
	Offset int
}
