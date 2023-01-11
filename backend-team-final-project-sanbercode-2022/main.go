package main

import (
	"fintech/configs"
	"fintech/docs"
	"fintech/routes"
	"fintech/utils"
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/
func main() {
	//programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Modal.in API"
	docs.SwaggerInfo.Description = "This is REST API Documentation for Modal.in"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = utils.Getenv("SWAGGER_HOST", "localhost:8080")
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// connect with database
	db := configs.ConnectDataBase()
	sqlDB, _ := db.DB()

	defer sqlDB.Close()

	// running router
	r := routes.SetupRouter(db)
	r.Run()
}
