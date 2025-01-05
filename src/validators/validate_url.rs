use regex::Regex;

/// Validate Email
/// TODO: Email address validation using regex is needed
pub fn validate_email(email: &str) -> Result<bool, String> {
    if email.is_empty() {
        return Err("Email is missing.".to_string());
    }
    Ok(true)
}

/// Validate URL
/// Only basic URL structure is checked
pub fn validate_url(url: &str) -> Result<bool, String> {
    let url_regex = Regex::new(r"^https?://[^\s/$.?#].[^\s]*$").unwrap();
    if !url_regex.is_match(url) {
        return Err("The URL format is invalid.".to_string());
    }
    Ok(true)
}

/// Validate WebHook URL
pub fn validate_webhook_url(url: &str) -> Result<bool, String> {
    if !url.is_empty() {
        return validate_url(url);
    }
    Ok(true)
}

/// Validate FallBack URL
pub fn validate_fallback_url(url: &str) -> Result<bool, String> {
    if !url.is_empty() {
        return validate_url(url);
    }
    Ok(true)
}
