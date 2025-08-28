package jwtutil

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID string, username string, accessSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // Access token expires in 15 minutes
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(userID string, refreshSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh token expires in 7 days
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func GenerateTokenPair(userID string, username string, accessSecret, refreshSecret []byte) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(userID, username, accessSecret)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateRefreshToken(userID, refreshSecret)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
