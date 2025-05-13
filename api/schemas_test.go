package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct_createShortUrlRequest(t *testing.T) {
	tests := []struct {
		name     string
		req      createShortUrlRequest
		wantErr  bool
		errField string
	}{
		{
			name:    "Valid request - only default fallback URL",
			req:     createShortUrlRequest{DefaultFallbackUrl: "http://example.com"},
			wantErr: false,
		},
		{
			name: "Valid request - all standard URL fields",
			req: createShortUrlRequest{
				IOSFallbackUrl:     "http://ios.store.com/app",
				AndroidFallbackUrl: "http://android.store.com/app",
				DefaultFallbackUrl: "http://example.com/default",
				WebhookUrl:         "https://webhook.site/token",
				OGImageUrl:         "https://example.com/image.png",
				OGTitle:            "Sample Title",
				OGDescription:      "Sample Description",
			},
			wantErr: false,
		},
		{
			name: "Valid request - deep links can be non-http if 'url' tag isn't too strict or if omitted",
			req: createShortUrlRequest{
				IOSDeepLink:        "myapp://product/123",
				AndroidDeepLink:    "androidapp://deeplink",
				DefaultFallbackUrl: "https://fallback.com",
			},
			wantErr:  false,
			errField: "",
		},
		{
			name:     "Invalid - DefaultFallbackUrl missing",
			req:      createShortUrlRequest{},
			wantErr:  true,
			errField: "DefaultFallbackUrl",
		},
		{
			name:     "Invalid - DefaultFallbackUrl not a URL",
			req:      createShortUrlRequest{DefaultFallbackUrl: "not-a-url-string"},
			wantErr:  true,
			errField: "DefaultFallbackUrl",
		},
		{
			name: "Invalid - WebhookUrl not a URL",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "http://example.com",
				WebhookUrl:         "invalid-url",
			},
			wantErr:  true,
			errField: "WebhookUrl",
		},
		{
			name: "Invalid - OGTitle too long (max 255)",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "http://example.com",
				OGTitle:            string(make([]byte, 256)),
			},
			wantErr:  true,
			errField: "OGTitle",
		},
		{
			name: "Valid - Deep links omitted (passes 'omitempty')",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "http://valid.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(&tt.req)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error for validation.")
				if tt.errField != "" && err != nil {
					assert.Contains(t, err.Error(), "'"+tt.errField+"'", "Error message should mention the failing field.")
				}
			} else {
				assert.NoError(t, err, "Expected no error for valid struct.")
			}
		})
	}
}

func TestValidateShortKey(t *testing.T) {
	tests := []struct {
		name     string
		shortKey string
		wantErr  bool
		errMsg   string
	}{
		{name: "Valid key - alphanumeric", shortKey: "ValidKey123", wantErr: false},
		{name: "Valid key - all numbers", shortKey: "789012", wantErr: false},
		{name: "Valid key - all letters", shortKey: "OnlyLetters", wantErr: false},
		{name: "Invalid - too short (2 chars)", shortKey: "a1", wantErr: true, errMsg: "short_key must be at least 3 characters long"},
		{name: "Invalid - empty string", shortKey: "", wantErr: true, errMsg: "short_key must be at least 3 characters long"},
		{name: "Invalid - contains hyphen", shortKey: "invalid-key", wantErr: true, errMsg: "short_key must contain only English letters and numbers"},
		{name: "Invalid - contains space", shortKey: "key with space", wantErr: true, errMsg: "short_key must contain only English letters and numbers"},
		{name: "Invalid - contains symbol", shortKey: "key_symbol", wantErr: true, errMsg: "short_key must contain only English letters and numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateShortKey(tt.shortKey)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error for short key validation.")
				if tt.errMsg != "" {
					assert.EqualError(t, err, tt.errMsg, "Error message mismatch.")
				}
			} else {
				assert.NoError(t, err, "Expected no error for valid short key.")
			}
		})
	}
}

func TestValidateStructWithCustomTypes(t *testing.T) {
	type customRequest struct {
		Email     string `json:"email" validate:"required,email"`
		Age       int    `json:"age" validate:"required,gte=18,lte=100"`
		Username  string `json:"username" validate:"required,min=3,max=20,alphanum"`
		Password  string `json:"password" validate:"required,min=8"`
		URL       string `json:"url" validate:"omitempty,url"`
		IPAddress string `json:"ipAddress" validate:"omitempty,ip"`
	}

	tests := []struct {
		name     string
		req      customRequest
		wantErr  bool
		errField string
	}{
		{
			name: "All fields valid",
			req: customRequest{
				Email:     "test@example.com",
				Age:       25,
				Username:  "testuser123",
				Password:  "password123",
				URL:       "https://example.com",
				IPAddress: "192.168.1.1",
			},
			wantErr: false,
		},
		{
			name: "Email format error",
			req: customRequest{
				Email:    "invalid-email",
				Age:      25,
				Username: "testuser123",
				Password: "password123",
			},
			wantErr:  true,
			errField: "Email",
		},
		{
			name: "Age range error (too low)",
			req: customRequest{
				Email:    "test@example.com",
				Age:      15,
				Username: "testuser123",
				Password: "password123",
			},
			wantErr:  true,
			errField: "Age",
		},
		{
			name: "Age range error (too high)",
			req: customRequest{
				Email:    "test@example.com",
				Age:      101,
				Username: "testuser123",
				Password: "password123",
			},
			wantErr:  true,
			errField: "Age",
		},
		{
			name: "Username length error (too short)",
			req: customRequest{
				Email:    "test@example.com",
				Age:      25,
				Username: "ab",
				Password: "password123",
			},
			wantErr:  true,
			errField: "Username",
		},
		{
			name: "Username format error (contains special characters)",
			req: customRequest{
				Email:    "test@example.com",
				Age:      25,
				Username: "user@name",
				Password: "password123",
			},
			wantErr:  true,
			errField: "Username",
		},
		{
			name: "Password length error",
			req: customRequest{
				Email:    "test@example.com",
				Age:      25,
				Username: "testuser123",
				Password: "pass",
			},
			wantErr:  true,
			errField: "Password",
		},
		{
			name: "URL format error",
			req: customRequest{
				Email:    "test@example.com",
				Age:      25,
				Username: "testuser123",
				Password: "password123",
				URL:      "invalid-url",
			},
			wantErr:  true,
			errField: "URL",
		},
		{
			name: "IP address format error",
			req: customRequest{
				Email:     "test@example.com",
				Age:       25,
				Username:  "testuser123",
				Password:  "password123",
				URL:       "https://example.com",
				IPAddress: "invalid-ip",
			},
			wantErr:  true,
			errField: "IPAddress",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(&tt.req)
			if tt.wantErr {
				assert.Error(t, err, "Validation error expected")
				if tt.errField != "" && err != nil {
					assert.Contains(t, err.Error(), "'"+tt.errField+"'", "Error message should contain the failed field")
				}
			} else {
				assert.NoError(t, err, "No error should occur for valid struct")
			}
		})
	}
}

func TestCreateShortUrlRequestValidation(t *testing.T) {
	tests := []struct {
		name     string
		req      createShortUrlRequest
		wantErr  bool
		errField string
	}{
		{
			name: "All fields valid",
			req: createShortUrlRequest{
				IOSDeepLink:        "https://example.com/ios",
				IOSFallbackUrl:     "https://example.com/ios-fallback",
				AndroidDeepLink:    "https://example.com/android",
				AndroidFallbackUrl: "https://example.com/android-fallback",
				DefaultFallbackUrl: "https://example.com/default",
				WebhookUrl:         "https://example.com/webhook",
				OGTitle:            "Test Title",
				OGDescription:      "Test Description",
				OGImageUrl:         "https://example.com/image.jpg",
			},
			wantErr: false,
		},
		{
			name: "Only required fields entered",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "https://example.com/default",
			},
			wantErr: false,
		},
		{
			name: "DefaultFallbackUrl missing",
			req: createShortUrlRequest{
				IOSDeepLink:     "https://example.com/ios",
				AndroidDeepLink: "https://example.com/android",
			},
			wantErr:  true,
			errField: "DefaultFallbackUrl",
		},
		{
			name: "Invalid URL format - IOSDeepLink",
			req: createShortUrlRequest{
				IOSDeepLink:        "invalid-url",
				DefaultFallbackUrl: "https://example.com/default",
			},
			wantErr:  true,
			errField: "IOSDeepLink",
		},
		{
			name: "Invalid URL format - DefaultFallbackUrl",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "invalid-url",
			},
			wantErr:  true,
			errField: "DefaultFallbackUrl",
		},
		{
			name: "OGTitle exceeds maximum length",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "https://example.com/default",
				OGTitle:            string(make([]byte, 256)),
			},
			wantErr:  true,
			errField: "OGTitle",
		},
		{
			name: "OGDescription exceeds maximum length",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "https://example.com/default",
				OGDescription:      string(make([]byte, 501)),
			},
			wantErr:  true,
			errField: "OGDescription",
		},
		{
			name: "Invalid URL format - OGImageUrl",
			req: createShortUrlRequest{
				DefaultFallbackUrl: "https://example.com/default",
				OGImageUrl:         "invalid-image-url",
			},
			wantErr:  true,
			errField: "OGImageUrl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(&tt.req)
			if tt.wantErr {
				assert.Error(t, err, "Validation error expected")
				if tt.errField != "" && err != nil {
					assert.Contains(t, err.Error(), "'"+tt.errField+"'", "Error message should contain the failed field")
				}
			} else {
				assert.NoError(t, err, "No error should occur for valid struct")
			}
		})
	}
}
