package utils

import (
	"time"

	"iscsi-web-panel/models"

	"github.com/gin-gonic/gin"
)

// Success returns a 200 JSON response
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
	})
}

// Error returns an error JSON response
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, models.APIResponse{
		Code:    code,
		Message: msg,
	})
}

// NowISO returns current time in ISO format
func NowISO() time.Time {
	return time.Now()
}

// FormatBytes formats bytes to human-readable string
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return string(rune(bytes+'0')) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return string([]byte{byte(bytes/div + '0'), '.', byte((bytes%div*10/div) + '0')}) + " " + []string{"KB", "MB", "GB", "TB", "PB"}[exp]
}
