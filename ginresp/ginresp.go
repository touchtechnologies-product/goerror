package ginresp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/devit-tel/goerror"
)

type ErrorValidate struct {
	FieldName string      `json:"fieldName"`
	Reason    string      `json:"reason"`
	Value     interface{} `json:"value"`
}

func RespWithError(c *gin.Context, err error) {
	if e, ok := err.(*goerror.GoError); ok {
		c.JSON(e.Status, gin.H{
			"type":    e.Code,
			"message": e.Msg + e.ExtendMsg,
		})

		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"type":    "UnknownType",
		"message": err.Error(),
	})
}

func RespValidateError(c *gin.Context, err error) {
	errValidates := make([]*ErrorValidate, 0)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, errValidate := range errs {
			errValidates = append(errValidates, &ErrorValidate{
				FieldName: errValidate.Field(),
				Reason:    errValidate.ActualTag(),
				Value:     errValidate.Param(),
			})
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"type":    "InvalidRequest",
		"message": "invalid request",
		"errors":  errValidates,
	})
}
