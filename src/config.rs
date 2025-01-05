use dotenv::dotenv;
use once_cell::sync::Lazy;
use std::env;

#[derive(Debug)]
pub struct Environment {
    pub server_protocol: String,
    pub server_host: String,
    pub server_port: String,
    pub email_address: String,
    pub email_user_name: String,
    pub email_password: String,
    pub email_host: String,
    pub email_port: String,
    pub jwt_secret: String,
}

// Lazy를 통해 최초 접근 시점에만 .env를 로드하고 Environment를 초기화
static ENVIRONMENTS: Lazy<Environment> = Lazy::new(|| {
    // .env 로드
    dotenv().ok();

    // 구조체 이름과 맞춰서 Environment { ... }로 수정
    Environment {
        server_protocol: env::var("SERVER_PROTOCOL").unwrap_or_else(|_| "http".to_string()),
        server_host: env::var("SERVER_HOST").unwrap_or_else(|_| "localhost".to_string()),
        server_port: env::var("SERVER_PORT").unwrap_or_else(|_| "8080".to_string()),
        email_address: env::var("EMAIL_ADDRESS").unwrap_or_else(|_| "".to_string()),
        email_user_name: env::var("EMAIL_USER_NAME").unwrap_or_else(|_| "".to_string()),
        email_password: env::var("EMAIL_PASSWORD").unwrap_or_else(|_| "".to_string()),
        email_host: env::var("EMAIL_HOST").unwrap_or_else(|_| "".to_string()),
        email_port: env::var("EMAIL_PORT").unwrap_or_else(|_| "".to_string()),
        jwt_secret: env::var("JWT_SECRET").unwrap_or_else(|_| "".to_string()),
    }
});

pub fn get_environments() -> &'static Environment {
    &ENVIRONMENTS
}