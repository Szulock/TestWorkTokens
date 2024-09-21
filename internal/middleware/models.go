package middleware

import (
	"github.com/dgrijalva/jwt-go"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.StandardClaims
}
