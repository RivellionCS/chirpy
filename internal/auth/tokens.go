package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	randKey := make([]byte, 32)
	rand.Read(randKey)
	encodedStr := hex.EncodeToString(randKey)
	return encodedStr
}