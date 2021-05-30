package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS = 200
	// code >= 10000 is error code
	ERROR = 10000
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Resp(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Fail(c *gin.Context, message string) {
	Resp(c, ERROR, map[string]interface{}{}, message)
}

func FailData(c *gin.Context, data interface{}, message string) {
	Resp(c, ERROR, data, message)
}

func FailError(c *gin.Context, err error) {
	Resp(c, ERROR, map[string]interface{}{}, err.Error())
}

// ---------------------------------------------------------------------------------------------------------------------

func Ok(c *gin.Context) {
	Resp(c, SUCCESS, map[string]interface{}{}, "success")
}

func Message(c *gin.Context, message string) {
	Resp(c, SUCCESS, map[string]interface{}{}, message)
}

func Data(c *gin.Context, data interface{}) {
	Resp(c, SUCCESS, data, "")
}

func DataMessage(c *gin.Context, data interface{}, message string) {
	Resp(c, SUCCESS, data, message)
}
