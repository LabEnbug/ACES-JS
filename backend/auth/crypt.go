package auth

import "golang.org/x/crypto/bcrypt"

func MakePassword(password string) string {
	encodedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encodedPassword)
}

func ValidatePassword(password string, encodedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(password))
	return err == nil
}
