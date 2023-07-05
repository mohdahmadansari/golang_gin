package main

import (
	"github.com/rs/zerolog/log"

	_ "github.com/mohdahmadansari/golang_gin/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohdahmadansari/golang_gin/database"
	"github.com/mohdahmadansari/golang_gin/router"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// "gorm.io/gorm"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	// println(db)
	// println(err)
}

// @title 	Go lang API
// @description 	This API developed using gin framework by Ahmad
// @host 	localhost:5000
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	var notFoundMessage = "API not found"
	r := gin.Default()

	// r.Static("/swagger-ui/", "./dist/swagger-ui")

	db, dbError := database.CreateConnection()

	if dbError != nil {
		notFoundMessage = "Database connection issue."
	} else {
		database.SetupMigration(db)
		r = router.AllRouter(r, db)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"success": "0", "message": notFoundMessage})
	})

	r.Run(":5000")

}
