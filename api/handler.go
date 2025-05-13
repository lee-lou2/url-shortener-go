package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
	"url-shortener-go/config"
	"url-shortener-go/model"
	"url-shortener-go/pkg"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// indexHandler Main page handler
// Renders the main page and generates a guest token.
// The guest token is used for authentication when calling APIs.
//
// GET /
// Response:
//   - 200: Main page HTML (index.html)
//   - 500: Error message when JWT generation fails
func indexHandler(c fiber.Ctx) error {
	t, err := pkg.GenToken("guest")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate JWT",
		})
	}
	// Save in encrypted cookie
	c.Cookie(&fiber.Cookie{Name: "token", Value: t})
	return c.Render("index", fiber.Map{"Token": t})
}

// createShortUrlHandler Short URL creation handler
// Validates the input URL information and creates a short URL.
// If the same URL information already exists, returns the existing short URL.
//
// POST /v1/urls
// Request Body:
//   - IOSDeepLink: iOS app deep link (optional)
//   - IOSFallbackUrl: URL to redirect when iOS app is not installed (optional)
//   - AndroidDeepLink: Android app deep link (optional)
//   - AndroidFallbackUrl: URL to redirect when Android app is not installed (optional)
//   - DefaultFallbackUrl: Default redirect URL (required)
//   - WebhookUrl: Webhook URL (optional)
//   - OGTitle: Open Graph title (optional)
//   - OGDescription: Open Graph description (optional)
//   - OGImageUrl: Open Graph image URL (optional)
//
// Response:
//   - 200: Short URL creation success {"message": "URL created successfully", "short_key": "shortKey"}
//   - 400: Invalid request format or validation failure
//   - 500: Internal server error
func createShortUrlHandler(c fiber.Ctx) error {
	var reqBody createShortUrlRequest
	if err := c.Bind().JSON(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// 1. Validation
	if err := ValidateStruct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 2. Generate hash
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s:%s:%s:%s",
		reqBody.IOSDeepLink,
		reqBody.IOSFallbackUrl,
		reqBody.AndroidDeepLink,
		reqBody.AndroidFallbackUrl,
		reqBody.DefaultFallbackUrl)))
	hashedValue := fmt.Sprintf("%x", hasher.Sum(nil))

	// 3. Check if exists
	db := config.GetDB()
	var url model.Url
	if err := db.Where("hashed_value = ?", hashedValue).First(&url).Error; err == nil {
		// If exists
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "URL already exists",
		})
	} else if err != gorm.ErrRecordNotFound {
		// If unknown error occurs
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to check URL",
		})
	}

	// 4. Create if not exists
	randKey := pkg.GenRandStr(2)
	url = model.Url{
		RandomKey:          randKey,
		IOSDeepLink:        reqBody.IOSDeepLink,
		IOSFallbackUrl:     reqBody.IOSFallbackUrl,
		AndroidDeepLink:    reqBody.AndroidDeepLink,
		AndroidFallbackUrl: reqBody.AndroidFallbackUrl,
		DefaultFallbackUrl: reqBody.DefaultFallbackUrl,
		HashedValue:        hashedValue,
		WebhookUrl:         reqBody.WebhookUrl,
		OGTitle:            reqBody.OGTitle,
		OGDescription:      reqBody.OGDescription,
		OGImageUrl:         reqBody.OGImageUrl,
		IsActive:           true,
	}
	if err := db.Create(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create URL",
		})
	}
	shortKey := pkg.MergeShortKey(randKey, uint64(url.ID))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "URL created successfully",
		"short_key": shortKey,
	})
}

// redirectToOriginalHandler Short URL redirect handler
// Takes the short URL key, looks up the original URL information, and renders the redirect page.
// First checks the Redis cache, and if not found, queries the DB.
// If a webhook URL is set, it calls the webhook asynchronously.
//
// GET /:short_key
// Path Parameters:
//   - short_key: Short URL key (required)
//
// Response:
//   - 200: Redirect page HTML (redirect.html)
//   - 400: Invalid short URL key format
//   - 404: Short URL not found
//   - 500: Internal server error
func redirectToOriginalHandler(c fiber.Ctx) error {
	shortKey := c.Params("short_key")

	// 1. Validation
	if err := validateShortKey(shortKey); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 2. Check cache
	var url model.Url
	cache := config.GetCache(c.Context())
	cacheKey := fmt.Sprintf("urls:%s", shortKey)
	if cachedVal, err := cache.Get(c.Context(), cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cachedVal), &url); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to unmarshal URL",
			})
		}

		// Call webhook
		go url.SendWebHook(shortKey, c.GetReqHeaders()["User-Agent"][0])
		return c.Render("redirect", fiber.Map{"Object": url})
	}

	// 3. If not in cache, query DB
	db := config.GetDB()
	id, randKey := pkg.SplitShortKey(shortKey)
	if err := db.Where("id = ?", id).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "URL not found",
		})
	} else if url.RandomKey != randKey {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "URL not found",
		})
	}

	// 4. Save to cache
	if data, err := json.Marshal(url); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to marshal URL",
		})
	} else {
		cache.Set(c.Context(), cacheKey, data, 1*time.Hour)
	}

	// 5. Call webhook
	go url.SendWebHook(shortKey, c.GetReqHeaders()["User-Agent"][0])

	// 6. Redirect to success page
	return c.Render("redirect", fiber.Map{"Object": url})
}
