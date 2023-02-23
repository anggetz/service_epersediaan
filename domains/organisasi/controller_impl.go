package organisasi

import (
	"github.com/gin-gonic/gin"
)

type OrganisasiControllerImpl struct{}

func NewOrganisasiController() OrganisasiController {
	return &OrganisasiControllerImpl{}
}

func (u *OrganisasiControllerImpl) Get(ctx *gin.Context) {
	// not impelement yet

}
