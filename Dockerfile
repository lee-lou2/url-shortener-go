FROM rust:1.83-slim

RUN apt-get update && apt-get install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && apt-get install -y \
    perl \
    libssl-dev \
    libfindbin-libs-perl

ENV OPENSSL_NO_VENDOR=1
