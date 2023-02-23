package domains

import "github.com/gin-gonic/gin"

type ControllerCrud interface {
	Get(ctx *gin.Context)
}
