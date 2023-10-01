package funcs

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	secretKey     = "YourSecretKey" // 用于签名和验证JWT的密钥，请确保这是一个复杂的随机字符串
	tokenDuration = 72              // Token有效时间（小时）
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * tokenDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
