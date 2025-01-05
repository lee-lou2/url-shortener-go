use crate::config;
use crate::schemas::short_url_schemas::{CreateUrlRequest, CreateUrlResponse};
use crate::state::AppState;
use crate::utils::converter::{id_to_key, merge_short_key};
use crate::utils::generator::generate_random_string;
use crate::validators::validate_url::{validate_email, validate_fallback_url, validate_url, validate_webhook_url};
use axum::extract::State;
use axum::{
    http::StatusCode,
    response::IntoResponse,
    Json,
};
use lettre::{transport::smtp::authentication::Credentials, Message, SmtpTransport, Transport};
use scraper::Html as ScraperHtml;
use sha2::{Digest, Sha256};
use chrono::{Utc, Duration};


/// Send Email
/// Send a verification URL via email
async fn send_email(email: String, code: String) -> Result<(), lettre::transport::smtp::Error> {
    let envs = config::get_environments();
    let protocol = &envs.server_protocol;
    let host = &envs.server_host;
    let port = &envs.server_port;
    let domain = if port != "80" && port != "443" {
        format!("{}:{}", host, port)
    } else { host.to_string() };
    let email_body = format!(
        "{}://{}/v1/verify/{}\n\nThis code is valid for 5 minutes.", protocol, domain, code
    );

    let from_email = &envs.email_address;
    let user_name = &envs.email_user_name;
    let password = &envs.email_password;
    let email_host = &envs.email_host;
    let email_port = &envs.email_port;
    let creds = Credentials::new(user_name.to_string(), password.to_string());

    let mailer = SmtpTransport::relay(email_host)?
        .credentials(creds)
        .port(email_port.parse().unwrap())
        .build();

    let email = Message::builder()
        .from(from_email.parse().unwrap())
        .to(email.parse().unwrap())
        .subject("[F-IT] Verification for Short Link Creation")
        .header(lettre::message::header::ContentType::TEXT_PLAIN)
        .body(email_body.as_bytes().to_vec())
        .unwrap();

    match mailer.send(&email) {
        Ok(_) => Ok(()),
        Err(e) => Err(e),
    }
}

/// Extract Head HTML
/// Retrieve information of the head tag from the HTML
fn extract_head_html(html: &str) -> String {
    let document = ScraperHtml::parse_document(html);
    let selector = scraper::Selector::parse("head").unwrap();
    if let Some(head) = document.select(&selector).next() {
        head.html()
    } else {
        String::new()
    }
}

/// URL shortening handler
pub async fn create_short_url_handler(
    State(state): State<AppState>,
    Json(payload): Json<CreateUrlRequest>,
) -> impl IntoResponse {
    // Data validation
    fn validate_data(payload: &CreateUrlRequest) -> Result<bool, String> {
        validate_email(&payload.email)?;
        validate_url(&payload.default_fallback_url)?;
        validate_webhook_url(&payload.webhook_url)?;
        validate_fallback_url(&payload.default_fallback_url)?;
        Ok(true)
    }

    if let Err(e) = validate_data(&payload) {
        return (StatusCode::BAD_REQUEST, e).into_response();
    }

    // Generate unique ID
    let mut hasher = Sha256::new();
    hasher.update(format!(
        "{}{}{}{}{}",
        payload.ios_deep_link,
        payload.ios_fallback_url,
        payload.android_deep_link,
        payload.android_fallback_url,
        payload.default_fallback_url
    ));
    let hashed_value = format!("{:x}", hasher.finalize());

    // If hashed_value already exists, return it as is
    match sqlx::query!(
        r#"
        SELECT
            id, email, random_key, is_verified
        FROM urls
        WHERE hashed_value = ?1
            AND is_deleted = 0
        "#,
        hashed_value
    )
        .fetch_optional(&state.db_pool)
        .await {
        Ok(Some(record)) => {
            if record.is_verified == 1 {
                return (StatusCode::CONFLICT, "Email is already verified.").into_response();
            }
            let id = record.id.unwrap();
            let unique_key = match id_to_key(id) {
                Some(key) => key,
                None => return (StatusCode::INTERNAL_SERVER_ERROR, "ID conversion failed").into_response(),
            };
            let short_key = merge_short_key(&record.random_key, &unique_key);

            // Add to email authentication table
            let code = generate_random_string(8);
            let expires_at = Utc::now() + Duration::minutes(5);
            let expires_at_str = expires_at.naive_utc().to_string();
            let _ = sqlx::query!(
                r#"
                INSERT INTO email_auth (short_key, code, expires_at) VALUES (?, ?, ?)
                "#,
                short_key,
                code,
                expires_at_str
            )
                .execute(&state.db_pool)
                .await;
            tokio::spawn(async move {
                if let Err(e) = send_email(record.email, code).await {
                    eprintln!("Failed to send email: {}", e);
                }
            });
            let response = CreateUrlResponse { is_created: false };
            return (StatusCode::CREATED, Json(response)).into_response();
        },
        Ok(None) => {
            // pass
        },
        Err(e) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                format!("Failed to query database: {}", e),
            ).into_response();
        }
    }

    // create pending short url
    let random_key = generate_random_string(4);
    let result = sqlx::query!(
        r#"
        INSERT INTO urls (
            random_key, email, ios_deep_link, ios_fallback_url,
            android_deep_link, android_fallback_url, default_fallback_url,
            hashed_value, webhook_url, head_html
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        "#,
        random_key,
        payload.email,
        payload.ios_deep_link,
        payload.ios_fallback_url,
        payload.android_deep_link,
        payload.android_fallback_url,
        payload.default_fallback_url,
        hashed_value,
        payload.webhook_url,
        payload.head_html
    )
        .execute(&state.db_pool)
        .await;

    match result {
        Ok(result) => {
            let id = result.last_insert_rowid();
            let unique_key = match id_to_key(id) {
                Some(key) => key,
                None => return (StatusCode::INTERNAL_SERVER_ERROR, "ID conversion failed").into_response(),
            };
            // Add to cache
            let short_key = merge_short_key(&random_key, &unique_key);
            let code = generate_random_string(8);
            let expires_at = Utc::now() + Duration::minutes(5);
            let expires_at_str = expires_at.naive_utc().to_string();
            let _ = sqlx::query!(
                r#"
                INSERT INTO email_auth (short_key, code, expires_at) VALUES (?, ?, ?)
                "#,
                short_key,
                code,
                expires_at_str
            )
                .execute(&state.db_pool)
                .await;
            tokio::spawn(async move {
                if let Err(e) = send_email(payload.email, code).await {
                    println!("Failed to send email: {}", e);
                }

                // If Head is not provided during creation request, fetch the default URL's Head code and apply it
                if payload.head_html.is_empty() {
                    let client = reqwest::Client::new();
                    match client.get(&payload.default_fallback_url).send().await {
                        Ok(response) => {
                            if let Ok(html) = response.text().await {
                                // Use cloned state
                                let head_html = extract_head_html(&html);
                                let _ = sqlx::query!(
                                    r#"
                                    UPDATE urls SET head_html = ?1 WHERE id = ?2
                                    "#,
                                    head_html,
                                    id,
                                )
                                    .execute(&state.db_pool)
                                    .await;
                            }
                        }
                        Err(e) => println!("Failed to fetch head HTML: {}", e),
                    }
                }
            });
            let response = CreateUrlResponse {
                is_created: true,
            };
            (StatusCode::CREATED, Json(response)).into_response()
        }
        Err(_) => {
            println!("Save failed");
            println!("payload: {:?}", payload);
            println!("random_key: {}", random_key);
            println!("hashed_value: {}", hashed_value);
            (StatusCode::INTERNAL_SERVER_ERROR, "Failed to save").into_response()
        }
    }
}
