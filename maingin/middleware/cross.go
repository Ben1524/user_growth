package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var AllowOrigin = map[string]bool{
	"http://a.site.com": true, // 允许的跨域源
	"http://b.site.com": true,
	"http://web.com":    true,
	"http://12.0.0.1":   true,
	"http://localhost":  true, // 允许的跨域源
}

func CrossMiddleware(c *gin.Context) {
	origin := c.GetHeader("Origin")
	log.Printf("CrossMiddleware Origin=%s\n", origin)
	if AllowOrigin[origin] {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}
