package services_test

import (
	"testing"

	"github.com/mohdahmadansari/golang_gin/services"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tokenString, refreshTokenString, err := services.GenerateToken("mohdahmad")
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)
	assert.NotEmpty(t, refreshTokenString)
}

func TestGenerateNonAuthToken(t *testing.T) {
	tokenString, err := services.GenerateNonAuthToken("mohdahmad")
	assert.Nil(t, err)
	assert.NotEmpty(t, tokenString)
}
func TestDecodeNonAuthToken(t *testing.T) {
	username := "mohdahmad"
	tokenString, _ := services.GenerateNonAuthToken(username)
	decodeString, err := services.DecodeNonAuthToken(tokenString)
	assert.Nil(t, err)
	assert.EqualValues(t, decodeString, username)
}

func TestDecodeToken(t *testing.T) {
	username := "mohdahmad"
	tokenString, _, _ := services.GenerateToken(username)
	decodeString, err := services.DecodeToken(tokenString)
	assert.Nil(t, err)
	assert.EqualValues(t, decodeString, username)
}

func TestDecodeRefreshToken(t *testing.T) {
	username := "mohdahmad"
	_, refreshTokenString, _ := services.GenerateToken(username)
	decodeString, err := services.DecodeRefreshToken(refreshTokenString)
	assert.Nil(t, err)
	assert.EqualValues(t, decodeString, username)
}
