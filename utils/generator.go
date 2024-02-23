package utils

import (
	"fmt"
	math "math/rand"
	"time"

	"github.com/Pallinder/go-randomdata"
)

func GenerateRandomUsername() string {
	word := randomdata.SillyName() // Generate a random word

	for len(word) < 8 || len(word) > 12 {
		word = randomdata.SillyName()
	}
	// Generate a random number (between 100 and 999)
	number := math.Intn(900) + 100

	// Concatenate the word and number
	return fmt.Sprintf("%s%d", word, number)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	math.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[math.Intn(len(charset))]
	}
	return string(b)
}
