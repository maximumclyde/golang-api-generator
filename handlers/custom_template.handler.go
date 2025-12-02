package handlers

import (
	"github.com/maximumclyde/golang-api-generator/store"

	"github.com/gin-gonic/gin"
)

// #region types
type CustomTemplateHandler struct {
	Store      *store.Store
	RouteGroup string
}

func NewCustomTemplateHandler(s *store.Store) *CustomTemplateHandler {
	//#region new handler
	handler := &CustomTemplateHandler{
		Store:      s,
		RouteGroup: "templates",
	}
	return handler
}

func (h *CustomTemplateHandler) RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	//#region register routes
	// pb := public.Group("/" + h.RouteGroup)
	// pr := protected.Group("/" + h.RouteGroup)
}
