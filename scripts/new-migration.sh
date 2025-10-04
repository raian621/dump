#!/usr/bin/env bash

MIGRATION_DIR='migrations'
MIGRATION_RECORD_FILENAME='_records.txt'

mkdir -p db/migrations

if [ "$0" = "help" ]; then
  echo "Usage: $0 <migration_name>"
  exit
fi

MIGRATION_NAME="$1"
MIGRATION_FILENAME="${MIGRATION_NAME}.sql"

echo "${MIGRATION_FILENAME}" >> "${MIGRATION_DIR}/${MIGRATION_RECORD_FILENAME}"
touch "${MIGRATION_DIR}/${MIGRATION_FILENAME}"