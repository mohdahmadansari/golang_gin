package models

import (
	"gorm.io/gorm"
)

type Nurse struct {
	gorm.Model
	Created_by int    `json:"-"`
	Admin      Admin  `json:"admin,omitempty" gorm:"foreignKey:Created_by;references:ID"`
	Username   string `json:"username" binding:"required" gorm:"index:id,unique"`
	Password   string `json:"password" binding:"required"`

	Email   string `json:"email" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone"`
	Address string `json:"address" `
}

type NurseLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetAllNurses(db *gorm.DB) ([]Nurse, error) {
	var nurses []Nurse

	err := db.Model(&Nurse{}).Preload("Admin").Find(&nurses).Error
	return nurses, err
}

func GetOwnNurses(db *gorm.DB, ownerId int) ([]Nurse, error) {
	var nurses []Nurse
	err := db.Model(&Nurse{}).Preload("Admin").Where("created_by = ?", ownerId).Find(&nurses).Error
	return nurses, err
}

func GetNursesById(db *gorm.DB, nurse *Nurse) ([]Nurse, error) {
	var nurses []Nurse
	err := db.Model(&Nurse{}).Preload("Admin").Where("id = ?", nurse.ID).Find(&nurses).Error
	return nurses, err
}

func IsUsernameAvailable(db *gorm.DB, ownerId uint, username string) bool {
	var isAvail = false
	var nurse []Nurse
	if isDuplicate := db.Where("username = ?", username).Where("id != ?", ownerId).First(&nurse).Error; isDuplicate != nil {
		isAvail = true
	}
	return isAvail
}
