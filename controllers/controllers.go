package controllers

import "gorm.io/gorm"

type Controllers struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) *Controllers {
	return &Controllers{
		db: db,
	}
}
