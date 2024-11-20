package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	AuthType string `json:"auth_type"`
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}
