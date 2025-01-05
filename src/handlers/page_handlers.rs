use axum::response::Html;

/// Health check handler
/// GET /health
/// Response str
pub async fn health_check_handler() -> &'static str {
    "OK"
}

/// Main page handler
/// GET /
/// Response HTML
pub async fn index_handler() -> Html<&'static str> {
    Html(include_str!("../templates/index.html"))
}
