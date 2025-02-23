make build

export BACKUP_SOURCE_NAME="mysql"
export BACKUP_DB_NAME="testdb"
export BACKUP_USERNAME="testuser"
export BACKUP_PASSWORD="testpass"
export BACKUP_HOST="127.0.0.1"
export BACKUP_PORT="3306"
export BACKUP_DEST_NAME="mysql"
export BACKUP_DEST_PATH="test"

./bin/db-backup

cat /tmp/dump.sql