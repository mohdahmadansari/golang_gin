package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password []byte) string {
	// default cost is 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		ErrorPanic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Info().Msg("issue with password encryption, the panic gracefully")
		}
	}()

	return string(hashedPassword)
}

func PasswordCompare(password []byte, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	return err
}

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func SetTokenIdentifier(guard string, tbl_identifier string) string {
	var adminToken = guard + "::" + tbl_identifier
	return adminToken
}

func GetTokenIdentifier(guard string, token string) string {

	// log.Info().Msg("guard" + guard + " -- " + token)
	var validToken = ""
	var adminToken = guard + "::"
	var tokenSlice = strings.Split(token, adminToken)
	if len(tokenSlice) == 2 {
		validToken = tokenSlice[1]
	}
	// log.Info().Msg("tokenSlice" + tokenSlice[0] + tokenSlice[1])
	return validToken
}

func ResponseJsonError(c *gin.Context, code int, s string) {
	c.AbortWithStatusJSON(code, gin.H{"success": 0, "message": s})
	// return
}

func ResponseJsonSuccess(c *gin.Context, H map[string]any) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{})
}

func Console(msg any) {
	fmt.Println(msg)
}
