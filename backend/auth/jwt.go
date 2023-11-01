package auth

import (
	"backend/config"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

//var jwtSecret = []byte("acesacesaces")

func CreateToken(userId uint) (string, int64, error) {
	exp := time.Now().Add(config.TokenExpireTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"userId": userId,
		"exp":    exp, // 3 days
	})

	tokenString, err := token.SignedString(EcdsaPrivateKey)
	return tokenString, exp, err
}

func GetInfoFromToken(tokenString string) (uint, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return EcdsaPublicKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok1 := claims["userId"].(float64)
		exp, ok2 := claims["exp"].(float64)
		if ok1 && ok2 {
			return uint(userId), int64(exp), nil // 返回token中的userID和exp
		}
		return 0, 0, fmt.Errorf("userID or exp not found in token")
	}
	return 0, 0, err
}
