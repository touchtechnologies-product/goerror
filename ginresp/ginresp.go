package ginresp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	"github.com/devit-tel/goerror"
)

func RespWithError(c *gin.Context, err error) {
	if e, ok := err.(*goerror.GoError); ok {
		c.JSON(e.Status, gin.H{
			"type":    e.Code,
			"message": e.Msg + e.ExtendMsg,
			"errors":  e.GetReasons(),
		})

		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"type":    "UnknownType",
		"message": err.Error(),
	})
}

func RespValidateError(c *gin.Context, err error) {
	errValidates := make([]*goerror.Reason, 0)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, errValidate := range errs {
			errValidates = append(errValidates, &goerror.Reason{
				FieldName: errValidate.Namespace(),
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
