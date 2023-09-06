package main

import (
	user_handler "learn_oauth/app/user/handler"
	user_repo "learn_oauth/app/user/repository"
	user_usecase "learn_oauth/app/user/usecase"
	"learn_oauth/infrastructure"
	"learn_oauth/infrastructure/database"
	"learn_oauth/rest"

	"github.com/gin-gonic/gin"
)

func main() {
	//init env
	infrastructure.LoadEnv()

	//init db
	infrastructure.ConnectDatabase()

	database.Migrate()

	//init repository
	userRepo := user_repo.NewUserRepository(infrastructure.DB)

	//init usecase
	userUsecase := user_usecase.NewUserUsecase(userRepo)

	//init handler
	userHandler := user_handler.NewUserHandler(userUsecase)

	user_handler.NewGoogleOauthConfig()

	//initialize rest
	rest := rest.NewRest(gin.Default())

	//initialize healthcheck
	rest.HealthCheck()

	//initialize google auth
	rest.RouteGoogleAuth(userHandler)

	//initialize user route
	rest.RouteUser(userHandler)

	//run engine
	rest.Run()
}
