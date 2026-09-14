#!/usr/bin/env bash

set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/opt/anr-platform/infrastructure/storage/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/anr_platform_${TIMESTAMP}.sql.gz"

mkdir -p "$BACKUP_DIR"

docker compose -f /opt/anr-platform/docker-compose.yml exec -T postgres pg_dump -U "${DB_USER:-anr_user}" -d "${DB_NAME:-anr_platform}" | gzip > "$BACKUP_FILE"

find "$BACKUP_DIR" -type f -name "anr_platform_*.sql.gz" -mtime +"$RETENTION_DAYS" -delete
