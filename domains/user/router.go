package user

import (
	"pvg/simada/service-golang/domains"
	"pvg/simada/service-golang/networks"

	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func NewUserRouter() domains.Route {
	return &UserRouter{}
}

func (u *UserRouter) RegisterHandler(r *gin.RouterGroup) {

	authApi := r.Group("auth")
	{
		authApi.POST("/token", NewUserController().Token)
	}

	dataApi := r.Group("data")
	{
		dataApi.Use(networks.AuthJWTMiddleware())
		dataApi.GET("/get", NewUserController().Get)
		dataApi.POST("/iam", NewUserController().IAM)
	}

}
