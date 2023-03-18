package barang

import "github.com/gin-gonic/gin"

type Controller interface {
	GetAlatAngkut(*gin.Context)
	CheckNumberPlate(*gin.Context)
}
