package auth

import (
	"github.com/jinzhu/gorm"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

// URL data model
type User struct {
	gorm.Model
	Email string `sql:"type:VARCHAR(255) CHARACTER SET utf8 COLLATE utf8_bin;unique"`
	Name string
	Admin int `gorm:"default:0"`
	Token string
	Password string
}

func (u *User) IsAdmin() bool{
	return u.Admin == 1
}

func (u *User) CreateApiToken() (tokenString string, err error){
	currTime := time.Now().Unix()
	tk := Token{
		Email:  u.Email,
		StandardClaims: &jwt.StandardClaims{
			IssuedAt: currTime,
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err = token.SignedString([]byte("secret"))
	return
}