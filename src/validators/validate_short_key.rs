pub fn validate_short_key(input: &str) -> Result<bool, String> {
    // Check the length
    if input.len() < 4 {
        return Err("Short key must be at least 4 characters".to_string());
    }
    // Initialize condition flags
    let (mut has_digit, mut has_lowercase, mut has_uppercase) = (false, false, false);
    // Check each character
    input.chars().for_each(|c| {
        if c.is_ascii_digit() {
            has_digit = true;
        } else if c.is_ascii_lowercase() {
            has_lowercase = true;
        } else if c.is_ascii_uppercase() {
            has_uppercase = true;
        }
    });
    // Return true if all conditions are met
    Ok(has_digit && has_lowercase && has_uppercase)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_validate_short_key() {
        // Valid cases
        assert_eq!(validate_short_key("aB1c"), Ok(true)); // Contains lowercase, uppercase, and digit
        assert_eq!(validate_short_key("1aBc"), Ok(true)); // Contains all categories
        assert_eq!(validate_short_key("Password1"), Ok(true)); // Example with a typical password

        // Invalid cases
        assert_eq!(
            validate_short_key("abc"),
            Err("Short key must be at least 4 characters".to_string())
        ); // Only lowercase, less than 4 chars
        assert_eq!(
            validate_short_key("Abc"),
            Err("Short key must be at least 4 characters".to_string())
        ); // Only lowercase and uppercase
        assert_eq!(
            validate_short_key("123"),
            Err("Short key must be at least 4 characters".to_string())
        ); // Only digits, less than 4 chars
        assert_eq!(validate_short_key("1234"), Ok(false)); // Only digits, 4 chars
        assert_eq!(validate_short_key("abcd"), Ok(false)); // Only lowercase, 4 chars
        assert_eq!(validate_short_key("ABCD"), Ok(false)); // Only uppercase, 4 chars
        assert_eq!(validate_short_key("a1b2"), Ok(false)); // Only lowercase and digits
        assert_eq!(validate_short_key("A1B2"), Ok(false)); // Only uppercase and digits
        assert_eq!(validate_short_key("ABc1"), Ok(true)); // Valid case with all conditions satisfied
    }
}
