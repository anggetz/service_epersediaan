package main

import (
	"flag"
	"fmt"
	"os"
	"pvg/simada/service-epersediaan/docs"
	"pvg/simada/service-epersediaan/domains/organisasi"
	"pvg/simada/service-epersediaan/domains/user"

	// gin-swagger middleware

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swagger embed files

// @title           E-Persediaan API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8089
// @BasePath  /v1/public-api

// @securitydefinitions.apikey  ApiKeyAuth
// @in 		header
// @name 	Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {

	var configPath = flag.String("config", ".env.local", "please input config path")

	flag.Parse()

	godotenv.Load(*configPath)

	fmt.Println("Configured", *configPath)

	route := gin.New()

	docs.SwaggerInfo.Host = os.Getenv("EPERSEDIAAN_SWAGGER_HOST")

	v1 := route.Group("v1")
	{
		publicApi := v1.Group("public-api")
		{
			userApi := publicApi.Group("user")
			{
				user.NewUserRouter().RegisterHandler(userApi)
			}

			organisasiAPI := publicApi.Group("organisasi")
			{
				organisasi.NewUserRouter().RegisterHandler(organisasiAPI)
			}

			docApi := publicApi.Group("doc")
			{
				docApi.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			}

		}

	}

	fmt.Println("Server listening on port " + os.Getenv("EPERSEDIAAN_APP_PORT"))

	route.Run(":" + os.Getenv("EPERSEDIAAN_APP_PORT"))
}
