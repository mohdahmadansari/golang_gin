package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohdahmadansari/golang_gin/helpers"
	"github.com/mohdahmadansari/golang_gin/models"
	"github.com/mohdahmadansari/golang_gin/services"
	"gorm.io/gorm"
)

type NurseCtrl struct {
	db *gorm.DB
}

func NewNurseCtrl(db *gorm.DB) *NurseCtrl {
	return &NurseCtrl{
		db: db,
	}
}

// @Security BearerAuth
// @Summary Get all Nurses
// @Description Get all Nurses
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Nurse
// @Router /admin/nurse [get]
func (ctrl *NurseCtrl) Get(c *gin.Context) {

	results, err := models.GetAllNurses(ctrl.db)
	// var Admin = c.MustGet("AdminData").(models.Admin)
	// results, err := models.GetOwnNurses(ctrl.db, int(Admin.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "No records found."})
		return
	}
	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": 1, "message": "No records found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Nurse List.", "data": results})
}

// @Security BearerAuth
// @Summary Get Nurse Details
// @Description Get Nurse Details
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "ID"
// @Success 200 {object} models.Nurse
// @Failure 400 {string} string "400 StatusBadRequest"
// @Failure 404 {string} string "404 not found"
// @Router /admin/nurse/{id} [get]
func (ctrl *NurseCtrl) GetOne(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
		return
	}
	var nurse models.Nurse
	nurse.ID = uint(id)
	result, err := models.GetNursesById(ctrl.db, &nurse)

	if err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "No record found."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Nurse details.", "data": result})
}

// @Security BearerAuth
// @Summary Add New Nurse
// @Description Add New Nurse
// @Accept  json
// @Produce  json
// @Param role body models.Nurse true "data"
// @Success 200 {object} models.Nurse
// @Router /admin/nurse [post]
func (ctrl *NurseCtrl) Post(c *gin.Context) {
	var t models.Nurse

	var Admin = c.MustGet("AdminData").(models.Admin)
	t.Admin = Admin
	if err := c.ShouldBindJSON(&t); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
		return
	}

	var Nurse models.Nurse
	if isDuplicate := ctrl.db.First(&Nurse, "username = ?", t.Username).Error; isDuplicate == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "Username does not exists."})
		return
	}
	t.Password, _ = helpers.GeneratePasswordHash([]byte(t.Password))
	if err_create := ctrl.db.Create(&t).Error; err_create != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err_create.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Nurse created successfully", "data": t})
}

// @Security BearerAuth
// @Summary Update Nurse
// @Description Update Nurse
// @Accept  json
// @Produce  json
// @Param role body models.Nurse true "data"
// @Param   id     path    int     true        "ID"
// @Success 200 {object} models.Nurse
// @Failure 400 {string} string "400 StatusBadRequest"
// @Failure 404 {string} string "404 not found"
// @Router /admin/nurse/{id} [put]
func (ctrl *NurseCtrl) Put(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
		return
	}

	updateNurse := make(map[string]interface{})
	updateNurse["id"] = uint(id)

	c.ShouldBindJSON(&updateNurse)
	fmt.Println(updateNurse)
	// if err := c.ShouldBindJSON(&updateNurse); err != nil {
	// c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
	// return
	// }

	var nurse models.Nurse
	nurse.ID = uint(id)
	result, err := models.GetNursesById(ctrl.db, &nurse)

	if err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "No record found."})
		return
	}

	username, has_username := updateNurse["username"]
	if has_username && !models.IsUsernameAvailable(ctrl.db, nurse.ID, username.(string)) {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "Username already exists."})
		return
	}

	if password, okay := updateNurse["password"]; okay {
		newPassword, _ := helpers.GeneratePasswordHash([]byte(password.(string)))
		updateNurse["password"] = newPassword
	}

	ctrl.db.Model(&nurse).Updates(&updateNurse)

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Nurse updated successfully", "data": updateNurse})
}

// @Security BearerAuth
// @Summary Delete Nurse
// @Description Delete Nurse
// @Accept  json
// @Produce  json
// @Param   id     path    int     true        "ID"
// @Success 204
// @Failure 400 {string} string "400 StatusBadRequest"
// @Failure 404 {string} string "404 not found"
// @Router /admin/nurse/{id} [delete]
func (ctrl *NurseCtrl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
		return
	}

	var nurse models.Nurse
	nurse.ID = uint(id)
	result, err := models.GetNursesById(ctrl.db, &nurse)

	if err != nil || len(result) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "No record found."})
		return
	}

	ctrl.db.Model(&nurse).Unscoped().Delete(&nurse)

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Nurse deleted successfully"})
}

// @Summary Nurse Login
// @Description Login Nurses
// @Accept  json
// @Produce  json
// @Param role body models.NurseLogin true "data"
// @Success 200 {object} models.NurseLogin
// @Router /nurse/login [post]
func (ctrl *NurseCtrl) Login(c *gin.Context) {
	var t models.NurseLogin

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": err.Error()})
		return
	}

	var nurse models.Nurse

	if err := ctrl.db.First(&nurse, "username = ?", t.Username).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "User does not exists."})
		return
	}

	hashedPassword := []byte(nurse.Password)
	password := []byte(t.Password)

	err := helpers.PasswordCompare(password, hashedPassword)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "Invalid username and password."})
		return
	}

	var nurseToken = helpers.SetTokenIdentifier("nurse", nurse.Username)
	jwtToken, refreshToken, jwtErr := services.GenerateToken(nurseToken)

	if jwtErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "There are issue in access."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Login Successfully!", "token": jwtToken, "refresh_token": refreshToken})
}

// @Security BearerAuth
// @Summary Get Profile details
// @Description Nurse profile details
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Nurse
// @Failure 400 {string} string "400 StatusBadRequest"
// @Failure 404 {string} string "404 not found"
// @Router /nurse/profile [get]
func (ctrl *NurseCtrl) GetProfile(c *gin.Context) {

	var nurse = c.MustGet("NurseData").(models.Nurse)

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Nurse details.", "data": nurse})
}

// @Security BearerAuth
// @Summary Update Nurse profile
// @Description Update Nurse profile
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Nurse
// @Param role body models.Nurse true "Update Profile"
// @Failure 400 {string} string "400 StatusBadRequest"
// @Failure 404 {string} string "404 not found"
// @Router /nurse/profile [post]
func (ctrl *NurseCtrl) UpdateProfile(c *gin.Context) {
	var nurse = c.MustGet("NurseData").(models.Nurse)

	// var r models.Nurse
	// r.ID = uint(id)
	updateNurse := make(map[string]interface{})
	updateNurse["id"] = nurse.ID

	c.ShouldBindJSON(&updateNurse)

	username, has_username := updateNurse["username"]
	if has_username && !models.IsUsernameAvailable(ctrl.db, nurse.ID, username.(string)) {
		c.JSON(http.StatusBadRequest, gin.H{"success": 0, "message": "Username already exists."})
		return
	}

	if password, okay := updateNurse["password"]; okay {
		newPassword, _ := helpers.GeneratePasswordHash([]byte(password.(string)))
		updateNurse["password"] = newPassword
	}

	ctrl.db.Model(&nurse).Updates(&updateNurse)

	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Profile Updated successfully", "data": updateNurse})
}

// @Security BearerAuth
// @Summary Get Own Nurses
// @Description Get Own Nurses
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Nurse
// @Router /admin/nurse/own [get]
func (ctrl *NurseCtrl) Getown(c *gin.Context) {

	// results, err := models.GetAllNurses(ctrl.db)
	var Admin = c.MustGet("AdminData").(models.Admin)
	results, err := models.GetOwnNurses(ctrl.db, int(Admin.ID))
	if err != nil || len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": 1, "message": "No records found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Nurse List.", "data": results})
}
