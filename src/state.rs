use sqlx::SqlitePool;
use std::collections::HashMap;
use std::time::Instant;
use tokio::sync::RwLock;
use std::sync::Arc;
use crate::schemas::short_url_schemas::ShortUrl;

pub struct CacheEntry {
    pub data: ShortUrl,
    pub expiry: Instant,
}

#[derive(Clone)]
pub struct AppState {
    pub db_pool: SqlitePool,
    pub cache: Arc<RwLock<HashMap<String, CacheEntry>>>,
}
