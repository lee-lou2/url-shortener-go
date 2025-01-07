use regex::Regex;

/// Validate Email
/// Basic email format validation
pub fn validate_email(email: &str) -> Result<bool, String> {
    if email.is_empty() {
        return Err("Email is missing.".to_string());
    }
    let email_regex = Regex::new(r#"(?i)^[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}$"#).unwrap();
    if !email_regex.is_match(email) {
        return Err("Invalid email format.".to_string());
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_validate_email() {
        // Valid cases
        assert_eq!(validate_email("test@example.com"), Ok(true));
        assert_eq!(
            validate_email("user.name+tag+sorting@example.com"),
            Ok(true)
        );
        assert_eq!(validate_email("user@subdomain.example.com"), Ok(true));
        assert_eq!(validate_email("USER@EXAMPLE.COM"), Ok(true));
        assert_eq!(validate_email("simple@example.co.uk"), Ok(true));

        // Invalid cases
        assert_eq!(validate_email(""), Err("Email is missing.".to_string()));
        assert_eq!(
            validate_email("plainaddress"),
            Err("Invalid email format.".to_string())
        );
        assert_eq!(
            validate_email("@missingusername.com"),
            Err("Invalid email format.".to_string())
        );
        assert_eq!(
            validate_email("username.com"),
            Err("Invalid email format.".to_string())
        );
        assert_eq!(
            validate_email("username@.com"),
            Err("Invalid email format.".to_string())
        );
        assert_eq!(
            validate_email("username@com"),
            Err("Invalid email format.".to_string())
        );
    }

    #[test]
    fn test_validate_url() {
        // Valid cases
        assert_eq!(validate_url("http://example.com"), Ok(true));
        assert_eq!(validate_url("https://www.example.com"), Ok(true));
        assert_eq!(
            validate_url("https://example.com/path?query=string"),
            Ok(true)
        );
        assert_eq!(validate_url("https://example.com:8080"), Ok(true));
        assert_eq!(validate_url("http://localhost"), Ok(true));

        // Invalid cases
        assert_eq!(
            validate_url(""),
            Err("The URL format is invalid.".to_string())
        );
        assert_eq!(
            validate_url("example.com"),
            Err("The URL format is invalid.".to_string())
        );
        assert_eq!(
            validate_url("ftp://example.com"),
            Err("The URL format is invalid.".to_string())
        );
        assert_eq!(
            validate_url("http//example.com"),
            Err("The URL format is invalid.".to_string())
        );
        assert_eq!(
            validate_url("http:/example.com"),
            Err("The URL format is invalid.".to_string())
        );
        assert_eq!(
            validate_url("://example.com"),
            Err("The URL format is invalid.".to_string())
        );
    }

    #[test]
    fn test_validate_webhook_url() {
        // Valid cases
        assert_eq!(validate_webhook_url("http://example.com"), Ok(true));
        assert_eq!(validate_webhook_url("https://example.com/path"), Ok(true));
        assert_eq!(validate_webhook_url(""), Ok(true)); // empty input is allowed

        // Invalid cases
        assert_eq!(
            validate_webhook_url("invalid_url"),
            Err("The URL format is invalid.".to_string())
        );
        assert_eq!(
            validate_webhook_url("ftp://example.com"),
            Err("The URL format is invalid.".to_string())
        );
    }

    #[test]
    fn test_validate_fallback_url() {
        // Reusing test cases for URL validation
        assert_eq!(validate_fallback_url("http://example.com"), Ok(true));
        assert_eq!(validate_fallback_url("https://example.com/path"), Ok(true));
        assert_eq!(validate_fallback_url(""), Ok(true)); // empty input is allowed

        // Invalid cases
        assert_eq!(
            validate_fallback_url("invalid_url"),
            Err("The URL format is invalid.".to_string())
        );
        assert_eq!(
            validate_fallback_url("example.com"),
            Err("The URL format is invalid.".to_string())
        );
    }
}
