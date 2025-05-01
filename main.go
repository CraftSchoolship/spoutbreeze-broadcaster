package main

import (
    "log"
	"os"
	"spoutbreeze/initializers"
    "spoutbreeze/routes"

	_ "spoutbreeze/docs" // This is for Swagger documentation

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
}

//  @title SpoutBreeze API
//  @version 1.0
//  @description This is a sample server for SpoutBreeze.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/
//	@Schemes	http

func main() {
	router := routes.SetupRouter()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1323" // Default port
	}
	err := router.Run(":" + PORT)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}