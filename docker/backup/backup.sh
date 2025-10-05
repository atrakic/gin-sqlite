#!/usr/bin/env bash

## Rolling backups using External storage
## https://litestream.io/external-storage/

set -ex -o pipefail

# Set the path to the backup directory
BACKUP_DIR="${BACKUP_DIR:-/var/backup}"

# This is the db that you want to backup
DB_FILE="${DB_FILE:-/var/tmp/database.db}"

###

# Create the backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Create a timestamp for the backup file
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")

# Create the backup file name
BACKUP_FILE="$BACKUP_DIR/db_backup_$TIMESTAMP.sqlite"

# Backup the sqlite db
#sqlite3 "$DB_FILE" ".backup $BACKUP_FILE"
sqlite3 "$DB_FILE" "VACUUM INTO '$BACKUP_FILE'"

# Check the integrity of the backup
sqlite3 "$BACKUP_FILE" 'PRAGMA integrity_check'

# gzip the backup file
gzip "$BACKUP_FILE"

if [ -n "$AWS_ACCESS_KEY_ID" ] && [ -n "$AWS_SECRET_ACCESS_KEY" ]; then

    if [ -n "$AWS_BACKUP_BUCKET" ]; then
        echo "Please set the AWS_BACKUP_BUCKET environment variable"
        exit 1
    fi

    # Install aws cli on alpine if not already installed
    which aws || pip install awscli

    # 1-day, rolling hourly backup
    aws s3 cp "$BACKUP_FILE.gz" s3://"$AWS_BACKUP_BUCKET"/backup-"$(date +%H)".gz

    # 1-month, rolling daily backup
    aws s3 cp "$BACKUP_FILE.gz" s3://"$AWS_BACKUP_BUCKET"/backup-"$(date +%d)".gz

    # 1-month, rolling hourly backup
    aws s3 cp "$BACKUP_FILE.gz" s3://"$AWS_BACKUP_BUCKET"/backup-"$(date +%d%H)".gz
    backup_status=$?

    echo "AWS backup executed"
fi

## Backup to Azure Blob Storage
# https://docs.microsoft.com/en-us/azure/storage/blobs/storage-quickstart-blobs-cli

if [ -n "$AZURE_STORAGE_ACCOUNT" ] && [ -n "$AZURE_STORAGE_KEY" ];
then

    if [ -n "$AZURE_CONTAINER" ]; then
        echo "Please set the AZURE_CONTAINER environment variable"
        exit 1
    fi

    # Install azure-cli on alpine
    which aws || pip install azure-cli

    # 1-day, rolling hourly backup
    az storage blob upload --account-name "$AZURE_STORAGE_ACCOUNT" --account-key "$AZURE_STORAGE_KEY" --container-name "$AZURE_CONTAINER" --file "$BACKUP_FILE.gz" --name backup-"$(date +%H)".gz

    # 1-month, rolling daily backup
    az storage blob upload --account-name "$AZURE_STORAGE_ACCOUNT" --account-key "$AZURE_STORAGE_KEY" --container-name "$AZURE_CONTAINER" --file "$BACKUP_FILE.gz" --name backup-"$(date +%d)".gz

    # 1-month, rolling hourly backup
    az storage blob upload --account-name "$AZURE_STORAGE_ACCOUNT" --account-key "$AZURE_STORAGE_KEY" --container-name "$AZURE_CONTAINER" --file "$BACKUP_FILE.gz" --name backup-"$(date +%d%H)".gz
    backup_status=$?

    echo "Azure backup executed"
fi

# https://deadmanssnitch.com/
# Notify dead man snitch if back up completed successfully.
if [[ -n "$snitch_url" ]]; then curl -d s="$backup_status" "$snitch_url" &> /dev/null; fi
