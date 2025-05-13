package pkg

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testJwtSecret = "test-secret-for-jwt-pkg"

func TestGenToken(t *testing.T) {
	originalSecret, secretIsSet := os.LookupEnv("JWT_SECRET")
	err := os.Setenv("JWT_SECRET", testJwtSecret)
	require.NoError(t, err)
	defer func() {
		if secretIsSet {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	userID := "testUserForGen"
	tokenString, err := GenToken(userID)

	require.NoError(t, err, "GenToken should not produce an error")
	require.NotEmpty(t, tokenString, "Generated token string should not be empty")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJwtSecret), nil
	})

	require.NoError(t, err, "Parsing generated token should not produce an error")
	assert.True(t, token.Valid, "Generated token should be valid")
	assert.Equal(t, userID, claims.Sub, "Subject in claims should match userID")
	assert.WithinDuration(t, time.Now().Add(10*time.Minute), claims.ExpiresAt.Time, 15*time.Second, "Expiration time should be around 10 minutes from now")
}

func TestParseToken(t *testing.T) {
	originalSecret, secretIsSet := os.LookupEnv("JWT_SECRET")
	err := os.Setenv("JWT_SECRET", testJwtSecret)
	require.NoError(t, err)
	defer func() {
		if secretIsSet {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	validUserID := "userToParse"

	generateTestToken := func(userID string, expirationTime time.Time, secret []byte) string {
		claimsMap := jwt.MapClaims{
			"sub": userID,
			"exp": expirationTime.Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsMap)
		signedToken, signErr := token.SignedString(secret)
		require.NoError(t, signErr)
		return signedToken
	}

	tests := []struct {
		name             string
		tokenStringFn    func() string
		setupEnvFn       func()
		cleanupEnvFn     func()
		expectedSub      string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "Valid token",
			tokenStringFn: func() string {
				return generateTestToken(validUserID, time.Now().Add(5*time.Minute), []byte(testJwtSecret))
			},
			expectedSub: validUserID,
			expectError: false,
		},
		{
			name: "Expired token",
			tokenStringFn: func() string {
				return generateTestToken(validUserID, time.Now().Add(-5*time.Minute), []byte(testJwtSecret))
			},
			expectError:      true,
			expectedErrorMsg: "JWT token validation failed",
		},
		{
			name: "Invalid signature (tampered token)",
			tokenStringFn: func() string {
				token := generateTestToken(validUserID, time.Now().Add(5*time.Minute), []byte(testJwtSecret))
				return token + "invalid-part"
			},
			expectError:      true,
			expectedErrorMsg: "JWT token validation failed",
		},
		{
			name: "Token signed with different secret",
			tokenStringFn: func() string {
				return generateTestToken(validUserID, time.Now().Add(5*time.Minute), []byte("another-secret-key"))
			},
			expectError:      true,
			expectedErrorMsg: "JWT token validation failed",
		},
		{
			name: "Malformed token string",
			tokenStringFn: func() string {
				return "not.a.valid.jwt.token"
			},
			expectError:      true,
			expectedErrorMsg: "JWT token validation failed",
		},
		{
			name: "No JWT_SECRET environment variable set when ParseToken is called",
			tokenStringFn: func() string {
				return generateTestToken(validUserID, time.Now().Add(5*time.Minute), []byte(testJwtSecret))
			},
			setupEnvFn: func() {
				os.Unsetenv("JWT_SECRET")
			},
			cleanupEnvFn: func() {
				os.Setenv("JWT_SECRET", testJwtSecret)
			},
			expectError:      true,
			expectedErrorMsg: "jwt secret environment variable not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnvFn != nil {
				tt.setupEnvFn()
			}
			if tt.cleanupEnvFn != nil {
				defer tt.cleanupEnvFn()
			}

			tokenStr := tt.tokenStringFn()
			claims, parseErr := ParseToken(tokenStr)

			if tt.expectError {
				require.Error(t, parseErr, "Expected an error for test case: %s", tt.name)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, parseErr.Error(), tt.expectedErrorMsg, "Error message mismatch for test case: %s", tt.name)
				}
				assert.Nil(t, claims, "Claims should be nil on error for test case: %s", tt.name)
			} else {
				require.NoError(t, parseErr, "Did not expect an error for test case: %s", tt.name)
				require.NotNil(t, claims, "Claims should not be nil for valid token for test case: %s", tt.name)
				assert.Equal(t, tt.expectedSub, claims.Sub, "Subject mismatch for test case: %s", tt.name)
			}
		})
	}
}

func TestGenToken2(t *testing.T) {
	testJwtSecret := "test-jwt-secret"
	validUserID := "user123"

	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", testJwtSecret)
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	tests := []struct {
		name        string
		userID      string
		setupEnvFn  func()
		expectError bool
	}{
		{
			name:        "Generate token with valid user ID",
			userID:      validUserID,
			expectError: false,
		},
		{
			name:        "Generate token with empty user ID",
			userID:      "",
			expectError: false,
		},
		{
			name:   "When JWT_SECRET environment variable is not set",
			userID: validUserID,
			setupEnvFn: func() {
				os.Unsetenv("JWT_SECRET")
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnvFn != nil {
				tt.setupEnvFn()
				defer os.Setenv("JWT_SECRET", testJwtSecret)
			}

			token, err := GenToken(tt.userID)

			if tt.expectError {
				assert.Error(t, err, "An error should occur when generating token")
				assert.Empty(t, token, "Token should be empty when error occurs")
			} else {
				assert.NoError(t, err, "No error should occur when generating token")
				assert.NotEmpty(t, token, "Generated token should not be empty")

				parts := strings.Split(token, ".")
				assert.Equal(t, 3, len(parts), "JWT token should consist of 3 parts")

				if os.Getenv("JWT_SECRET") != "" {
					claims, parseErr := ParseToken(token)
					assert.NoError(t, parseErr, "Generated token should be parsable")
					assert.Equal(t, tt.userID, claims.Sub, "The sub claim of the token should match the user ID")
				}
			}
		})
	}
}

func TestTokenLifecycle(t *testing.T) {
	testJwtSecret := "test-jwt-secret"
	validUserID := "user123"

	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", testJwtSecret)
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	token, err := GenToken(validUserID)
	require.NoError(t, err, "No error should occur when generating token")
	require.NotEmpty(t, token, "Generated token should not be empty")

	claims, parseErr := ParseToken(token)
	require.NoError(t, parseErr, "Generated token should be parsable")
	require.NotNil(t, claims, "Parsed claims should not be nil")

	assert.Equal(t, validUserID, claims.Sub, "The sub claim of the token should match the user ID")
	assert.NotNil(t, claims.ExpiresAt, "Expiration time should be set")
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()), "Expiration time should be after current time")

	expectedExpiry := time.Now().Add(10 * time.Minute)
	timeDiff := claims.ExpiresAt.Time.Sub(expectedExpiry)
	assert.Less(t, timeDiff.Abs(), 5*time.Second, "Expiration time should be set to about 10 minutes later")
}
