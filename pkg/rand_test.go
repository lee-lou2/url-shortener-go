package pkg

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenRandStr(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "Length 0", length: 0},
		{name: "Length 1", length: 1},
		{name: "Length 2 (used in app)", length: 2},
		{name: "Length 10", length: 10},
		{name: "Length 32", length: 32},
	}

	allowedCharsRegex := regexp.MustCompile("^[a-zA-Z0-9]*$")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenRandStr(tt.length)
			assert.Len(t, got, tt.length, "Generated string should have the specified length.")

			if tt.length > 0 {
				assert.True(t, allowedCharsRegex.MatchString(got), "Generated string should only contain alphanumeric characters.")
			} else {
				assert.Empty(t, got, "Generated string should be empty for length 0.")
			}
		})
	}

	if testing.Short() {
		t.Skip("Skipping pseudo-randomness check in short mode.")
	}

	const testLengthForRandomness = 8
	const numberOfSamples = 5
	generated := make(map[string]bool)

	for i := 0; i < numberOfSamples; i++ {
		s := GenRandStr(testLengthForRandomness)
		assert.Len(t, s, testLengthForRandomness)
		generated[s] = true
	}

	assert.GreaterOrEqual(t, len(generated), numberOfSamples-1, "Generated strings of the same length should generally be unique across multiple calls.")
}
