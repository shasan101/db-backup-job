package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/shasan101/db-backup-job/lib/uploaders"
)

func MakeBackupName() string {
	currentTime := time.Now()
	return fmt.Sprintf("backup-%s.tar.gz", currentTime.Format("02-01-2006-15_04"))
}

func GetStorageLayer(layerName string) func(context.Context, string, string) (uploaders.UploadBackup, error) {

	storageTypes := map[string]func(context.Context, string, string) (uploaders.UploadBackup, error){
		"aws": func(c context.Context, bucket, key string) (uploaders.UploadBackup, error) {
			return uploaders.InitAwsBucket(c, bucket, key)
		},
		"gcs": func(c context.Context, bucket, key string) (uploaders.UploadBackup, error) {
			return uploaders.InitGcpBucket(c, bucket, key)
		},
		"local": func(c context.Context, bucket, key string) (uploaders.UploadBackup, error) {
			return uploaders.InitLocal(c, bucket, key)
		},
	}
	return storageTypes[layerName]
}
