package middleware

import (
	"learn_oauth/domain"
	"learn_oauth/infrastructure"
	"learn_oauth/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthJWT(c *gin.Context) {
	userId, tokenExpireTime, err := util.ParsesAndValidateJWT(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to parse and validate token",
			"error":   err.Error(),
		})
	}

	if float64(time.Now().Unix()) > tokenExpireTime {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "your session has been expired",
			"error":   nil,
		})
		return
	}

	var user domain.User
	err = infrastructure.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "failed to get user",
			"error":   err,
		})
		return
	}

	c.Set("user", user)
	c.Next()
}
