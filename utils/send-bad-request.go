package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendBadRequest(ctx *gin.Context, err error) {
	SendResponse(ctx, HttpResponse{
		IsSuccessful: false,
		Status:       http.StatusBadRequest,
		Code:         "bad_request",
		Data: ErrorData{
			Error: err.Error(),
		},
	})
}
