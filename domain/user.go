package domain

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;notnull;primary_key"`
	Name  string    `json:"name" gorm:"type:varchar(50);notnull"`
	Email string    `json:"email" gorm:"type:varchar(50);notnull"`
}
