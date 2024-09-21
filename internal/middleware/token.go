package middleware

import (
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("HiMedodsHopeYouLikeThisCode")

func generateTokens(userID, ip string, db *sql.DB) (string, string, error) {

	accessClaims := Claims{
		UserID: userID,
		IP:     ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
		},
	}
	Logger.Info("Создаю access токен")
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	signedAccessToken, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	Logger.Info("Создаю и хэширую refresh токен")
	rawRefreshToken := base64.StdEncoding.EncodeToString([]byte(userID + ":" + ip))
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(rawRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	Logger.Info("Сохраняю refresh токен")
	_, err = db.Exec("INSERT INTO refresh_tokens (user_id, token, ip_address) VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET token = $2, ip_address = $3", userID, hashedRefreshToken, ip)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, rawRefreshToken, nil
}
