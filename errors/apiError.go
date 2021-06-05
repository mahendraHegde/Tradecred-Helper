package errors

import "github.com/gin-gonic/gin"

//ApiError is app error
type ApiError struct {
	Status  int
	Message string
	Meta    interface{} `json:",omitempty"`
}

func (error ApiError) Error() string {
	return error.Message
}

func BuildJsonError(c *gin.Context, message string, status int, meta interface{}) {
	c.JSON(int(status), ApiError{Status: status, Message: message, Meta: meta})
}
