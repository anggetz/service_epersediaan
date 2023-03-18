package barang

import (
	"pvg/simada/service-golang/domains"
	"pvg/simada/service-golang/networks"

	"github.com/gin-gonic/gin"
)

type Router struct{}

func NewRouter() domains.Route {
	return &Router{}
}

func (u *Router) RegisterHandler(r *gin.RouterGroup) {
	dataApi := r.Group("data").Use(networks.AuthJWTMiddleware())
	{
		dataApi.GET("/get-alat-angkut", NewController().GetAlatAngkut)
		dataApi.GET("/check-number-plate", NewController().CheckNumberPlate)
	}
}
