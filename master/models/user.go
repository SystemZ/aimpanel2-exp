package models

import (
	"github.com/jinzhu/gorm"
	"github.com/gofrs/uuid"
)

type User struct {
	gorm.Model

	ID uuid.UUID `gorm:"primarykey;type:char(36);column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}