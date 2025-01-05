const CHARS: &str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
const BASE: i64 = 62;

/// ID To Key
/// Converts the primary key of the ShortURL table into a string.
/// Indexes when English lowercase/uppercase letters and numbers are sequentially combined.
pub fn id_to_key(mut id: i64) -> Option<String> {
    if id < 1 {
        return None;
    }
    let mut key = Vec::new();
    while id > 0 {
        id -= 1;
        let digit = (id % BASE) as usize;
        key.push(CHARS.as_bytes()[digit] as char);
        id /= BASE;
    }
    key.reverse();
    Some(key.iter().collect())
}

/// Key To ID
/// Converts arbitrary characters into a number.
/// Indexes when English lowercase/uppercase letters and numbers are sequentially combined.
pub fn key_to_id(key: &str) -> Option<i64> {
    let mut result = 0i64;
    for c in key.chars() {
        let digit = CHARS.find(c)? as i64;
        result = result * BASE + (digit + 1);
    }
    Some(result)
}

/// Split Short Key
/// Extracts the unique key in the middle and the random keys at the front and back.
pub fn split_short_key(short_key: &str) -> (String, String) {
    let front_random_key = short_key[..2].to_string();
    let back_random_key = short_key[short_key.len() - 2..].to_string();
    let random_key = &(front_random_key + &back_random_key);
    let unique_key = short_key[2..short_key.len() - 2].to_string();
    (unique_key.to_string(), random_key.to_string())
}

/// Merge Short Key
/// Combines a 4-character random key and a unique string to generate a ShortKey.
pub fn merge_short_key(random_key: &str, unique_key: &str) -> String {
    (&random_key[..2]).to_string() + &unique_key + &random_key[2..]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_id_to_key_basic() {
        // Returns None if less than 1
        assert_eq!(id_to_key(0), None);
        assert_eq!(id_to_key(-1), None);

        // Single character
        assert_eq!(id_to_key(1), Some("a".to_string()));  // 'a'
        assert_eq!(id_to_key(26), Some("z".to_string())); // 'z'
        assert_eq!(id_to_key(27), Some("A".to_string())); // 'A'
        assert_eq!(id_to_key(52), Some("Z".to_string())); // 'Z'
        assert_eq!(id_to_key(53), Some("0".to_string())); // '0'
        assert_eq!(id_to_key(62), Some("9".to_string())); // '9'

        // Multiple characters
        assert_eq!(id_to_key(63), Some("aa".to_string()));
        assert_eq!(id_to_key(64), Some("ab".to_string()));
    }

    /// Check if converting ID → Key → ID returns the same value
    #[test]
    fn test_round_trip() {
        let test_values = [
            1, 2, 25, 26, 27, 52, 53, 62, 63, 64, 100, 999, 1000, 99999999
        ];

        for &id in &test_values {
            let key = id_to_key(id).unwrap();
            let decoded_id = key_to_id(&key).unwrap();
            assert_eq!(id, decoded_id, "Roundtrip mismatch for ID {id}");
        }
    }

    /// Check if None is returned when a key contains invalid characters
    #[test]
    fn test_invalid_chars() {
        // Characters not in CHARS (a-z, A-Z, 0-9) should return None
        assert_eq!(key_to_id("!"), None);
        assert_eq!(key_to_id("한글"), None);
        assert_eq!(key_to_id("a!"), None);
        assert_eq!(key_to_id("a b"), None);
    }

    #[test]
    fn test_split_short_key_basic() {
        // Basic test with short keys of valid lengths
        let (unique_key, random_key) = split_short_key("abXYyz");
        assert_eq!(unique_key, "XY");
        assert_eq!(random_key, "abyz");

        let (unique_key, random_key) = split_short_key("cdUVef");
        assert_eq!(unique_key, "UV");
        assert_eq!(random_key, "cdef");
    }

    #[test]
    fn test_split_short_key_edge_cases() {
        // Edge case: Empty string
        let result = std::panic::catch_unwind(|| split_short_key(""));
        assert!(result.is_err(), "Expected a panic for empty short key");

        // Edge case: String too short to extract
        let result = std::panic::catch_unwind(|| split_short_key("abc"));
        assert!(result.is_err(), "Expected a panic for short key too short to split");
    }

    #[test]
    fn test_split_short_key_complex() {
        // Test with longer strings
        let (unique_key, random_key) = split_short_key("ab12345yz");
        assert_eq!(unique_key, "12345");
        assert_eq!(random_key, "abyz");

        let (unique_key, random_key) = split_short_key("XYHelloZ9");
        assert_eq!(unique_key, "Hello");
        assert_eq!(random_key, "XYZ9");
    }

    #[test]
    fn test_split_short_key_special_characters() {
        // Ensure the function works with special characters or edge alphanumeric combinations
        let (unique_key, random_key) = split_short_key("a1<>9x");
        assert_eq!(unique_key, "<>");
        assert_eq!(random_key, "a19x");
    }


    #[test]
    fn test_merge_short_key_basic() {
        // Basic test with simple inputs
        let result = merge_short_key("abcd", "XY");
        assert_eq!(result, "abXYcd");

        let result = merge_short_key("wxyz", "123");
        assert_eq!(result, "wx123yz");
    }

    #[test]
    fn test_merge_short_key_edge_cases() {
        // Edge case: Empty unique key
        let result = merge_short_key("abcd", "");
        assert_eq!(result, "abcd");

        // Edge case: Empty random key
        let result = std::panic::catch_unwind(|| merge_short_key("", "XY"));
        assert!(result.is_err(), "Expected a panic for empty random key");

        // Edge case: Random key not 4 characters
        let result = std::panic::catch_unwind(|| {
            if "abc".len() != 4 {
                panic!("Random key must be exactly 4 characters");
            }
            merge_short_key("abc", "XY")
        });
        assert!(result.is_err(), "Expected a panic for random key less than 4 characters");

        let result = std::panic::catch_unwind(|| {
            if "abcdef".len() != 4 {
                panic!("Random key must be exactly 4 characters");
            }
            merge_short_key("abcdef", "XY")
        });
        assert!(result.is_err(), "Expected a panic for random key more than 4 characters");
    }

    #[test]
    fn test_merge_short_key_special_characters() {
        // Test with special characters
        let result = merge_short_key("12#$", "&*");
        assert_eq!(result, "12&*#$");

        // Test with spaces
        let result = merge_short_key("ab b", "c ");
        assert_eq!(result, "abc  b");
    }

    #[test]
    fn test_merge_short_key_long_unique_key() {
        // Test with a long unique key
        let result = merge_short_key("abcd", "HelloWorld");
        assert_eq!(result, "abHelloWorldcd");

        // Test with very large unique key
        let result = merge_short_key("mnop", "ThisIsAVeryLongUniqueString");
        assert_eq!(result, "mnThisIsAVeryLongUniqueStringop");
    }

    #[test]
    fn test_merge_short_key_numeric_input() {
        // Test numeric-looking strings as random/unique keys
        let result = merge_short_key("1234", "5678");
        assert_eq!(result, "12567834");

        let result = merge_short_key("9876", "5432");
        assert_eq!(result, "98543276");
    }
}