#!/bin/bash
IMAGE_NAME="rust-url-shortener"
IMAGE_TAG="latest"
INTERNAL_SERVER_PORT=3000
EXTERNAL_SERVER_PORT=3000

docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .
if docker ps -a | grep -q ${IMAGE_NAME}; then
  docker stop ${IMAGE_NAME}
  docker rm ${IMAGE_NAME}
fi

docker run --name ${IMAGE_NAME} \
  -v .:/app \
  -w /app \
  -d \
  -p ${INTERNAL_SERVER_PORT}:${EXTERNAL_SERVER_PORT} \
  ${IMAGE_NAME}:${IMAGE_TAG}