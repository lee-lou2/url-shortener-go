use crate::state::AppState;
use crate::state::CacheEntry;
use crate::utils::converter::{key_to_id, split_short_key};
use axum::{
    body::Body, extract::Path, extract::State, http::Request, http::StatusCode,
    response::Html, response::IntoResponse, response::Redirect,
};
use serde_json::json;
use std::time::Duration;
use std::time::Instant;
use crate::schemas::short_url_schemas::ShortUrl;

/// send web hook
/// Sends the User-Agent and ShortKey to the webhook URL registered by the user
async fn send_web_hook(user_agent: &str, short_key: &str, webhook_url: &str) {
    let client = reqwest::Client::new();
    if let Err(e) = client
        .post(webhook_url)
        .json(&json!({
            "short_key": short_key,
            "user_agent": user_agent,
        })).send().await
    {
        eprintln!("Failed to send webhook: {}", e);
    }
}

/// Update Success WebPage
/// Updates the information of the redirection web page
fn set_success_page(
    short_url: &ShortUrl,
) -> String {
    include_str!("../templates/redirect.html")
        .replace("{ios_deep_link}", &short_url.ios_deep_link)
        .replace("{ios_fallback_url}", &short_url.ios_fallback_url)
        .replace("{android_deep_link}", &short_url.android_deep_link)
        .replace("{android_fallback_url}", &short_url.android_fallback_url)
        .replace("{default_fallback_url}", &short_url.default_fallback_url)
        .replace("{head_html}", &short_url.head_html)
}

/// Redirect URL Handler
/// Redirects to the stored URL for the short link
pub async fn redirect_to_original_handler(
    Path(short_key): Path<String>,
    State(state): State<AppState>,
    req: Request<Body>,
) -> impl IntoResponse {
    // Check Legacy url
    if short_key.len() == 4 {
        let legacy_data = include_str!("../legacy.json");
        let legacy_data: serde_json::Value = serde_json::from_str(legacy_data).unwrap();
        if let Some(legacy_url) = legacy_data.get(&short_key) {
            if let Some(url) = legacy_url.as_str() {
                return (StatusCode::FOUND, Redirect::to(url)).into_response();
            }
        }
    }

    // Check cache
    let cache = state.cache.read().await;
    if let Some(entry) = cache.get(&short_key) {
        if entry.expiry > Instant::now() {
            // Parse entry.data
            let short_url = &entry.data;
            if short_url.webhook_url != "" {
                let webhook_url = short_url.webhook_url.to_string();
                tokio::spawn(async move {
                    let user_agent = req
                        .headers()
                        .get("User-Agent")
                        .and_then(|header_value| header_value.to_str().ok())
                        .unwrap_or("");
                    send_web_hook(&user_agent, &short_key, &webhook_url).await;
                });
            }
            let success_html = set_success_page(short_url);
            return (StatusCode::OK, Html(success_html)).into_response();
        }
    }
    drop(cache);

    let (request_unique_key, request_random_key) = split_short_key(&short_key);
    let url_id = key_to_id(&request_unique_key).unwrap_or(0);

    // // If not in cache, query the database
    match sqlx::query!(
        r#"
        SELECT
            random_key, ios_deep_link, ios_fallback_url,
            android_deep_link, android_fallback_url,
            default_fallback_url, webhook_url, head_html
        FROM urls
        WHERE id = ?
            AND is_deleted = 0
            AND is_verified = 1
        "#,
        url_id
    )
        .fetch_optional(&state.db_pool)
        .await {
        Ok(Some(record)) => {
            if record.random_key != request_random_key {
                return (StatusCode::NOT_FOUND, "URL not found").into_response();
            }
            let short_url = ShortUrl {
                ios_deep_link: record.ios_deep_link.unwrap_or(String::from("")),
                ios_fallback_url: record.ios_fallback_url.unwrap_or(String::from("")),
                android_deep_link: record.android_deep_link.unwrap_or(String::from("")),
                android_fallback_url: record.android_fallback_url.unwrap_or(String::from("")),
                default_fallback_url: record.default_fallback_url,
                webhook_url: record.webhook_url.unwrap_or(String::from("")),
                head_html: record.head_html.unwrap_or(String::from("")),
            };
            state.cache.write().await.insert(short_key.clone(), CacheEntry {
                data: short_url.clone(),
                expiry: Instant::now() + Duration::from_secs(3600),
            });
            if short_url.webhook_url != "" {
                // Send the webhook
                let short_url_clone = short_url.clone();
                tokio::spawn(async move {
                    let user_agent = req
                        .headers()
                        .get("User-Agent")
                        .and_then(|header_value| header_value.to_str().ok())
                        .unwrap_or("");
                    let webhook_url = short_url_clone.webhook_url.to_string();
                    send_web_hook(&user_agent, &short_key, &webhook_url).await;
                });
            }
            let success_html = set_success_page(&short_url);
            (StatusCode::OK, Html(success_html)).into_response()
        }
        Ok(None) => {
            (StatusCode::NOT_FOUND, "URL not found").into_response()
        }
        Err(e) => {
            eprintln!("DB Error: {:?}", e);
            StatusCode::INTERNAL_SERVER_ERROR.into_response()
        }
    }
}
