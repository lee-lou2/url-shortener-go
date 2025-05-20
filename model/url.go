package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"url-shortener-go/config"

	"gorm.io/gorm"
)

// Url is a model struct that stores shortened URL information.
// This struct supports URL redirection, deep links, open graph metadata, and webhook functionality.
type Url struct {
	gorm.Model                // Includes GORM default fields (ID, CreatedAt, UpdatedAt, DeletedAt)
	RandomKey          string `json:"random_key" gorm:"not null;type:varchar(2)"`                // 2-character string storing the random part of the shortened URL
	IOSDeepLink        string `json:"ios_deep_link" gorm:"null;type:text"`                       // Deep link URL connecting to iOS app
	IOSFallbackUrl     string `json:"ios_fallback_url" gorm:"null;type:text"`                    // Fallback URL to use when iOS deep link fails
	AndroidDeepLink    string `json:"android_deep_link" gorm:"null;type:text"`                   // Deep link URL connecting to Android app
	AndroidFallbackUrl string `json:"android_fallback_url" gorm:"null;type:text"`                // Fallback URL to use when Android deep link fails
	DefaultFallbackUrl string `json:"default_fallback_url" gorm:"not null;type:text"`            // Default redirection URL (required field)
	HashedValue        string `json:"hashed_value" gorm:"not null;index;type:text"`              // Hash value to prevent URL duplication
	WebhookUrl         string `json:"webhook_url" gorm:"null;type:text"`                         // Webhook URL to be called when the URL is accessed
	OGTitle            string `json:"og_title" gorm:"null;type:varchar(255)"`                    // Open Graph title
	OGDescription      string `json:"og_description" gorm:"null;type:text"`                      // Open Graph description
	OGImageUrl         string `json:"og_image_url" gorm:"null;type:text"`                        // Open Graph image URL
	IsActive           bool   `json:"is_active" gorm:"default:true;index;not null;type:boolean"` // URL activation status
}

// TableName specifies the table name to be used by GORM.
//
// This function maps the Url model to the 'urls' table.
//
// Returns:
//   - string: table name "urls"
func (u *Url) TableName() string {
	return "urls"
}

// SendWebHook sends access information to the configured webhook URL when a shortened URL is accessed.
//
// Parameters:
//   - shortKey: The key of the accessed shortened URL
//   - userAgent: The User-Agent string of the accessing user
//
// Returns:
//   - error: Error that occurred during webhook transmission, nil on success
//
// Operation:
//  1. Immediately returns nil if webhook URL is not set
//  2. Converts shortKey and userAgent to JSON format
//  3. Sends data to webhook URL via HTTP POST request
//  4. Returns an error if response status code is not 200
func (u *Url) SendWebHook(shortKey, userAgent string) error {
	if u.WebhookUrl == "" {
		return nil
	}

	body := map[string]string{
		"short_key":  shortKey,
		"user_agent": userAgent,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}

	resp, err := http.Post(u.WebhookUrl, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("webhook failed with status code %d", resp.StatusCode)
	}

	return nil
}

// init function runs automatically when the package is loaded,
// performing database migration to create the table corresponding to the Url model.
func init() {
	if config.GetEnv("RUN_MIGRATIONS", "true") == "true" {
		db := config.GetDB()
		if db != nil {
			_ = db.AutoMigrate(&Url{})
		}
	}
}
