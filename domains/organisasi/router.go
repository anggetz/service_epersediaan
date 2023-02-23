package organisasi

import (
	"pvg/simada/service-epersediaan/domains"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func NewUserRouter() domains.Route {
	return &UserRouter{}
}

func (u *UserRouter) RegisterHandler(r *gin.RouterGroup) {
}
