package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePasswordHash(t *testing.T) {
	_, err := GeneratePasswordHash([]byte("ahmad"))
	assert.Nil(t, err)
}

func TestPasswordCompare(t *testing.T) {
	password := "ahmad"
	hasPass, _ := GeneratePasswordHash([]byte(password))
	err := PasswordCompare([]byte(password), []byte(hasPass))
	assert.Nil(t, err)
}

func TestSetTokenIdentifier(t *testing.T) {
	str := SetTokenIdentifier("auth", "admin_table")
	assert.Equal(t, str, "auth::admin_table")
}
func TestGetTokenIdentifier(t *testing.T) {
	guard := "admin"
	token := "123wqqeqweqweq12312kl31lk3kl123kl12j3kl12j3klj123kl"
	fullToken := "admin::" + token

	validToken := GetTokenIdentifier(guard, fullToken)
	assert.Equal(t, validToken, token)
}
