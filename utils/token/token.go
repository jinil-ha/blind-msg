package token

import (
	"math/rand"
	"time"
)

const tokenLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate return new token with length
func Generate(length int) string {
	b := make([]byte, length)
	l := len(tokenLetters)

	for i := range b {
		b[i] = tokenLetters[rand.Intn(l)]
	}
	return string(b)
}
