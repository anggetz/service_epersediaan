package barang

import "github.com/gin-gonic/gin"

type Controller interface {
	GetRegisteredDataTransportation(*gin.Context)
	GetAlatAngkut(*gin.Context)
	CheckNumberPlate(*gin.Context)
}
