package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohdahmadansari/golang_gin/database"
	"github.com/stretchr/testify/assert"
)

func TestWelcome(t *testing.T) {

	connStr, _ := database.GetConnectionString()
	db, _ := database.CreateConnection(connStr)
	ctr := NewController(db)

	r := gin.Default()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.GET("/", ctr.Welcome)
	r.ServeHTTP(rr, req)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, rr.Code)
}
