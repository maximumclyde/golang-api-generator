package router

import (
	"github.com/maximumclyde/golang-api-generator/handlers"
	"github.com/maximumclyde/golang-api-generator/store"

	"github.com/gin-gonic/gin"
)

func NewRouter(s *store.Store) *gin.Engine {
	g := gin.New()

	publicRoutes := g.Group("")
	protectedRoutes := g.Group("")

	//#region Template
	templateHandler := handlers.NewTemplateHandler(s)
	templateHandler.RegisterRoutes(publicRoutes, protectedRoutes)

	return g
}
