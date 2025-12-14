package user

import (
	"UserStorage/dbhandler"
	"UserStorage/models"
	"UserStorage/queueHandler"
	"UserStorage/secutiry"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewUserHandler(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	testUser := models.User{
		Email:    "test@test.pl",
		Username: "test",
		Age:      21,
		Password: "test",
		Files:    nil,
	}
	MockJsonPost(ctx, testUser, "")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testDB.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
	testMQ.EXPECT().Publish(gomock.Any())
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testObj.CreateUser(ctx)
	assert.Equal(t, w.Code, http.StatusCreated)

}

func TestNewUserHandlerErr18(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	testUser := models.User{
		Email:    "test@test.pl",
		Username: "test",
		Age:      1,
		Password: "test",
		Files:    nil,
	}
	MockJsonPost(ctx, testUser, "")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testObj.CreateUser(ctx)
	assert.Equal(t, w.Code, http.StatusBadRequest)
	msgErr := msgErr{}
	err := json.NewDecoder(w.Body).Decode(&msgErr)
	assert.NoError(t, err)
	assert.Equal(t, msgErr.Err, "User age is less than 18")
}

func TestNewUserHandlerErrNoCreate(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	testUser := models.User{
		Email:    "test@test.pl",
		Username: "test",
		Age:      21,
		Password: "test",
		Files:    nil,
	}
	MockJsonPost(ctx, testUser, "")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(fmt.Errorf("user exist"))
	testObj.CreateUser(ctx)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	msgErr := msgErr{}
	err := json.NewDecoder(w.Body).Decode(&msgErr)
	assert.NoError(t, err)
	assert.Equal(t, msgErr.Err, "user exist")
}

func TestGetUsr(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	testUser := models.User{
		Email:    "test@test.pl",
		Username: "test",
		Age:      21,
		Password: "test",
		Files:    nil,
	}
	MockJsonPost(ctx, testUser, "")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
	testMQ.EXPECT().Publish(gomock.Any())
	testObj.CreateUser(ctx)
	assert.Equal(t, w.Code, http.StatusCreated)
	w = httptest.NewRecorder()
	ctx = GetTestGinContext(w)
	params := gin.Params{gin.Param{Key: "id", Value: testUser.Email}}
	urlTest := url.Values{}
	urlTest.Add("id", testUser.Email)
	MockJsonGet(ctx, params, urlTest)
	testDB.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(models.User{}, fmt.Errorf("user not found"))
	testObj.GetUser(ctx)
	assert.Equal(t, w.Code, http.StatusNotFound)
	msgErr := msgErr{}
	err := json.NewDecoder(w.Body).Decode(&msgErr)
	assert.NoError(t, err)
	assert.Equal(t, msgErr.Err, "user not found", "")
}

func TestGetUserNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	testUser := models.User{
		Email:    "test@test.pl",
		Username: "test",
		Age:      21,
		Password: "test",
		Files:    nil,
	}
	MockJsonPost(ctx, testUser, "")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
	testMQ.EXPECT().Publish(gomock.Any())
	testObj.CreateUser(ctx)
	assert.Equal(t, w.Code, http.StatusCreated)
	w = httptest.NewRecorder()
	ctx = GetTestGinContext(w)
	params := gin.Params{gin.Param{Key: "id", Value: testUser.Email}}
	urlTest := url.Values{}
	urlTest.Add("id", testUser.Email)
	MockJsonGet(ctx, params, urlTest)
	retUsr := testUser
	retUsr.Password = "$2a$10$tuHoVjOL8nzbYreznZZGv.lIiF5ET0glp07i7iDNb/ataCE08u.lC"
	testDB.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(retUsr, nil)
	testObj.GetUser(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	msgErr := models.User{}
	err := json.NewDecoder(w.Body).Decode(&msgErr)
	assert.NoError(t, err)
	assert.Equal(t, msgErr.Password, "")
}

func TestGetAllUsers(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	MockJsonGet(ctx, gin.Params{}, url.Values{})
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().GetUsers(gomock.Any()).Return([]models.User{{
		Email:    "test2@email.com",
		Username: "test1",
		Age:      22,
		Password: "",
		Files:    nil,
	}, {
		Email:    "test2@email.com",
		Username: "test2",
		Age:      55,
		Password: "",
		Files:    nil,
	},
	}, nil)
	testObj.GetAllUsers(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	var msgErr []models.User
	err := json.NewDecoder(w.Body).Decode(&msgErr)
	assert.NoError(t, err)
	assert.Equal(t, len(msgErr), 2)
}

func TestUpdUser(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	testUser := models.User{
		Email:    "test@test.pl",
		Username: "test",
		Age:      21,
		Password: "test",
		Files:    nil,
	}
	MockJsonPost(ctx, testUser, "")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testDB.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)
	testMQ.EXPECT().Publish(gomock.Any())
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testObj.UpdateUser(ctx)
	assert.Equal(t, w.Code, http.StatusOK)

}

func TestUpdUserErr(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	testUser := models.User{
		Email:    "test@test.pl",
		Username: "test",
		Age:      21,
		Password: "test",
		Files:    nil,
	}
	MockJsonPost(ctx, testUser, "")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testDB.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(fmt.Errorf("user not found"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testObj.UpdateUser(ctx)
	assert.Equal(t, w.Code, http.StatusInternalServerError)

}

func TestDeleteUsr(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	MockJsonDelete(ctx, gin.Params{gin.Param{Key: "id", Value: "test@email.com"}})
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(nil)
	testMQ.EXPECT().Publish(gomock.Any())
	testObj.DeleteUser(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	msgOut := msgInf{}
	err := json.NewDecoder(w.Body).Decode(&msgOut)
	assert.NoError(t, err)
	assert.Equal(t, msgOut.Message, "User deleted", "")
}

func TestDeleteUsrErr(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	MockJsonDelete(ctx, gin.Params{gin.Param{Key: "id", Value: "test@email.com"}})
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Return(fmt.Errorf("user not found"))
	testObj.DeleteUser(ctx)
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	msgOut := msgErr{}
	err := json.NewDecoder(w.Body).Decode(&msgOut)
	assert.NoError(t, err)
	assert.Equal(t, msgOut.Err, "user not found", "")
}

func TestGetUsrFiles(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	MockJsonGet(ctx, gin.Params{gin.Param{
		Key:   "id",
		Value: "test@mail.com",
	}}, url.Values{})
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().GetUserFiles(gomock.Any(), gomock.Any()).Return([]models.File{{Name: "testFile1"}, {Name: "testFile2"}}, nil)
	testObj.GetUserFiles(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	var files []models.File
	err := json.NewDecoder(w.Body).Decode(&files)
	assert.NoError(t, err)
	assert.Equal(t, len(files), 2)
}

func TestGetUsrFilesErr(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	MockJsonGet(ctx, gin.Params{gin.Param{
		Key:   "id",
		Value: "test@mail.com",
	}}, url.Values{})
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().GetUserFiles(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("user not found"))
	testObj.GetUserFiles(ctx)
	assert.Equal(t, w.Code, http.StatusNotFound)
	var msgErrOut msgErr
	err := json.NewDecoder(w.Body).Decode(&msgErrOut)
	assert.NoError(t, err)
	assert.Equal(t, msgErrOut.Err, "user not found", "")
}

func TestAddUsrFile(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	MockJsonPost(ctx, models.File{Name: "testFile"}, "test@email.com")
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().AddFileToUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	testObj.AddFileToUser(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	var msgErrOut msgInf
	err := json.NewDecoder(w.Body).Decode(&msgErrOut)
	assert.NoError(t, err)
	assert.Equal(t, msgErrOut.Message, "file added", "")
}

func TestDeleteUSerFiles(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := GetTestGinContext(w)
	var logger = logrus.New()
	MockJsonDelete(ctx, gin.Params{gin.Param{
		Key:   "id",
		Value: "test@mail.com",
	}})
	ctrl := gomock.NewController(t)
	testDB := dbhandler.NewMockDBHandler(ctrl)
	testMQ := queueHandler.NewMockQueueHandler(ctrl)
	auth := secutiry.NewAuthObj([]byte("test"))
	testObj := NewUserHandler(logger, testDB, testMQ, auth)
	testDB.EXPECT().DeleteFilesFromUser(gomock.Any(), gomock.Any()).Return(nil)
	testObj.DeleteFilesFromUser(ctx)
	assert.Equal(t, w.Code, http.StatusOK)
	var msgErrOut msgInf
	err := json.NewDecoder(w.Body).Decode(&msgErrOut)
	assert.NoError(t, err)
	assert.Equal(t, msgErrOut.Message, "files deleted", "")
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

func MockJsonPost(c *gin.Context, content interface{}, id string) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	if id != "" {
		c.Params = gin.Params{gin.Param{Key: "id", Value: id}}
	}

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

type msgErr struct {
	Err string `json:"error"`
}

type msgInf struct {
	Message string `json:"message"`
}

func MockJsonGet(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("id", "")

	// set path params
	c.Params = params

	// set query params
	c.Request.URL.RawQuery = u.Encode()
}

func MockJsonDelete(c *gin.Context, params gin.Params) {
	c.Request.Method = "DELETE"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)
	c.Params = params
}
