package utils

import (
	"github.com/gin-gonic/gin"
)

func SendResponse(ctx *gin.Context, response HttpResponse) {
	ctx.JSON(response.Status, response)
}
