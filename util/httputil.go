package util

import (
	"fmt"
	"net/url"
	"os"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

// NewError example
func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

// check required query string

func IsRequiredKeyAvail(keys []string, urls url.Values) error {
	for _, key := range keys {
		if !urls.Has(key) {
			return fmt.Errorf("%s is required", key)
		}
	}

	return nil
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

func GetUsername(ctx *gin.Context) (string, error) {
	authorizations := ctx.Request.Header["Authorization"]
	if len(authorizations) == 0 {
		return "", fmt.Errorf("header authorization empty")
	}

	token, err := jwt.Parse(authorizations[0], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err != nil {
			return "", err
		}

		return claims["username"].(string), nil

	} else {

		return "", fmt.Errorf("cannot claims jwt token")
	}
}
