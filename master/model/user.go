package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type User struct {
	ID uuid.UUID `json:"id" gorm:"primary_key;type:varchar(36)"`

	Username string `gorm:"column:username" json:"username"`

	PasswordHash string `gorm:"column:password_hash" json:"-"`

	Email string `gorm:"column:email" json:"email"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

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

//func (u *User) GetHost(db *gorm.DB, hostId string) *Host {
//	var host Host
//	if db.Where("id = ? and user_id = ?", hostId, u.ID).First(&host).RecordNotFound() {
//		return nil
//	}
//	return &host
//}
//
//func (u *User) GetHosts(db *gorm.DB, hostId string) *[]Host {
//	var hosts []Host
//
//	if db.Where("user_id = ?", u.ID).Find(&hosts).RecordNotFound() {
//		return nil
//	}
//	return &hosts
//}
