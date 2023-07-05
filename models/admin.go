package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"index:id,unique"`
	Password string `json:"password" binding:"required"`
}

type AdminLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
