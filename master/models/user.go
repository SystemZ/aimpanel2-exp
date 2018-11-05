package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
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

func (u *User) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		log.Println(err)
	}
	return string(bytes)
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"uid":      u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString, err
}
