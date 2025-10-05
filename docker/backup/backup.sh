#!/usr/bin/env bash

## Rolling backups using External storage
## https://litestream.io/external-storage/

set -ex -o pipefail

# Configuration
BACKUP_DIR="${BACKUP_DIR:-/var/backup}"
DB_FILE="${DB_FILE:-/var/tmp/database.db}"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
BACKUP_FILE="$BACKUP_DIR/db_backup_$TIMESTAMP.sqlite"

# Functions
validate_environment() {
    if [ ! -f "$DB_FILE" ]; then
        echo "Database file $DB_FILE does not exist. Exiting."
        exit 1
    fi
}

create_local_backup() {
    echo "Creating local backup..."
    mkdir -p "$BACKUP_DIR"

    # Backup the sqlite db using VACUUM INTO for better compression
    sqlite3 "$DB_FILE" "VACUUM INTO '$BACKUP_FILE'"

    # Check the integrity of the backup
    sqlite3 "$BACKUP_FILE" 'PRAGMA integrity_check'

    # Compress the backup file
    gzip "$BACKUP_FILE"
    echo "Local backup created: $BACKUP_FILE.gz"
}

install_tool() {
    local tool=$1
    local package=$2

    if ! which "$tool" > /dev/null 2>&1; then
        echo "Installing $tool..."
        pip install "$package"
    fi
}

upload_rolling_backups() {
    local upload_cmd="$1"

    # 1-day, rolling hourly backup
    eval "$upload_cmd backup-$(date +%H).gz"

    # 1-month, rolling daily backup
    eval "$upload_cmd backup-$(date +%d).gz"

    # 1-month, rolling hourly backup
    eval "$upload_cmd backup-$(date +%d%H).gz"
}

backup_to_aws() {
    if [ -z "$AWS_BACKUP_BUCKET" ]; then
        echo "Please set the AWS_BACKUP_BUCKET environment variable"
        exit 1
    fi

    install_tool "aws" "awscli"

    local base_cmd="aws s3 cp \"$BACKUP_FILE.gz\" s3://$AWS_BACKUP_BUCKET/"
    upload_rolling_backups "$base_cmd"

    echo "AWS backup executed"
}

backup_to_azure() {
    if [ -z "$AZURE_CONTAINER" ]; then
        echo "Please set the AZURE_CONTAINER environment variable"
        exit 1
    fi

    install_tool "az" "azure-cli"

    local base_cmd="az storage blob upload --account-name \"$AZURE_STORAGE_ACCOUNT\" --account-key \"$AZURE_STORAGE_KEY\" --container-name \"$AZURE_CONTAINER\" --file \"$BACKUP_FILE.gz\" --name"
    upload_rolling_backups "$base_cmd"

    echo "Azure backup executed"
}

# https://deadmanssnitch.com/
notify_healthcheck() {
    local status=$1
    if [[ -n "$snitch_url" ]]; then
        curl -d s="$status" "$snitch_url" &> /dev/null
    fi
}

main() {
    validate_environment
    create_local_backup

    local backup_status=0

    # AWS backup
    if [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ]; then
        backup_to_aws
        backup_status=$?
    fi

    # Azure backup
    if [ -n "$AZURE_STORAGE_ACCOUNT" ] && [ -n "$AZURE_STORAGE_KEY" ]; then
        backup_to_azure
        backup_status=$?
    fi

    # Health check notification
    notify_healthcheck "$backup_status"
}

# Run main function
main "$@"
