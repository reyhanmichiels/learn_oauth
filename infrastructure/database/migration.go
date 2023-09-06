package database

import (
	"learn_oauth/domain"
	"learn_oauth/infrastructure"
)

func Migrate() {
	infrastructure.DB.AutoMigrate(&domain.User{})
}