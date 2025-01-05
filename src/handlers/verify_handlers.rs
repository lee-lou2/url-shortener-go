use crate::config;
use crate::state::AppState;
use crate::utils::converter::{key_to_id, split_short_key};
use axum::extract::State;
use axum::{extract::Path, http::StatusCode, response::Html, response::IntoResponse};


/// Verify Email Handler
/// Handler to verify the verification link sent via email
pub async fn verify_email_handler(
    State(state): State<AppState>,
    Path(code): Path<String>,
) -> impl IntoResponse {
    // Find short_key using the verification code
    match sqlx::query!(
        r#"
        SELECT short_key FROM email_auth WHERE code = ?1 AND expires_at > datetime('now')
        "#,
        code
    )
        .fetch_one(&state.db_pool)
        .await {
        Ok(record) => {
            let (unique_key, random_key) = split_short_key(&record.short_key);
            let url_id = key_to_id(&unique_key).unwrap_or(0);

            // Update URL verification status
            let result = sqlx::query!(
                r#"
                UPDATE urls SET is_verified = true WHERE random_key = ?1 AND id = ?2
                "#,
                random_key,
                url_id
            )
                .execute(&state.db_pool)
                .await;

            match result {
                Ok(_) => {
                    // Delete the verified code
                    let envs = config::get_environments();
                    let protocol = &envs.server_protocol;
                    let host = &envs.server_host;
                    let port = &envs.server_port;
                    let domain = if port != "80" && port != "443" {
                        format!("{}:{}", host, port)
                    } else { host.to_string() };
                    let short_url = format!("{}://{}/{}", protocol, domain, record.short_key);
                    let _ = sqlx::query!(
                        r#"
                        DELETE FROM email_auth WHERE code = ?1
                        "#,
                        code
                    )
                        .execute(&state.db_pool)
                        .await;
                    let success_html =
                        include_str!("../templates/verify/success.html").replace("{short_url}", &short_url);
                    (StatusCode::OK, Html(success_html)).into_response()
                }
                Err(_) => {
                    (
                        StatusCode::INTERNAL_SERVER_ERROR,
                        Html(include_str!("../templates/verify/error.html"))
                    ).into_response()
                }
            }
        }
        Err(_) => {
            (
                StatusCode::INTERNAL_SERVER_ERROR,
                Html(include_str!("../templates/verify/failed.html")),
            ).into_response()
        }
    }
}
