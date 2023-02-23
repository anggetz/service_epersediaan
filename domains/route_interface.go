package domains

import "github.com/gin-gonic/gin"

type Route interface {
	RegisterHandler(*gin.RouterGroup)
}
