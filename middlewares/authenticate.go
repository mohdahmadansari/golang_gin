package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mohdahmadansari/golang_gin/helpers"
	"github.com/mohdahmadansari/golang_gin/models"
	"github.com/mohdahmadansari/golang_gin/services"
	"gorm.io/gorm"
)

func Authenticate(guard string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		requiredToken := c.Request.Header["Authorization"]
		fmt.Println(requiredToken)
		if len(requiredToken) == 0 {
			helpers.ResponseJsonError(c, 403, "Invalid token.")
			return
		}

		decodedToken, _ := services.DecodeToken(requiredToken[0])
		// helpers.ResponseJsonError(c, 403, err.Error())
		username := helpers.GetTokenIdentifier(guard, decodedToken)

		if username == "" {
			helpers.ResponseJsonError(c, 403, "Access denied.")
			return
		}

		if guard == "admin" {
			var Admin models.Admin

			if err := db.First(&Admin, "username = ?", username).Error; err != nil {
				helpers.ResponseJsonError(c, 403, "User does not exists.")
				return
			}

			c.Set("AdminData", Admin)
		}

		if guard == "nurse" {
			var nurse models.Nurse

			if err := db.First(&nurse, "username = ?", username).Error; err != nil {
				helpers.ResponseJsonError(c, 403, "User does not exists.")
				return
			}

			c.Set("NurseData", nurse)
		}
		c.Next()
	}
}
