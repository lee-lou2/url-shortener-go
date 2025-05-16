package pkg

import (
	"github.com/jxskiss/base62"
)

// SplitShortKey splits a short URL key into an ID and a random key.
//
// This function separates the first character and the last character of the short URL key as the random key,
// and the middle part as the base62 encoded ID.
//
// Parameters:
//   - shortKey: The short URL key string to split
//
// Returns:
//   - uint64: The decoded ID value
//   - string: The 2-character random key (first character + last character)
//
// Example:
//   - "a123b" -> 123, "ab" (ID: 123, Random key: "ab")
func SplitShortKey(shortKey string) (uint64, string) {
	frontRandomKey := shortKey[:1]
	backRandomKey := shortKey[len(shortKey)-1:]
	randKey := frontRandomKey + backRandomKey
	uniqueKey := shortKey[1 : len(shortKey)-1]
	id, _ := base62.StdEncoding.ParseUint([]byte(uniqueKey))
	return id, randKey
}

// MergeShortKey combines a random key and an ID to create a short URL key.
//
// This function takes a 2-character random key and an ID to generate a short URL key.
// The first character of the random key is placed at the beginning of the result, and the second character at the end.
// The ID is encoded in base62 and placed in the middle.
//
// Parameters:
//   - randKey: A 2-character random key string
//   - id: The ID value to encode
//
// Returns:
//   - string: The generated short URL key
//
// Example:
//   - "ab", 123 -> "a123b" (Random key: "ab", ID: 123)
func MergeShortKey(randKey string, id uint64) string {
	return randKey[:1] + string(base62.StdEncoding.FormatUint(id)) + randKey[1:]
}
