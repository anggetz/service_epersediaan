package networks

import (
	"fmt"
	"net/http"
	"os"
	"pvg/simada/service-golang/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizations := c.Request.Header["Authorization"]
		if len(authorizations) == 0 {
			util.NewError(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			c.Abort()
			return
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

			util.NewError(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			if err != nil {
				util.NewError(c, http.StatusUnauthorized, err)
				c.Abort()
				return
			}

			if claims["expiry_date"].(float64) < float64(time.Now().Unix()) {
				util.NewError(c, http.StatusUnauthorized, fmt.Errorf("Token expired"))
				c.Abort()
				return
			}
		} else {
			util.NewError(c, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			c.Abort()
			return
		}

		c.Next()
	}
}
