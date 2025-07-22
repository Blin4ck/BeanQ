package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret   []byte
	tokenTTL time.Duration
}

func NewJWTService(secret string, tokenTTL time.Duration) *JWTService {
	return &JWTService{
		secret:   []byte(secret),
		tokenTTL: tokenTTL,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWTToken создает JWT токен для пользователя с ролью
func (j *JWTService) GenerateJWTToken(userID string, role string) (string, error) {
	expirationTime := time.Now().Add(j.tokenTTL)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ParseJWTToken парсит и валидирует JWT токен
func (j *JWTService) ParseJWTToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
