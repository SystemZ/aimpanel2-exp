package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// User represents the user for this application
//
// swagger:model
type User struct {
	// ID of the user
	//
	// required: true
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)"`

	// Username of the user
	//
	// required: true
	Username string `gorm:"column:username"`

	// Password of the user
	//
	// required: true
	PasswordHash string `gorm:"column:password_hash"`

	// Email of the user
	//
	// required: true
	Email string `gorm:"column:email"`

	// Created at timestamp
	CreatedAt time.Time

	// Updated at timestamp
	UpdatedAt time.Time

	// Deleted at timestamp
	DeletedAt *time.Time

	//TODO: plan_id
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuidGen, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	scope.SetColumn("ID", uuidGen)
	return nil
}
