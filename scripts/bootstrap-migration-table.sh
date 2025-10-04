#!/usr/bin/env bash

MIGRATION_PATH=migrations/bootstrap-migration-table.sql

PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} \
  -d ${DB_NAME} -a -f ${MIGRATION_PATH}
