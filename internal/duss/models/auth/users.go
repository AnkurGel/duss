package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

// User data model
type User struct {
	gorm.Model
	Email string `sql:"type:VARCHAR(255) CHARACTER SET utf8 COLLATE utf8_bin;unique" json:"email_id"`
	Name string `json:"name"`
	Admin int `gorm:"default:0" json:"is_admin"`
	Token string `json:"token_id"`
	Password string `json:"password,omitempty"`
}

// IsAdmin verifies the admin role of the user
func (u *User) IsAdmin() bool{
	return u.Admin == 1
}

// CreateAPIToken creates token for user
func (u *User) CreateAPIToken() (tokenString string, err error){
	currTime := time.Now().Unix()
	tk := Token{
		Email:  u.Email,
		Name: u.Name,
		StandardClaims: &jwt.StandardClaims{
			IssuedAt: currTime,
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err = token.SignedString([]byte(viper.GetString("JwtSecret")))
	return
}