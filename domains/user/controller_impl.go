package user

import (
	"fmt"
	"net/http"
	"os"
	"pvg/simada/service-epersediaan/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/golang-jwt/jwt/v5"
)

type UserControllerImpl struct{}

func NewUserController() UserController {
	return &UserControllerImpl{}
}

// Token   godoc
// @Summary      Users
// @Description  get users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  UserModel
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Security 	 ApiKeyAuth
// @Router       /user/data/get [get]
func (u *UserControllerImpl) Get(ctx *gin.Context) {
	// not impelement yet
}

// Token   godoc
// @Summary      Token
// @Description  get a token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        message	body	RequestToken	true	"Authorization Credential"
// @Success      200  {object}  ResponseToken
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /user/token [post]
func (u *UserControllerImpl) Token(ctx *gin.Context) {
	requestToken := RequestToken{}
	err := ctx.BindJSON(&requestToken)

	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	userModel, err := NewRepository().GetByUsername(requestToken.Username)
	if err != nil {
		if err == pg.ErrNoRows {
			err = fmt.Errorf("username or password does not match")
		}
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	if !util.IsSame(requestToken.Password, userModel.Password) {
		util.NewError(ctx, http.StatusBadRequest, fmt.Errorf("username or password does not match"))
		return
	}

	today := time.Now()

	// create jwt here
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":    userModel.Username,
		"expiry_date": time.Date(today.Year(), today.Month(), today.Day(), today.Hour(), today.Minute()+int(time.Minute*30), today.Second(), 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))

	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusAccepted, ResponseToken{
		Token: tokenString,
	})
}
