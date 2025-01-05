use serde::{Deserialize, Serialize};

/// URL request struct
#[derive(Deserialize, Debug)]
pub struct CreateUrlRequest {
    pub email: String,
    #[serde(rename = "iosDeepLink")]
    pub ios_deep_link: String,
    #[serde(rename = "iosFallbackUrl")]
    pub ios_fallback_url: String,
    #[serde(rename = "androidDeepLink")]
    pub android_deep_link: String,
    #[serde(rename = "androidFallbackUrl")]
    pub android_fallback_url: String,
    #[serde(rename = "defaultFallbackUrl")]
    pub default_fallback_url: String,
    #[serde(rename = "webhookUrl")]
    pub webhook_url: String,
    #[serde(rename = "headHtml")]
    pub head_html: String,
}

/// URL response struct
#[derive(Serialize)]
pub struct CreateUrlResponse {
    pub is_created: bool,
}

/// ShortUrl struct
#[derive(Clone)]
pub struct ShortUrl {
    pub ios_deep_link: String,
    pub ios_fallback_url: String,
    pub android_deep_link: String,
    pub android_fallback_url: String,
    pub default_fallback_url: String,
    pub webhook_url: String,
    pub head_html: String,
}
