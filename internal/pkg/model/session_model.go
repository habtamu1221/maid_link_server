package model

import "github.com/dgrijalva/jwt-go"

type Session struct {
	jwt.StandardClaims
	UserID string
	Role   int
}
