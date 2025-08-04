package response

import (
	"evaframe/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
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
	logger.L().LogIf(err)
	c.JSON(http.StatusBadRequest, Response{
		Message: message,
		Error:   err.Error(),
	})
}

func Error(c *gin.Context, err error, message string) {
	logger.L().LogIf(err)

	// error 类型为『数据库未找到内容』
	if err == gorm.ErrRecordNotFound {
		Abort404(c, message)
		return
	}

	c.JSON(http.StatusOK, Response{
		Message: message,
		Error:   err.Error(),
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Message: message,
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Message: message,
	})
}
