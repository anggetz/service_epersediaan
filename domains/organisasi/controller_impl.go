package organisasi

import (
	"fmt"
	"math"
	"net/http"
	"pvg/simada/service-epersediaan/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

type OrganisasiControllerImpl struct{}

func NewOrganisasiController() OrganisasiController {
	return &OrganisasiControllerImpl{}
}

// Token   		 godoc
// @Summary      Organisasi
// @Description  get m_organisasi
// @Tags         organisasi
// @Accept       json
// @Produce      json
// @Param        take	query	int			false	"take"
// @Param        page	query	int			false	"page"
// @Param        search	query	string		false	"search"
// @Success      200  {object}  OrganisasiModel
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Security 	 ApiKeyAuth
// @Router       /organisasi/data/get [get]
func (u *OrganisasiControllerImpl) Get(ctx *gin.Context) {
	params := ParamPagination{}

	params.take, _ = strconv.Atoi(ctx.Request.URL.Query().Get("take"))
	params.page, _ = strconv.Atoi(ctx.Request.URL.Query().Get("page"))
	params.search = ctx.Query("search")

	print(params.take)

	if params.take == 0 {
		params.take = 15
	}

	if params.page == 0 {
		params.page = 1
	}

	var offset int = 0
	if params.page == 1 {
		offset = 0
	} else {
		offset = (params.page - 1) * params.take
	}

	organisasi, err := NewRepository().GetAllOrganisasi(params.page, params.take, offset, params.search)
	if err != nil {
		if err == pg.ErrNoRows {
			err = fmt.Errorf("error")
		}
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	var ln int = len(organisasi)

	ctx.JSON(http.StatusOK, gin.H{
		"data":      organisasi,
		"dataTotal": ln,
		"totalPage": int(math.Ceil(float64(ln) / float64(params.take))),
		"page":      params.page,
	})
}
