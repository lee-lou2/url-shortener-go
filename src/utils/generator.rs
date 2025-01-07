/// Generate Random String
/// Generates a random string of a specified length
pub fn generate_random_string(length: usize) -> String {
    if length == 0 {
        return String::new();
    }
    let chars = [
        'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r',
        's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
        'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1',
        '2', '3', '4', '5', '6', '7', '8', '9',
    ];
    nanoid::nanoid!(length, &chars)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_generate_random_string_length() {
        let length = 10;
        let random_string = generate_random_string(length);
        assert_eq!(random_string.len(), length);
    }

    #[test]
    fn test_generate_random_string_is_unique() {
        let length = 16;
        let string1 = generate_random_string(length);
        let string2 = generate_random_string(length);
        assert_ne!(string1, string2, "Generated strings should be unique.");
    }

    #[test]
    fn test_generate_random_string_with_zero_length() {
        let length = 0;
        let random_string = generate_random_string(length);
        assert_eq!(
            random_string.len(),
            length,
            "Random string with length 0 should be empty."
        );
    }

    #[test]
    fn test_generate_random_string_contains_only_valid_chars() {
        let length = 32;
        let random_string = generate_random_string(length);
        for ch in random_string.chars() {
            assert!(
                ch.is_alphanumeric(),
                "Character '{}' is not alphanumeric.",
                ch
            );
        }
    }

    #[test]
    fn test_generate_random_string_large_length() {
        let length = 1000;
        let random_string = generate_random_string(length);
        assert_eq!(
            random_string.len(),
            length,
            "Random string of 1000 characters should have correct length."
        );
    }
}
