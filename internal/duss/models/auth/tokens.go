package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
)


//Token data model
type Token struct {
	Email  string
	Name string
	*jwt.StandardClaims
}
