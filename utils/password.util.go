package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	return string(hashBytes), err
}

func MatchHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
