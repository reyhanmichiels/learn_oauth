package rest

import (
	user_handler "learn_oauth/app/user/handler"
	"learn_oauth/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	gin *gin.Engine
}

func NewRest(gin *gin.Engine) *Rest {
	return &Rest{
		gin: gin,
	}
}

func (r *Rest) HealthCheck() {
	r.gin.GET("/api/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	})
}

func (r *Rest) RouteGoogleAuth(userHandler *user_handler.UserHandler) {
	r.gin.GET("/auth/google/login", userHandler.OauthGoogleLoginHandler)
	r.gin.GET("/auth/google/callback", userHandler.OauthGoogleCallbackHandler)
}

func (r *Rest) RouteUser(userHandler *user_handler.UserHandler) {
	r.gin.GET("/api/v1/user", middleware.AuthJWT, userHandler.GetLoginUser)
}

func (r *Rest) Run() {
	r.gin.Run()
}
