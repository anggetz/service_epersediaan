package main

import (
	"flag"
	"fmt"
	"os"
	"pvg/simada/service-golang/docs"
	"pvg/simada/service-golang/domains/barang"
	"pvg/simada/service-golang/domains/organisasi"
	"pvg/simada/service-golang/domains/penyusutan"
	"pvg/simada/service-golang/domains/user"
	"time"

	"github.com/nats-io/nats.go"

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

	var configPath = flag.String("config", ".env", "please input config path")

	flag.Parse()

	godotenv.Load(*configPath)

	// connect to nats
	nc, err := nats.Connect(fmt.Sprintf("%s:%s", os.Getenv("NATS_HOST"), os.Getenv("NATS_PORT")))
	if err != nil {
		panic(err)
	}

	// register the subcription
	penyusutan.NewChannel(nc).RegisterCalcPenyusutan()

	defer nc.Drain()

	defer nc.Close()

	// fmt.Println(string(msg.Data))

	fmt.Println("Configured", *configPath)

	route := gin.New()

	docs.SwaggerInfo.Host = os.Getenv("EPERSEDIAAN_SWAGGER_HOST")

	v1 := route.Group("v1")
	{
		publicApi := v1.Group("public-api")
		{
			barangApi := publicApi.Group("barang")
			{
				barang.NewRouter().RegisterHandler(barangApi)
			}

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

	penyusutan.NewUseCase().CalcPenyusutan(186, "3", time.Date(2013, 01, 10, 0, 0, 0, 0, time.Local))

	fmt.Println("Server listening on port " + os.Getenv("EPERSEDIAAN_APP_PORT"))

	route.Run(":" + os.Getenv("EPERSEDIAAN_APP_PORT"))
}
