package ginresp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/devit-tel/goerror"
)

func cleanResponse(responseBytes []byte) string {
	return strings.Replace(string(responseBytes), "\n", "", -1)
}

func TestRespWithError_WithAftError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	gin.SetMode("test")
	router := gin.New()
	router.GET("/user", func(c *gin.Context) {
		RespWithError(c, goerror.DefineBadRequest("InvalidRequest", "Username is required"))
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, `{"message":"Username is required","type":"InvalidRequest"}`, cleanResponse(res.Body.Bytes()))
}

func TestRespWithError_WithDefaultError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/user", nil)
	res := httptest.NewRecorder()

	gin.SetMode("test")
	router := gin.New()
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
	expectErrorJson := `{"errors":[{"fieldName":"Age","reason":"lte","value":"130"}],"message":"invalid request","type":"InvalidRequest"}`

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
		UserID string `json:"userID"`
		Name   string `json:"name" binding:"required"`
		Age    int    `json:"age" binding:"gte=0,lte=130"`
	}{}

	req := newRequestWithBody(t, reqData)
	res := httptest.NewRecorder()

	gin.SetMode("test")
	router := gin.New()
	router.POST("/user", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&emptyStruct); err != nil {
			RespValidateError(c, err)
		}
	})

	router.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
	require.Equal(t, expectErrorJson, cleanResponse(res.Body.Bytes()))
}
