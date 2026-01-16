package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RayID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rayID := uuid.Must(uuid.NewV7()).String()
		c.Set("ray_id", rayID)
		c.Header("X-Ray-ID", rayID)
		c.Next()
	}
}
