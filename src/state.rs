use crate::schemas::short_url_schemas::ShortUrl;
use sqlx::SqlitePool;
use std::collections::HashMap;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tokio::sync::RwLock;
use tokio::time::interval;

pub struct CacheEntry {
    pub data: ShortUrl,
    pub expiry: Instant,
}

#[derive(Clone)]
pub struct AppState {
    pub db_pool: SqlitePool,
    pub cache: Arc<RwLock<HashMap<String, CacheEntry>>>,
}

impl AppState {
    pub fn new(db_pool: SqlitePool) -> Self {
        let state = Self {
            db_pool,
            cache: Arc::new(RwLock::new(HashMap::new())),
        };
        let state_clone = state.clone();
        tokio::spawn(async move {
            let duration = Duration::from_secs(60);
            loop {
                interval(duration).tick().await;
                let mut cache = state_clone.cache.write().await;
                let now = Instant::now();
                cache.retain(|_, entry| entry.expiry > now);
            }
        });
        state
    }
}
