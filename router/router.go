package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mohdahmadansari/golang_gin/controllers"
	"github.com/mohdahmadansari/golang_gin/database"
	"github.com/mohdahmadansari/golang_gin/middlewares"
)

func AllRouter(r *gin.Engine) *gin.Engine {

	var notFoundMessage = "invalid route"
	connStr, _ := database.GetConnectionString()
	db, dbError := database.CreateConnection(connStr)

	if dbError != nil {
		notFoundMessage = "Database connection failed."
	} else {
		database.SetupMigration(db)
		database.SeedDatabase(db)
	}

	ctr := controllers.NewController(db)

	api := r.Group("/api")
	{
		api.GET("/", ctr.Welcome)

		api.POST("/login", controllers.NewAdminCtr(db).Login)

		adminApi := api.Group("/admin")
		adminApi.Use(middlewares.Authenticate("admin", db))
		{
			adminApi.GET("/", controllers.NewAdminCtr(db).Dashboard)
			adminApi.GET("/dashboard", controllers.NewAdminCtr(db).Dashboard)

			adminNursesApi := adminApi.Group("/nurse")
			{
				adminNursesApi.GET("", controllers.NewNurseCtrl(db).Get)
				adminNursesApi.POST("", controllers.NewNurseCtrl(db).Post)
				adminNursesApi.GET(":id", controllers.NewNurseCtrl(db).GetOne)
				adminNursesApi.PUT(":id", controllers.NewNurseCtrl(db).Put)
				adminNursesApi.DELETE(":id", controllers.NewNurseCtrl(db).Delete)

				adminNursesApi.GET("/own", controllers.NewNurseCtrl(db).Getown)
			}

		}

		nurseApi := api.Group("/nurse")
		nurseApi.POST("/login", controllers.NewNurseCtrl(db).Login)
		nurseApi.Use(middlewares.Authenticate("nurse", db))
		{
			nurseApi.GET("/profile", controllers.NewNurseCtrl(db).GetProfile)
			nurseApi.POST("/profile", controllers.NewNurseCtrl(db).UpdateProfile)
		}

	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"success": "0", "message": notFoundMessage})
	})

	return r
}
