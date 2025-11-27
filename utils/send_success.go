package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSuccess(ctx *gin.Context, data any) {
	SendResponse(ctx, HttpResponse{
		IsSuccessful: true,
		Status:       http.StatusOK,
		Code:         "",
		Data:         data,
	})
}
