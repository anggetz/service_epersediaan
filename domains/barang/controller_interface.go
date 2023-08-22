package barang

import "github.com/gin-gonic/gin"

type Controller interface {
	GetRegisteredDataTransportation(*gin.Context)
	GetAlatAngkut(*gin.Context)
	GetNonAlatAngkut(*gin.Context)
	CheckNumberPlate(*gin.Context)
}
