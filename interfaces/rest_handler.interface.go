package interfaces

import "github.com/gin-gonic/gin"

type IRestHandler interface {
	RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup)
}
