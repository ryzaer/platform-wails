package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CARA PAKAI
// GenerateAccessToken()
// GenerateAccessToken(24 * time.Hour)
// GenerateAccessToken(24*time.Hour, "ABC123")
// GenerateAccessToken(map[string]any{
// 	"user": "USR001",
// })
// GenerateAccessToken(
// 	30*24*time.Hour,
// 	map[string]any{
// 		"user": "USR001",
// 		"role": "admin",
// 	},
// )

var accessSecret = []byte("CHANGE_ME_TO_RANDOM_SECRET")

func GenerateAccessToken(args ...any) (string, error) {

	expired := 30 * 24 * time.Hour // default

	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(expired).Unix(),
	}

	for _, arg := range args {

		switch v := arg.(type) {

		case time.Duration:
			expired = v
			claims["exp"] = now.Add(expired).Unix()

		case string:
			claims["code"] = v

		case map[string]any:
			for k, val := range v {
				claims[k] = val
			}
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func ReadAccessToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return accessSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
