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

func ValidateToken(tokenString string, cfg *config.JWTConfig) (jwt.MapClaims, error) {
	keyFunc := func(secret string) jwt.Keyfunc {
		return func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		}
	}

	// Try parsing with access secret first, then fall back to refresh secret.
	token, err := jwt.Parse(tokenString, keyFunc(cfg.AccessJWTSecret))
	if err != nil {
		token, err = jwt.Parse(tokenString, keyFunc(cfg.RefreshJWTSecret))
		if err != nil {
			return nil, err
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
