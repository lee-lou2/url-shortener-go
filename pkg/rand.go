package pkg

import (
	"math/rand"
)

// GenRandStr generates a random string of the specified length
//
// This function generates a random string containing uppercase and lowercase letters and numbers.
// It is primarily used for generating random keys in URL shortening services.
//
// Parameters:
//   - length: The length of the random string to generate
//
// Returns:
//   - The generated random string
//   - Returns an empty string if length is 0
//
// Notes:
//   - This function uses the math/rand package, so it is not cryptographically secure.
//   - For security-sensitive purposes, it is recommended to use the crypto/rand package.
func GenRandStr(length int) string {
	if length == 0 {
		return ""
	}

	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}
