use axum::{
    body::Body,
    http::{header::AUTHORIZATION, Request},
    middleware::Next,
    response::IntoResponse,
};
use jsonwebtoken::{decode, DecodingKey, Validation};
use serde::{Deserialize, Serialize};

/// JWT Token Claims
#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct Claims {
    pub sub: String,
    pub exp: usize,
}

/// JWT Authentication Middleware
/// If a token exists, extract claims; if not, just move to the next step
pub async fn jwt_auth_middleware(mut req: Request<Body>, next: Next) -> impl IntoResponse {
    let mut claims = Claims {
        sub: "".to_string(),
        exp: 0,
    };
    if let Some(auth_value) = req.headers().get(AUTHORIZATION) {
        if let Ok(auth_str) = auth_value.to_str() {
            if let Some(token) = auth_str.strip_prefix("Bearer ") {
                let envs = crate::config::get_environments();
                if let Ok(data) = decode::<Claims>(
                    token,
                    &DecodingKey::from_secret(&envs.jwt_secret.as_bytes()),
                    &Validation::default(),
                ) {
                    claims = data.claims;
                } else {
                    println!("JWT token validation failed");
                }
            }
        }
    }
    req.extensions_mut().insert(claims);
    next.run(req).await
}
