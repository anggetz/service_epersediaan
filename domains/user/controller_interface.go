package user

import (
	"pvg/simada/service-epersediaan/domains"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	domains.ControllerCrud
	Token(ctx *gin.Context)
}
