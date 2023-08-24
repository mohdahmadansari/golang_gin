package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohdahmadansari/golang_gin/database"
	"github.com/mohdahmadansari/golang_gin/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type LoginResponse struct {
	Message       string `json:"message"`
	Refresh_token string `json:"refresh_token"`
	Token         string `json:"token"`
}

func TestAdmin(t *testing.T) {

	db := getDb()

	r := gin.Default()
	req, err := createRequest("GET", "/admin", "")
	rr := httptest.NewRecorder()
	r.GET("/admin", NewAdminCtr(db).Admin)
	r.ServeHTTP(rr, req)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, rr.Code)
}

func TestLogin_invalidRequest(t *testing.T) {

	db := getDb()
	r := gin.Default()

	// invalid request
	jsonBody := `{"username": "incorrect_username"}`
	req, _ := createRequest("POST", "/login", jsonBody)
	resInvalid := httptest.NewRecorder()
	r.POST("/login", NewAdminCtr(db).Login)
	r.ServeHTTP(resInvalid, req)

	assert.EqualValues(t, http.StatusBadRequest, resInvalid.Code)
}
func TestLogin_invalidUser(t *testing.T) {

	db := getDb()
	r := gin.Default()

	jsonBody := `{"username": "incorrect_username", "password": "passssss"}`
	req, _ := createRequest("POST", "/login", jsonBody)
	res := httptest.NewRecorder()
	r.POST("/login", NewAdminCtr(db).Login)
	r.ServeHTTP(res, req)

	assert.EqualValues(t, http.StatusBadRequest, res.Code)
}

func TestLogin_incorrectPassword(t *testing.T) {

	db := getDb()
	r := gin.Default()

	jsonBody := `{"username": "ahmad", "password": "passssss"}`
	req, _ := createRequest("POST", "/login", jsonBody)
	res := httptest.NewRecorder()
	r.POST("/login", NewAdminCtr(db).Login)
	r.ServeHTTP(res, req)

	assert.EqualValues(t, http.StatusBadRequest, res.Code)
}

func TestLogin_success(t *testing.T) {

	db := getDb()
	database.SeedDatabase(db)
	r := gin.Default()

	jsonBody := `{"username": "ahmad", "password": "click123"}`
	req, _ := createRequest("POST", "/login", jsonBody)
	res := httptest.NewRecorder()
	r.POST("/login", NewAdminCtr(db).Login)
	r.ServeHTTP(res, req)
	fmt.Println(res.Body)
	var loginResponse LoginResponse
	err := json.Unmarshal(res.Body.Bytes(), &loginResponse)

	fmt.Println(loginResponse)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.Code)
	assert.NotEmpty(t, loginResponse.Token)
}

func TestDashboard_validUser(t *testing.T) {

	db := getDb()

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	NewAdminCtr(db).Dashboard(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)

}

func TestDashboard_invalidUser(t *testing.T) {

	db := getDb()

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, true)
	NewAdminCtr(db).Dashboard(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func createRequest(method string, route string, jsonBody string) (req *http.Request, err error) {
	var body io.Reader = nil
	if jsonBody != "" {
		body = bytes.NewBufferString(jsonBody)
	}
	req, err = http.NewRequest(method, route, body)
	req.Header.Add("Content-Type", "application/json")
	// println(err)
	return
}

func getDb() (db *gorm.DB) {
	connStr, _ := database.GetConnectionString()
	dbObj, _ := database.CreateConnection(connStr)
	database.SetupMigration(dbObj)
	database.SeedDatabase(dbObj)

	return dbObj
}

func MockAdmin(c *gin.Context, blank bool) {
	var admin models.Admin
	if !blank {
		admin.Email = "ahmad@gmail.com"
		admin.Password = "click123"
		admin.Username = "ahmad"
		admin.ID = 1
	}
	c.Set("AdminData", admin)
}

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}
