# db-backup-job

<u><b>NOTE:</b></u> Practice and/or Just for Fun. Any, and all, feedback is welcome.

The `db-backup-job` project is designed to facilitate the process of creating and managing database backups. This project supports multiple database types and storage solutions, making it a versatile tool for database administrators and developers.

## Features

- **Database Support**: Currently supports MySQL and MongoDB databases.
- **Storage Solutions**: Supports local storage, AWS S3, and Google Cloud Storage (GCS).
- **Backup Compression**: Compresses backups to save storage space.
- **Environment Configuration**: Uses environment variables for configuration, making it easy to integrate with different environments.

## Getting Started

### Prerequisites

- Go 1.16 or later
- Access to the databases you want to back up
- Access to the storage solutions you want to use (local, AWS S3, GCS)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/shasan101/db-backup-job.git
    cd db-backup-job
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

### Configuration

Configure the environment variables required for the backup job. You can set these variables in your environment or use a `.env` file.

```sh
export BACKUP_SOURCE_NAME=mysql
export BACKUP_DB_NAME=your_db_name
export BACKUP_USERNAME=your_username
export BACKUP_PASSWORD=your_password
export BACKUP_HOST=your_host
export BACKUP_PORT=your_port
export BACKUP_DEST_TYPE=s3
export BACKUP_DEST_NAME=your_bucket_name
export BACKUP_DEST_PATH=your_backup_path
```


### Build
use the Makefile command: ```make build```

### Run
./bin/db-backup