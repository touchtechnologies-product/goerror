## RespWithError

wrap error from goerror to json object response

```go
	router := gin.New()
	router.GET("/user", func(c *gin.Context) {
		RespWithError(c, goerror.DefineBadRequest("InvalidRequest", "Username is required"))
	})
    
    ...
```



## RespValidateError

wrap error from gin validator to json object response

```go
    router.POST("/loginJSON", func(c *gin.Context) {
        var json Login
        if err := c.ShouldBindJSON(&json); err != nil {
            ginresp.RespValidateError(c, err)
            return
        }
    
        ...
```


example object response (validate error)
```json
{
    "errors": [
        {"fieldName": "Age", "reason": "lte", "value": "130"},
        {"fieldName": "Name", "reason": "required", "value": ""}
    ],
    "message": "invalid request",
    "type": "InvalidRequest"
}
```