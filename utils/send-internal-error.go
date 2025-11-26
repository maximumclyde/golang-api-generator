package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendInternalError(ctx *gin.Context, err error) {
	SendResponse(ctx, HttpResponse{
		IsSuccessful: false,
		Status:       http.StatusInternalServerError,
		Code:         "internal_server_error",
		Data: ErrorData{
			Error: err.Error(),
		},
	})
}
