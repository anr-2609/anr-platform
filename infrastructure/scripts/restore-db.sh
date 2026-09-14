#!/usr/bin/env bash

set -euo pipefail

BACKUP_FILE="${1:-}"

if [ -z "$BACKUP_FILE" ] || [ ! -f "$BACKUP_FILE" ]; then
  echo "Usage: $0 <path_to_backup_file.sql.gz>"
  exit 1
fi

gunzip -c "$BACKUP_FILE" | docker compose -f /opt/anr-platform/docker-compose.yml exec -T postgres psql -U "${DB_USER:-anr_user}" -d "${DB_NAME:-anr_platform}"
