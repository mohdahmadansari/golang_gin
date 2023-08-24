package models_test

import (
	"testing"

	"github.com/mohdahmadansari/golang_gin/database"
	"github.com/mohdahmadansari/golang_gin/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetAllNurses(t *testing.T) {
	db := getDb()
	nurses, err := models.GetAllNurses(db)
	assert.Nil(t, err)
	assert.NotEmpty(t, nurses)
}

func TestGetOwnNurses(t *testing.T) {
	db := getDb()
	nurses, err := models.GetOwnNurses(db, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, nurses)
}

func TestGetNursesById(t *testing.T) {
	var n models.Nurse
	n.ID = 1
	db := getDb()
	nurses, err := models.GetNursesById(db, &n)
	assert.Nil(t, err)
	assert.NotEmpty(t, nurses)
}

func TestIsUsernameAvailable(t *testing.T) {

	db := getDb()

	invalid := models.IsUsernameAvailable(db, uint(1), "ahmad1")
	valid := models.IsUsernameAvailable(db, uint(1), "ahmad1231231")

	assert.False(t, invalid)
	assert.True(t, valid)
}

func getDb() (db *gorm.DB) {
	connStr, _ := database.GetConnectionString()
	dbObj, _ := database.CreateConnection(connStr)
	database.SetupMigration(dbObj)
	database.SeedDatabase(dbObj)

	return dbObj
}
