package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"log"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)"`
	Username string    `gorm:"column:username"`
	Password string    `gorm:"column:password"`
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)
	return nil
}
