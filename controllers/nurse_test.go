package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohdahmadansari/golang_gin/helpers"
	"github.com/mohdahmadansari/golang_gin/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGet_err(t *testing.T) {

	db := getDb()
	db.Migrator().DropTable(&models.Nurse{})

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	NewNurseCtrl(db).Get(ctx)

	db.AutoMigrate(&models.Nurse{})

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func TestGet_zero_row(t *testing.T) {

	db := getDb()

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	NewNurseCtrl(db).Get(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)

}

func TestGet_with_results(t *testing.T) {

	var defaultPassword = "click123"
	newHashedPassword, _ := helpers.GeneratePasswordHash([]byte(defaultPassword))
	nurse := []models.Nurse{
		{Email: "ahmadweb2011@gmail.com", Username: "ahmad", Password: newHashedPassword, Name: "Ahmad", Phone: "9718224606", Address: "159/2"},
		{Email: "ahmadweb2012@gmail.com", Username: "ahmad1", Password: newHashedPassword, Name: "Ahmad1", Phone: "9718224606", Address: "159/2"},
	}

	db := getDb()
	db.Migrator().DropTable(&models.Nurse{})
	db.AutoMigrate(&models.Nurse{})
	var admin models.Admin
	admin.ID = 1
	for _, a := range nurse {
		a.Admin = admin
		db.Create(&a)
	}

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	NewNurseCtrl(db).Get(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)

}

func TestGetOne_invalid_request(t *testing.T) {

	db := getDb()

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	NewNurseCtrl(db).GetOne(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func TestGetOne_valid_request(t *testing.T) {

	db := getDb()

	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}
	u := url.Values{}

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockingJsonGet(ctx, params, u)
	NewNurseCtrl(db).GetOne(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)

}

func TestGetOne_valid_request_invalid_id(t *testing.T) {

	db := getDb()

	params := []gin.Param{
		{
			Key:   "id",
			Value: "1000000011233",
		},
	}
	u := url.Values{}

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockingJsonGet(ctx, params, u)
	NewNurseCtrl(db).GetOne(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func TestPut_invalid_request(t *testing.T) {

	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	NewNurseCtrl(db).Put(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func TestPut_valid_request(t *testing.T) {

	db := getDb()

	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockJsonPut(ctx, nil, params)
	NewNurseCtrl(db).Put(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)

}

func TestPut_valid_request_invalid_id(t *testing.T) {

	db := getDb()

	params := []gin.Param{
		{
			Key:   "id",
			Value: "1000654545645600",
		},
	}

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockJsonPut(ctx, nil, params)
	NewNurseCtrl(db).Put(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)

}

func TestPut_chk_duplicate_username_update(t *testing.T) {
	db := getDb()

	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}
	nurseData := &models.Nurse{Username: "ahmad1"} // added other existed username

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockJsonPut(ctx, &nurseData, params)
	NewNurseCtrl(db).Put(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestPut_update_password(t *testing.T) {
	db := getDb()

	params := []gin.Param{
		{
			Key:   "id",
			Value: "1",
		},
	}
	nurseData := &models.Nurse{Password: "click123"} // added other existed username

	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockJsonPut(ctx, &nurseData, params)
	NewNurseCtrl(db).Put(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestPost_invalid_request(t *testing.T) {

	nurseData := &models.Nurse{Name: "Ahmad"}

	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockJsonPost(ctx, &nurseData, false)
	NewNurseCtrl(db).Post(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}
func TestPost_valid_request_same_username(t *testing.T) {

	// admin := &models.Admin{Email: "ahmad@gmail.com", Username: "ahmad", Password: "click123"}
	nurseData := &models.Nurse{Name: "Ahmad New", Username: "ahmad", Password: "click123", Email: "a@gmail.com", Phone: "9718224606", Address: "159/222222"}
	nurseData.Admin.ID = 1
	nurseData.Admin.Email = "ahmad@gmail.com"
	nurseData.Admin.Username = "ahmad"
	nurseData.Admin.Password = "click123"
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockJsonPost(ctx, &nurseData, false)
	NewNurseCtrl(db).Post(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestPost_create_nurse_succsssfully(t *testing.T) {

	now := time.Now()
	unixMilli := now.UnixMilli()

	// admin := &models.Admin{Email: "ahmad@gmail.com", Username: "ahmad", Password: "click123"}
	nurseData := &models.Nurse{Name: "Ahmad New", Username: "ahmad555" + strconv.Itoa(int(unixMilli)), Password: "click123", Email: "a@gmail.com", Phone: "9718224606", Address: "159/222222"}
	nurseData.Admin.ID = 1
	nurseData.Admin.Email = "ahmad@gmail.com"
	nurseData.Admin.Username = "ahmad"
	nurseData.Admin.Password = "click123"
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockJsonPost(ctx, &nurseData, false)
	NewNurseCtrl(db).Post(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestLogin_nurse_invalid_request(t *testing.T) {

	nurse := &models.NurseLogin{Username: "ahmad"}
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockJsonPost(ctx, &nurse, false)
	NewNurseCtrl(db).Login(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestLogin_nurse_invalid_username(t *testing.T) {

	nurse := &models.NurseLogin{Username: "ahmad123", Password: "click123"}
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockJsonPost(ctx, &nurse, false)
	NewNurseCtrl(db).Login(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestLogin_nurse_invalid_password(t *testing.T) {

	nurse := &models.NurseLogin{Username: "ahmad", Password: "click1123"}
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockJsonPost(ctx, &nurse, false)
	NewNurseCtrl(db).Login(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestLogin_nurse_success(t *testing.T) {

	nurse := &models.NurseLogin{Username: "ahmad", Password: "click123"}
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockJsonPost(ctx, &nurse, false)
	NewNurseCtrl(db).Login(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestUpdateProfile(t *testing.T) {

	// nurseData := &models.Nurse{Name: "Ahmad update profile test"}
	updateNurse := make(map[string]interface{})
	updateNurse["name"] = "Ahmad update profile test"
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockNurse(ctx, false)
	MockJsonPost(ctx, &updateNurse, false)
	NewNurseCtrl(db).UpdateProfile(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestUpdateProfile_chk_username_update(t *testing.T) {

	updateNurse := make(map[string]interface{})
	updateNurse["username"] = "ahmad1" //already exists username corresponding to ID:2
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockNurse(ctx, false)
	MockJsonPost(ctx, &updateNurse, false)
	NewNurseCtrl(db).UpdateProfile(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProfile_password_update(t *testing.T) {

	updateNurse := make(map[string]interface{})
	updateNurse["password"] = "ahmad123"
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockNurse(ctx, false)
	MockJsonPost(ctx, &updateNurse, false)
	NewNurseCtrl(db).UpdateProfile(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestGetProfile(t *testing.T) {

	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockNurse(ctx, false)
	MockingJsonGet(ctx, nil, nil)
	NewNurseCtrl(db).GetProfile(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func TestGetown_admin(t *testing.T) {

	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockingJsonGet(ctx, nil, nil)
	NewNurseCtrl(db).Getown(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}
func TestDelete_without_id(t *testing.T) {

	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockingJsonDelete(ctx, nil, nil)
	NewNurseCtrl(db).Delete(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestDelete_invalid_nurse(t *testing.T) {
	params := []gin.Param{
		{
			Key:   "id",
			Value: "1231231231232300",
		},
	}
	db := getDb()
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockingJsonDelete(ctx, params, nil)
	NewNurseCtrl(db).Delete(ctx)

	assert.EqualValues(t, http.StatusBadRequest, w.Code)
}

func TestDelete_success(t *testing.T) {
	// creating new nurse and delete that
	TestPost_create_nurse_succsssfully(t)
	var nurseLast models.Nurse

	db := getDb()
	db.Last(&nurseLast)

	params := []gin.Param{
		{
			Key:   "id",
			Value: strconv.Itoa(int(nurseLast.ID)),
		},
	}
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	MockAdmin(ctx, false)
	MockingJsonDelete(ctx, params, nil)
	NewNurseCtrl(db).Delete(ctx)

	assert.EqualValues(t, http.StatusOK, w.Code)
}

func MockNurse(c *gin.Context, blank bool) {
	var nurse models.Nurse
	if !blank {
		nurse.Username = "ahmad"
		nurse.ID = 1
	}
	c.Set("NurseData", nurse)
}

func MockingJsonGet(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	// if setAdmin {
	// 	MockAdmin(c, false, db)
	// }
	c.Params = params
	c.Request.URL.RawQuery = u.Encode()
}
func MockingJsonDelete(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "DELETE"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = params
	c.Request.URL.RawQuery = u.Encode()
}
func MockJsonPost(c *gin.Context, content interface{}, adminData bool) {
	var admin models.Admin
	admin.Username = "ahmad"
	admin.ID = 1

	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)
	if adminData {
		c.Set("AdminData", &admin)
	}

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
func MockingJsonPost(c *gin.Context, content interface{}, setAdmin bool) {
	var admin models.Admin
	admin.Username = "ahmad"
	admin.ID = 1

	c.Request.Method = http.MethodPost
	// c.Request.Header.Set("Content-Type", "application/json")
	c.Header("Content-Type", "application/vnd.api+json")
	c.Set("AdminData", admin)
	// body := bytes.NewBufferString(jsonBody)
	// body := new(bytes.Buffer)
	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonbytes))
	// body := bytes.NewBufferString(`"address":"159/2","email":"ahmadweb2011@gmail.com","name":"Ahmad","password":"ahmad","phone":"9718224606","username":"ahmad"`)
	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	// c.BindJSON(jsonbytes)
	// c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	c.Request.Body = io.NopCloser(bytes.NewBufferString("{\"address\":\"159/2\",\"email\":\"ahmadweb2011@gmail.com\",\"name\":\"Ahmad\",\"password\":\"ahmad\",\"phone\":\"9718224606\",\"username\":\"ahmad\"}"))
	// c.Request.Body = io.NopCloser(body)
	// c.Request.Body = io.NopCloser(body)
	// c.ContentType()

	// fmt.Println(c.Request.Body)
}

func MockJsonPut(c *gin.Context, content interface{}, params gin.Params) {
	c.Request.Method = "PUT"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)
	c.Params = params

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

func MakeRequest(method string, url string, jsonBody string, isAuthenticatedRequest bool, db *gorm.DB, ctrMethod func(c *gin.Context)) *httptest.ResponseRecorder {
	r := gin.Default()
	// requestBody, _ := json.Marshal(body)
	// request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	var body io.Reader = nil
	if jsonBody != "" {
		body = bytes.NewBufferString(jsonBody)
	}
	request, _ := http.NewRequest(method, url, body)
	request.Header.Add("Content-Type", "application/json")
	if isAuthenticatedRequest {
		request.Header.Add("Authorization", "Bearer "+BearerToken(db))
	}
	writer := httptest.NewRecorder()
	if method == "POST" {
		r.POST(url, ctrMethod)
	} else if method == "GET" {
		r.GET(url, ctrMethod)
	} else if method == "PUT" {
		r.PUT(url, ctrMethod)
	} else if method == "DELETE" {
		r.DELETE(url, ctrMethod)
	}

	r.ServeHTTP(writer, request)
	return writer
}

func BearerToken(db *gorm.DB) string {
	// adminAuth := models.AdminLogin{
	// 	Username: "ahmad",
	// 	Password: "click123",
	// }
	jsonBody := `{"username": "ahmad", "password": "click123"}`
	writer := MakeRequest("POST", "/login", jsonBody, false, db, NewAdminCtr(db).Login)
	var response map[string]string
	json.Unmarshal(writer.Body.Bytes(), &response)
	fmt.Println(response["token"])
	return response["token"]
}
