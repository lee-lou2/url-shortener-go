#!/bin/bash

DB_FILE="sqlite3.db"

# 데이터베이스 파일 확인 및 생성
if [ ! -f "$DB_FILE" ]; then
  echo "Database file does not exist. Creating..."
  sqlite3 "$DB_FILE" <<EOF
CREATE TABLE IF NOT EXISTS urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    random_key VARCHAR(4) NOT NULL,
    email VARCHAR(255) NOT NULL,
    ios_deep_link TEXT NULL,
    ios_fallback_url TEXT NULL,
    android_deep_link TEXT NULL,
    android_fallback_url TEXT NULL,
    default_fallback_url TEXT NOT NULL,
    hashed_value TEXT NOT NULL UNIQUE,
    webhook_url TEXT NULL,
    head_html TEXT NULL,
    is_verified INTEGER NOT NULL DEFAULT 0,
    is_deleted INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_hashed_value ON urls (hashed_value);

CREATE TABLE IF NOT EXISTS email_auth (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    short_key VARCHAR(10) NOT NULL,
    code VARCHAR(8) NOT NULL,
    expires_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS code ON email_auth (code);
EOF
  echo "Database initialized."
else
  echo "Database file already exists."
fi
