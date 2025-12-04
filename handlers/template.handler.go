package handlers

import (
	"errors"

	"github.com/maximumclyde/golang-api-generator/models"
	"github.com/maximumclyde/golang-api-generator/store"
	"github.com/maximumclyde/golang-api-generator/utils"

	"github.com/gin-gonic/gin"
)

// #region types
type TemplateHandler struct {
	Store      *store.Store
	RouteGroup string
}

func NewTemplateHandler[T models.Template](s *store.Store) *TemplateHandler {
	//#region new handler
	handler := &TemplateHandler{
		Store:      s,
		RouteGroup: "templates",
	}
	return handler
}

func (h *TemplateHandler) Create(ctx *gin.Context) {
	//#region create
	var body = &models.Template{}
	err := ctx.ShouldBindJSON(body)

	if err != nil {
		utils.SendBadRequest(ctx, err)
		return
	}

	err = h.Store.Services.Template.Create(ctx, body)
	if err != nil {
		utils.SendInternalError(ctx, err)
		return
	}

	utils.SendSuccess(ctx, nil)
}

func (h *TemplateHandler) GetById(ctx *gin.Context) {
	//#region get id
	id, exists := ctx.Params.Get("id")
	if !exists {
		utils.SendBadRequest(ctx, errors.New("param id is missing"))
		return
	}

	data, err := h.Store.Services.Template.GetById(ctx, id)
	if err != nil {
		utils.SendInternalError(ctx, err)
		return
	}

	utils.SendSuccess(ctx, data)
}

func (h *TemplateHandler) Find(ctx *gin.Context) {
	//#region find
	query := &models.TemplateQuery{}
	err := ctx.ShouldBind(query)
	if err != nil {
		utils.SendBadRequest(ctx, err)
		return
	}

	data, err := h.Store.Services.Template.Find(ctx)
	if err != nil {
		utils.SendInternalError(ctx, err)
		return
	}

	utils.SendSuccess(ctx, data)
}

func (h *TemplateHandler) Patch(ctx *gin.Context) {
	//#region patch
	id, exists := ctx.Params.Get("id")
	if !exists {
		utils.SendBadRequest(ctx, errors.New("param id is missing"))
		return
	}

	data := &map[string]any{}
	err := ctx.ShouldBindJSON(data)
	if err != nil {
		utils.SendBadRequest(ctx, err)
		return
	}

	err = h.Store.Services.Template.Patch(ctx, id, data)
	if err != nil {
		utils.SendInternalError(ctx, err)
		return
	}

	utils.SendSuccess(ctx, nil)
}

func (h *TemplateHandler) Remove(ctx *gin.Context) {
	//#region remove
	id, exists := ctx.Params.Get("id")
	if !exists {
		utils.SendBadRequest(ctx, errors.New("param id is missing"))
		return
	}

	err := h.Store.Services.Template.Remove(ctx, id)
	if err != nil {
		utils.SendInternalError(ctx, err)
		return
	}

	utils.SendSuccess(ctx, nil)
}

func (h *TemplateHandler) RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup) {
	//#region register routes
	pb := public.Group("/" + h.RouteGroup)
	{
		pb.GET("/", h.Find)
		pb.GET("/:id", h.GetById)
	}

	pr := protected.Group("/" + h.RouteGroup)
	{
		pr.POST("/", h.Create)
		pr.PATCH("/:id", h.Patch)
		pr.DELETE("/:id", h.Remove)
	}
}
