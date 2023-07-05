package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/mohdahmadansari/golang_gin/helpers"
	"github.com/mohdahmadansari/golang_gin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
}

func CreateConnection() (conn *gorm.DB, e error) {

	queryString := os.Getenv("MYSQL_USERNAME") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE_NAME") + "?charset=utf8&parseTime=True&loc=Local"

	println(queryString)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       queryString,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})

	defer func() {
		if r := recover(); r != nil {
			log.Info().Msg("recovered from the panic")
		}
	}()

	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Error %s when opening DB\n", err))
	}

	return db, err
}

func SetupMigration(db *gorm.DB) {
	// db.AutoMigrate(&Admin{}, &User{}, &Message{})
	db.AutoMigrate(&models.Admin{}, &models.Nurse{})
	if err := db.AutoMigrate(&models.Admin{}); err == nil && db.Migrator().HasTable(&models.Admin{}) {
		if err := db.First(&models.Admin{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			seedAdmin(db)
		}
	}

}

func seedAdmin(db *gorm.DB) {
	var defaultPassword = "click123"
	newHashedPassword := helpers.GeneratePasswordHash([]byte(defaultPassword))
	admins := []models.Admin{
		{Email: "ahmadweb2011@gmail.com", Username: "ahmad", Password: newHashedPassword},
		{Email: "manish@gmail.com", Username: "manish", Password: newHashedPassword},
	}
	for _, a := range admins {
		db.Create(&a)
	}
}
