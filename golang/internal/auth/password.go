package auth

import (
	"github.com/alexedwards/argon2id"
	"fmt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("Error parsing hash: %v", err)
		
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error){
	return argon2id.ComparePasswordAndHash(password, hash)
}