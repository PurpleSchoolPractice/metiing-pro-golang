package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email  string
	UserID uint
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type JWT struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret:          secret,
		AccessTokenTTL:  time.Minute * 4,
		RefreshTokenTTL: time.Minute * 5,
	}
}

func (j *JWT) generateToken(data JWTData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   data.Email,
		"user_id": data.UserID,
		"exp":     time.Now().Add(j.AccessTokenTTL).Unix(),
	})
	return t.SignedString([]byte(j.Secret))
}

func (j *JWT) generateRefreshToken(data JWTData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": data.Email,
		"exp":   time.Now().Add(j.RefreshTokenTTL).Unix(),
		"type":  "refresh",
	})
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) GenerateTokenPair(data JWTData) (*TokenPair, error) {
	accessToken, err := j.generateToken(data)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateRefreshToken(data)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWT) ParseToken(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil || !t.Valid {
		return false, nil
	}

	claims := t.Claims.(jwt.MapClaims)

	email, ok1 := claims["email"].(string)
	userIDFloat, ok2 := claims["user_id"].(float64) // JWT хранит числа как float64

	if !ok1 || !ok2 {
		return false, nil
	}

	return true, &JWTData{
		Email:  email,
		UserID: uint(userIDFloat),
	}
}

func (j *JWT) ParseRefreshToken(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return false, nil
	}

	claims := t.Claims.(jwt.MapClaims)
	tokenType, ok := claims["type"]
	if !ok || tokenType != "refresh" {
		return false, nil
	}

	email := claims["email"].(string)
	return t.Valid, &JWTData{
		Email: email,
	}
}
