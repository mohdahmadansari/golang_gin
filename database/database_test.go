package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionString(t *testing.T) {

	connStr, err := GetConnectionString()
	assert.NoError(t, err)
	assert.NotEmpty(t, connStr)
}
func TestCreateConnection(t *testing.T) {

	connStr, _ := GetConnectionString()

	db, dbError := CreateConnection(connStr)
	assert.NoError(t, dbError)
	assert.NotEmpty(t, db)
}

func TestDropTables(t *testing.T) {

	connStr, _ := GetConnectionString()
	db, _ := CreateConnection(connStr)

	err := DropTables(db)
	assert.Nil(t, err)
}
func TestSetupMigration(t *testing.T) {

	connStr, _ := GetConnectionString()
	db, _ := CreateConnection(connStr)

	DropTables(db)

	_, beforeMigrationErr := SeedDatabase(db)
	assert.NotNil(t, beforeMigrationErr)

	err := SetupMigration(db)
	assert.NoError(t, err)
}

func TestSeedAdmin(t *testing.T) {

	connStr, _ := GetConnectionString()
	db, _ := CreateConnection(connStr)
	tx, err := SeedDatabase(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, tx)

	txExist, errExist := SeedDatabase(db) // check already seeded data
	assert.NoError(t, errExist)
	assert.NotEmpty(t, txExist)
}
