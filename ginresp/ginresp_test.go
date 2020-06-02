package ginresp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/validator.v9"

	"github.com/devit-tel/goerror"
)

func setupGin() *gin.Engine {
	gin.SetMode("test")
	router := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return router
}

func cleanResponse(responseBytes []byte) string {
	return strings.Replace(string(responseBytes), "\n", "", -1)
}

func TestRespWithError_WithAftError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	router := setupGin()
	router.GET("/user", func(c *gin.Context) {
		RespWithError(c, goerror.DefineBadRequest("InvalidRequest", "Username is required"))
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, `{"errors":[],"message":"Username is required","type":"InvalidRequest"}`, cleanResponse(res.Body.Bytes()))
}

func TestRespWithError_WithDefaultError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	router := setupGin()
	router.GET("/user", func(c *gin.Context) {
		RespWithError(c, errors.New("default error"))
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
	require.Equal(t, `{"message":"default error","type":"UnknownType"}`, cleanResponse(res.Body.Bytes()))
}

func newRequestWithBody(t *testing.T, jsonData interface{}) *http.Request {
	data, err := json.Marshal(jsonData)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
	require.NoError(t, err)

	return req
}

func TestRespValidateError(t *testing.T) {
	expectErrorJson := `{"errors":[{"fieldName":"age","reason":"lte","value":"130"}],"message":"invalid request","type":"InvalidRequest"}`

	reqData := struct {
		UserID string `json:"userID"`
		Name   string `json:"name"`
		Age    int    `json:"age"`
	}{
		UserID: "data1",
		Name:   "tester",
		Age:    500,
	}

	emptyStruct := struct {
		UserID string `json:"userId"`
		Name   string `json:"name" binding:"required"`
		Age    int    `json:"age" form:"age" binding:"gte=0,lte=130"`
	}{}

	req := newRequestWithBody(t, reqData)
	res := httptest.NewRecorder()

	router := setupGin()
	router.POST("/user", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&emptyStruct); err != nil {
			RespValidateError(c, err)
		}
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, expectErrorJson, cleanResponse(res.Body.Bytes()))
}

func TestRespWithErrorReasons(t *testing.T) {
	expectErrorJson := `{"errors":[{"fieldName":"username","reason":"username already exist"},{"fieldName":"phone","reason":"phone number already exist","value":"0598881111"}],"message":"user is already exist","type":"UserExist"}`

	req := newRequestWithBody(t, nil)
	res := httptest.NewRecorder()

	router := setupGin()
	router.POST("/user", func(c *gin.Context) {
		e := goerror.DefineBadRequest("UserExist", "user is already exist")
		e.AddReason("username", "username already exist", nil)
		e.AddReason("phone", "phone number already exist", "0598881111")

		RespWithError(c, e)
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, expectErrorJson, cleanResponse(res.Body.Bytes()))
}
