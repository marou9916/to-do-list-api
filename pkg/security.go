package pkg

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken() string {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	return hex.EncodeToString(bytes)
}
