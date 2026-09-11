#!/usr/bin/env bash

set -euo pipefail

DOMAINS=("api.anr-studio.com")
EMAIL="${CERTBOT_EMAIL:-admin@anr-studio.com}"
DATA_PATH="./infrastructure/gateway/certbot"

mkdir -p "$DATA_PATH/www"
mkdir -p "$DATA_PATH/conf"

for domain in "${DOMAINS[@]}"; do
  CERT_DIR="$DATA_PATH/conf/live/$domain"
  if [ ! -f "$CERT_DIR/fullchain.pem" ]; then
    mkdir -p "$CERT_DIR"
    openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
      -keyout "$CERT_DIR/privkey.pem" \
      -out "$CERT_DIR/fullchain.pem" \
      -subj "/CN=localhost"
  fi
done

if [ "${1:-}" = "--dummy-only" ]; then
  exit 0
fi

docker compose up -d gateway

for domain in "${DOMAINS[@]}"; do
  docker compose run --rm --entrypoint "\
    certbot certonly --webroot -w /var/www/certbot \
    --email $EMAIL \
    -d $domain \
    --rsa-key-size 4096 \
    --agree-tos \
    --force-renewal \
    --non-interactive" certbot
done

docker compose exec gateway nginx -s reload
