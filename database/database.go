package database

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mohdahmadansari/golang_gin/helpers"
	"github.com/mohdahmadansari/golang_gin/models"
	"github.com/mohdahmadansari/golang_gin/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
}

func CreateConnection(connStr string) (conn *gorm.DB, e error) {

	// queryString := os.Getenv("MYSQL_USERNAME") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE_NAME") + "?charset=utf8&parseTime=True&loc=Local"
	queryString := connStr

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

func SetupMigration(db *gorm.DB) (err error) {
	// db.AutoMigrate(&Admin{}, &User{}, &Message{})
	db.AutoMigrate(&models.Admin{}, &models.Nurse{})
	err = db.AutoMigrate(&models.Admin{})
	// if err := db.AutoMigrate(&models.Admin{}); err == nil && db.Migrator().HasTable(&models.Admin{}) {
	// 	SeedAdmin(db)
	// }
	return
}

func DropTables(db *gorm.DB) (err error) {
	err = db.Migrator().DropTable(&models.Admin{}, &models.Nurse{})
	return
}

func SeedDatabase(db *gorm.DB) (tx *gorm.DB, err error) {

	if !db.Migrator().HasTable(&models.Admin{}) {
		err = errors.New("migration failed")
	}
	if errtbl := db.First(&models.Admin{}).Error; errors.Is(errtbl, gorm.ErrRecordNotFound) {
		var defaultPassword = "click123"
		newHashedPassword, _ := helpers.GeneratePasswordHash([]byte(defaultPassword))
		admins := []models.Admin{
			{Email: "ahmadweb2011@gmail.com", Username: "ahmad", Password: newHashedPassword},
			{Email: "manish@gmail.com", Username: "manish", Password: newHashedPassword},
		}
		for _, a := range admins {
			tx = db.Create(&a)
			// println(tx)
		}
	} else {
		tx = db.First(&models.Admin{})
	}

	return
}

// func GetConnectionString() (queryString string, err error) {
// 	c, err := util.LoadConfig("../", "app")
// 	if err != nil {
// 		log.Fatal().Msg("Error loading app.env file")
// 	}
// 	// queryString := os.Getenv("MYSQL_USERNAME") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE_NAME") + "?charset=utf8&parseTime=True&loc=Local"
// 	queryString = c.DbUsername + ":" + c.DbPassword + "@tcp(" + c.DbHost + ":" + c.DbPort + ")/" + c.DbName + "?charset=utf8&parseTime=True&loc=Local"
// 	println(queryString)
// 	return queryString, err
// }

func GetConnectionString() (queryString string, err error) {
	c := loadConfigFile("./")
	// queryString := os.Getenv("MYSQL_USERNAME") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE_NAME") + "?charset=utf8&parseTime=True&loc=Local"
	queryString = c.DbUsername + ":" + c.DbPassword + "@tcp(" + c.DbHost + ":" + c.DbPort + ")/" + c.DbName + "?charset=utf8&parseTime=True&loc=Local"
	println(queryString)
	return queryString, err
}

func loadConfigFile(path string) util.Config {
	c, err := util.LoadConfig(path, "app")
	if err != nil {
		return loadConfigFile("../")
	}
	return c
}
