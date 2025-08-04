package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   error  `json:"error,omitempty"`
}

type PageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int64  `json:"total"`
	Offset  int    `json:"offset"`
	Limit   int    `json:"limit"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Message: "success",
		Data:    data,
	})
}

func Page(c *gin.Context, data any, total int64, offset, limit int) {
	c.JSON(http.StatusOK, PageResponse{
		Message: "success",
		Data:    data,
		Total:   total,
		Offset:  offset,
		Limit:   limit,
	})
}

func Abort404(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Message: message,
	})
}

func Abort403(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Message: message,
	})
}

func Abort500(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Message: message,
	})
}

func BadRequest(c *gin.Context, err error, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Message: message,
		Error:   err,
	})
}

func Error(c *gin.Context, err error, message string) {

	// error 类型为『数据库未找到内容』
	if err == gorm.ErrRecordNotFound {
		Abort404(c, message)
		return
	}

	c.JSON(http.StatusOK, Response{
		Message: message,
		Error:   err,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
	})
}
