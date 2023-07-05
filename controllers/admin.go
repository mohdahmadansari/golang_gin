package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohdahmadansari/golang_gin/helpers"
	"github.com/mohdahmadansari/golang_gin/models"
	"github.com/mohdahmadansari/golang_gin/services"
	"gorm.io/gorm"
)

type AdminCtr struct {
	db *gorm.DB
}

func NewAdminCtr(db *gorm.DB) *AdminCtr {
	return &AdminCtr{
		db: db,
	}
}

// @Summary Admin Login
// @Description Admin Login
// @Accept  json
// @Produce  json
// @Param Login body models.AdminLogin true "Login Request"
// @Success 200 {object} models.AdminLogin
// @Router /login [post]
func (ctrl *AdminCtr) Login(c *gin.Context) {
	var t models.AdminLogin

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
		return
	}

	var Admin models.Admin

	if err := ctrl.db.First(&Admin, "username = ?", t.Username).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "User does not exists."})
		return
	}

	hashedPassword := []byte(Admin.Password)
	password := []byte(t.Password)

	err := helpers.PasswordCompare(password, hashedPassword)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "Invalid username and password."})
		return
	}

	var adminToken = helpers.SetTokenIdentifier("admin", Admin.Username)
	jwtToken, refreshToken, jwtErr := services.GenerateToken(adminToken)

	if jwtErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "There are issue in access."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Login Successfully!", "token": jwtToken, "refresh_token": refreshToken})
}

func (ctrl *AdminCtr) Admin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Welcome to go Admin API."})
}

func (ctrl *AdminCtr) Dashboard(c *gin.Context) {
	// helpers.Console(c.Value("AdminData"))
	Admin := c.MustGet("AdminData").(models.Admin)

	if Admin.Username == "" {
		helpers.ResponseJsonError(c, 403, "Admin does not exists.")
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Admin Authenticated for Dashboard landing API.", "Data": Admin})
}
