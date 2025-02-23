# Build stage
FROM rust:1.83 as builder

# Install necessary build dependencies
RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    perl \
    libfindbin-libs-perl

# OpenSSL configuration (use system OpenSSL during build)
ENV OPENSSL_NO_VENDOR=1

# Set working directory and copy source code
WORKDIR /usr/src/app
COPY Cargo.toml Cargo.lock ./
COPY src ./src
COPY .env ./
COPY sqlite3.db ./

# Build the app (release mode)
RUN cargo build --release

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies (only libssl3 is needed)
RUN apt-get update && apt-get install -y \
    libssl3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary generated in the build stage
COPY --from=builder /usr/src/app/target/release/rust-url-shortener /usr/local/bin/
COPY --from=builder /usr/src/app/.env /app/

# Set working directory
WORKDIR /app

# Execution command
CMD ["/usr/local/bin/rust-url-shortener"]