package additional_galxe_methods

import (
	"math/rand"
	"time"
)

func GeneratePersonaWidget() string {
	charSet := "abcdefghijklmnopqrstuvwxyz0123456789"
	var output string
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 16; i++ {
		randomIndex := rand.Intn(len(charSet))
		output += string(charSet[randomIndex])
	}

	return output
}
