package user

import (
	"pvg/simada/service-golang/domains"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	domains.ControllerCrud
	Token(ctx *gin.Context)
	IAM(ctx *gin.Context)
}
