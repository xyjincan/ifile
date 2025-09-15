package tool

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func passAdminAPI(c *gin.Context) bool {
	passUserAPI(c)
	clientIP := c.Request.RemoteAddr
	if !strings.HasPrefix(clientIP, "127.0.0.1") {
		fmt.Println("denied", clientIP)
		return false
	}
	return true
}

func passUserAPI(c *gin.Context) bool {
	clientIP := c.Request.RemoteAddr
	if strings.HasPrefix(clientIP, "192.168.") {
		fmt.Println("clientIP", clientIP, "内网IP")
	} else {
		fmt.Println("clientIP", clientIP)
	}
	return true
}
