package service

import "github.com/gin-gonic/gin"

// GetUserID return user id
func GetUserID(c *gin.Context) uint64 {
	if c == nil {
		return 0
	}

	// uid must be named the same as the uid in middleware/auth
	if v, exists := c.Get("uid"); exists {
		uid, ok := v.(uint64)
		if !ok {
			return 0
		}

		return uid
	}
	return 0
}
