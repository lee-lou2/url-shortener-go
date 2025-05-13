package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is a struct that defines the claim information to be included in the JWT token.
// The Sub field represents the subject of the token (typically the user ID).
// RegisteredClaims includes standard JWT claims (issue time, expiration time, etc.).
type Claims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

// GenToken generates a JWT token based on the user ID.
//
// Parameters:
//   - userID: User identifier to be included in the token
//
// Returns:
//   - string: Signed JWT token string
//   - error: Error that occurred during token generation
//
// The token is signed with the HS256 algorithm and has an expiration time of 10 minutes.
// It uses the JWT_SECRET environment variable to sign the token.
func GenToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ParseToken parses and validates a JWT token string.
//
// Parameters:
//   - tokenString: JWT token string to validate
//
// Returns:
//   - *Claims: Pointer to a struct containing the parsed claim information
//   - error: Error that occurred during token parsing or validation
//
// This function performs the following validations:
//   - Checks if the JWT_SECRET environment variable is set
//   - Verifies that the token's signing method matches expectations
//   - Checks if the token has not expired
//   - Verifies that the token is valid
//
// If an error occurs, it returns nil Claims with a specific error message.
func ParseToken(tokenString string) (*Claims, error) {
	var claims Claims

	// Check JWT secret key
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		return nil, fmt.Errorf("jwt secret environment variable not set")
	}

	// Parse and validate token
	if token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	}); err != nil && err == jwt.ErrTokenExpired {
		// Token has expired
		return nil, fmt.Errorf("JWT token has expired")
	} else if err != nil {
		// Token validation failed
		return nil, fmt.Errorf("JWT token validation failed")
	} else if !token.Valid {
		// Token is invalid
		return nil, fmt.Errorf("JWT token is invalid")
	} else if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		// Token has expired
		return nil, fmt.Errorf("JWT token has expired")
	} else {
		return &claims, nil
	}
}
