package utils

import (
	"errors"
	"time"

	"github.com/Perlishnov/gotrainingproject/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct{
	secret []byte
	exp time.Duration
}

func NewJWTUtil(cfg *config.Config) *JWTUtil  {
	return &JWTUtil{
		secret: []byte(cfg.JWTSecret),
		exp: time.Duration(cfg.JWTExpirationHours) * time.Hour,
	}
}

type Claims struct{
	UserID string `json:"user_id"`
	Email string `json:"email"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTUtil) GenerateToken(userID string, email, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.exp)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error)  {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid{
		return nil, errors.New("Invalid token")
	}
	return claims, nil
}

