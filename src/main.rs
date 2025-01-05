mod handlers;
mod schemas;
mod state;
mod utils;
mod validators;
mod config;

use axum::{
    routing::{get, post},
    Router,
};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;
use tower_http::trace::TraceLayer;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

#[tokio::main]
async fn main() -> Result<(), sqlx::Error> {
    // Initialize DB
    let db_pool = sqlx::sqlite::SqlitePoolOptions::new()
        .max_connections(10)
        .connect("sqlite://sqlite3.db")
        .await
        .expect("Failed to create pool");
    let state = state::AppState {
        db_pool,
        cache: Arc::new(RwLock::new(HashMap::new())),
    };

    // Initialize logger
    tracing_subscriber::registry()
        .with(tracing_subscriber::fmt::layer())
        .init();

    // Configure router
    let app = Router::new()
        .route("/", get(handlers::page_handlers::index_handler))
        .route("/health", get(handlers::page_handlers::health_check_handler))
        .route("/v1/urls", post(handlers::short_url_handlers::create_short_url_handler))
        .route("/v1/verify/{code}", get(handlers::verify_handlers::verify_email_handler))
        .route("/{short_key}", get(handlers::redirect_handlers::redirect_to_original_handler))
        .with_state(state)
        .layer(TraceLayer::new_for_http());

    // Start the server
    let envs = config::get_environments();
    let protocol = &envs.server_protocol;
    let host = &envs.server_host;
    let port = &envs.server_port;
    let listener = tokio::net::TcpListener::bind(format!("{}:{}", host, port))
        .await?;
    println!("Server running on {}://{}:{}", protocol, host, port);
    axum::serve(listener, app).await?;
    Ok(())
}
