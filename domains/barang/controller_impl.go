package barang

import (
	"math"
	"net/http"
	"pvg/simada/service-golang/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ControllerImpl struct{}

func NewController() Controller {
	return &ControllerImpl{}
}

// Token   godoc
// @Summary      Barang
// @Description  get alat angkut master
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        take	query	int			false	"take"
// @Param        page	query	int			false	"page"
// @Param        search	query	string		false	"search"
// @Success      200  {array}  	Model
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Security 	 ApiKeyAuth
// @Router       /barang/data/get-alat-angkut [get]
func (c *ControllerImpl) GetAlatAngkut(ctx *gin.Context) {
	params := ParamPagination{}

	params.take, _ = strconv.Atoi(ctx.Request.URL.Query().Get("take"))
	params.page, _ = strconv.Atoi(ctx.Request.URL.Query().Get("page"))
	params.search = ctx.Query("search")

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

	databarang, totalData, err := NewUseCase().GetApelMaster(params.take, offset, params.search)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":      databarang,
		"dataTotal": totalData,
		"totalPage": int(math.Ceil(float64(totalData) / float64(params.take))),
		"page":      params.page,
	})
}

// Token   godoc
// @Summary      Barang
// @Description  get alat angkut master
// @Tags         barang
// @Accept       json
// @Produce      json
// @Param        number_plate	query	string		true	"number_plate"
// @Param        chassis_number	query	string		false	"chassis_number"
// @Param        machine_number	query	string		false	"machine_number"
// @Param        pidopd			query	int			false	"pidopd"
// @Param        sub_pidopd		query	int			false	"sub_pidopd"
// @Param        pidupt			query	int			false	"pidupt"
// @Success      200  {object}  MesinModel
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Security 	 ApiKeyAuth
// @Router       /barang/data/check-number-plate [get]
func (c *ControllerImpl) CheckNumberPlate(ctx *gin.Context) {
	params := ParamCheckNumberPlate{}

	errHttpRequired := util.IsRequiredKeyAvail([]string{"number_plate"}, ctx.Request.URL.Query())
	if errHttpRequired != nil {
		util.NewError(ctx, http.StatusBadRequest, errHttpRequired)
		return
	}

	params.NumberPlate = ctx.Request.URL.Query().Get("number_plate")
	params.ChassisNumber = ctx.Request.URL.Query().Get("chassis_number")
	params.MachineNumber = ctx.Request.URL.Query().Get("machine_number")
	params.Pidopd, _ = strconv.Atoi(ctx.Request.URL.Query().Get("pidopd"))
	params.SubPidopd, _ = strconv.Atoi(ctx.Request.URL.Query().Get("sub_pidopd"))
	params.Pidupt, _ = strconv.Atoi(ctx.Request.URL.Query().Get("pidupt"))

	databarang, err := NewUseCase().CheckPlatNumberChassisNumberAndMachineNumber(params.NumberPlate, params.ChassisNumber, params.MachineNumber, params.Pidopd, params.SubPidopd, params.Pidupt)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": databarang,
	})
}
