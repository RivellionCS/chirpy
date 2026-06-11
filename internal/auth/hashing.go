package auth

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		return "", err
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	check, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("error checking password and hash: %s", err)
		return false, err
	}
	return check, nil
}