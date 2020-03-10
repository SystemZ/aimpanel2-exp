package model

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type User struct {
	Base

	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
	//TODO: plan_id
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
