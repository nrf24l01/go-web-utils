package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nrf24l01/go-web-utils/config"
)

func GenerateAccessToken(claims jwt.MapClaims, cfg *config.JWTConfig) (string, error) {
	claims["exp"] = time.Now().Add(time.Duration(cfg.AccessTokenExpiryMinutes) * time.Minute).Unix() // Access token expires in configured minutes
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.AccessJWTSecret))
}

func GenerateRefreshToken(claims jwt.MapClaims, cfg *config.JWTConfig) (string, error) {
	claims["exp"] = time.Now().Add(time.Duration(cfg.RefreshTokenExpiryMinutes) * time.Minute).Unix() // Refresh token expires in configured minutes
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.RefreshJWTSecret))
}

func GenerateTokenPair(accessClaims, refreshClaims jwt.MapClaims, cfg *config.JWTConfig) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(accessClaims, cfg)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateRefreshToken(refreshClaims, cfg)
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
