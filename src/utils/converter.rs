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
