package models

import (
	"github.com/gofrs/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)"`
	Username string    `gorm:"column:username"`
	Password string    `gorm:"column:password"`
}
