package main

import (
	"api-generator/store"

	"github.com/gin-gonic/gin"
)

func NewRouter(s *store.Store) *gin.Engine {
	//#region new
	g := gin.New()

	return g
}
