package api

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Create a validator instance
var validate = validator.New()

// createShortUrlRequest Short URL creation request structure
// Uses validator tags to define validation rules.
//
// Fields:
//   - IOSDeepLink: iOS app deep link URL (optional, must be URL format)
//   - IOSFallbackUrl: URL to redirect when iOS app is not installed (optional, must be URL format)
//   - AndroidDeepLink: Android app deep link URL (optional, must be URL format)
//   - AndroidFallbackUrl: URL to redirect when Android app is not installed (optional, must be URL format)
//   - DefaultFallbackUrl: Default redirect URL (required, must be URL format)
//   - WebhookUrl: Webhook URL (optional, must be URL format)
//   - OGTitle: Open Graph title (optional)
//   - OGDescription: Open Graph description (optional)
//   - OGImageUrl: Open Graph image URL (optional, must be URL format)
type createShortUrlRequest struct {
	IOSDeepLink        string `json:"iosDeepLink" validate:"omitempty,url"`
	IOSFallbackUrl     string `json:"iosFallbackUrl" validate:"omitempty,url"`
	AndroidDeepLink    string `json:"androidDeepLink" validate:"omitempty,url"`
	AndroidFallbackUrl string `json:"androidFallbackUrl" validate:"omitempty,url"`
	DefaultFallbackUrl string `json:"defaultFallbackUrl" validate:"required,url"`
	WebhookUrl         string `json:"webhookUrl" validate:"omitempty,url"`
	OGTitle            string `json:"ogTitle" validate:"omitempty,max=255"`       // Example: Added maximum length limit
	OGDescription      string `json:"ogDescription" validate:"omitempty,max=500"` // Example: Added maximum length limit
	OGImageUrl         string `json:"ogImageUrl" validate:"omitempty,url"`
}

// ValidateStruct function validates the passed structure.
// Fiber handlers can call this function to validate request data.
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			fieldError := validationErrors[0]
			return errors.New("Validation failed on field '" + fieldError.Field() + "' with tag '" + fieldError.Tag() + "'")
		}
		return err
	}
	return nil
}

// validateShortKey Short URL key validation function (existing code)
// This function is not directly used in the createShortUrlRequest structure,
// but if there was a shortKey field in the structure, validator tags could be used in a similar way.
// Example: `validate:"required,min=3,alphanum"`
func validateShortKey(shortKey string) error {
	if len(shortKey) < 3 {
		return errors.New("short_key must be at least 3 characters long")
	}
	matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", shortKey)
	if err != nil || !matched {
		return errors.New("short_key must contain only English letters and numbers")
	}
	return nil
}
