package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

//var jwtSecret = []byte("acesacesaces")

func CreateToken(userId int) (string, int64, error) {
	exp := time.Now().Add(time.Hour * 72).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"userId": userId,
		"exp":    exp, // 3 days
	})

	tokenString, err := token.SignedString(EcdsaPrivateKey)
	return tokenString, exp, err
}

func GetInfoFromToken(tokenString string) (int, int64, error) {
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
			return int(userId), int64(exp), nil // 返回token中的userID和exp
		}
		return -1, -1, fmt.Errorf("userID or exp not found in token")
	}
	return -1, -1, err
}
