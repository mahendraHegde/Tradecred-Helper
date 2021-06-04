package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mahendraHegde/tradecred-notifier/errors"
)

func getDealsValidator(c *gin.Context) {
	var getDealsQury GetDealsQS
	if err := c.BindQuery(&getDealsQury); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.Set("query", getDealsQury)
	var body Credentials
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.Set("body", body)
	var header Headers
	if err := c.BindHeader(&header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	c.Set("headers", header)
	c.Next()
}

func Authenticator(headers Headers, secret string) error {
	if headers.ApiKey != secret {
		return errors.ApiError{Status: 401, Message: "UnAuthenticated"}
	}
	return nil
}
