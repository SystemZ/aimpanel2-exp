package model

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"1238206236281802752"`

	Username     string `bson:"username" json:"username"`
	PasswordHash string `bson:"password_hash" json:"password_hash"`
	Email        string `bson:"email" json:"email"`
	//TODO: plan_id
}

func (u *User) GetCollectionName() string {
	return userCollection
}

func (u *User) GetID() primitive.ObjectID {
	return u.ID
}

func (u *User) SetID(id primitive.ObjectID) {
	u.ID = id
}

func (u *User) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		log.Println(err)
	}
	return string(bytes)
}

func (u *User) IsPasswordOk(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"uid":      u.ID.Hex(),
		"username": u.Username,
		"email":    u.Email,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString, err
}

func CheckIfUserExist(username string) bool {
	count, err := DB.Collection(userCollection).CountDocuments(context.TODO(),
		bson.D{{Key: "username", Value: username}})
	if err != nil {
		logrus.Error(err)
	}

	if count > 0 {
		return true
	}

	return false
}

func GetUserById(id primitive.ObjectID) (*User, error) {
	var user User

	err := DB.Collection(userCollection).FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User

	err := DB.Collection(userCollection).FindOne(context.TODO(), bson.D{{Key: "username", Value: username}}).
		Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
