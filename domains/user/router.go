package user

import (
	"pvg/simada/service-epersediaan/domains"
	"pvg/simada/service-epersediaan/networks"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func NewUserRouter() domains.Route {
	return &UserRouter{}
}

func (u *UserRouter) RegisterHandler(r *gin.RouterGroup) {

	r.POST("/token", NewUserController().Token)

	dataApi := r.Group("data").Use(networks.AuthJWTMiddleware())
	{
		dataApi.GET("/get", NewUserController().Get)
	}

}
