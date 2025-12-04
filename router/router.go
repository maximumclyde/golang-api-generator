package router

import (
	"api-generator/handlers"
	"api-generator/store"

	"github.com/gin-gonic/gin"
)

func NewRouter(s *store.Store) *gin.Engine {
	//#region new
	g := gin.New()

	publicRoutes := g.Group("")
	protectedRoutes := g.Group("")

	//#region Template
	templateHandler := handlers.NewTemplateHandler(s)
	templateHandler.RegisterRoutes(publicRoutes, protectedRoutes)

	return g
}
