package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
	lib "github.com/shasan101/db-backup-job/lib"
)

type DbEnv struct {
	SourceName     string `json:"source_name" envconfig:"SOURCE_NAME" required:"true"` // mysql, mongo, etc
	DbName         string `json:"db_name" envconfig:"DB_NAME" required:"true"`
	User           string `json:"user" envconfig:"USERNAME" required:"true"`
	Password       string `json:"password" envconfig:"PASSWORD" required:"true"`
	Host           string `json:"host" envconfig:"HOST" required:"true"`
	Port           string `json:"port" envconfig:"PORT" required:"true"`
	DestType       string `json:"dest_type" envconfig:"DEST_TYPE" default:"local" required:"true"` // local, s3, gcs
	DestName       string `json:"dest_name" envconfig:"DEST_NAME" required:"true"`                 // storage bucket name
	DestPath       string `json:"dest_path" envconfig:"DEST_PATH" required:"true"`                 // specifiy the path to store the backup in the bucket
	CollectionName string `json:"collection_name" envconfig:"COLLECTION" ignored:"true"`
}

func main() {

	ctx := context.Background()
	programLevel := new(slog.LevelVar)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))

	slog.SetDefault(logger)

	var dbEnv DbEnv
	err := envconfig.Process("BACKUP", &dbEnv)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("read env vars are: %+v\n", dbEnv)
	switch dbEnv.SourceName {
	case "mysql":
		// call the mysql handler
		mysql, err := lib.InitMySql(dbEnv.DbName, dbEnv.User, dbEnv.Password, dbEnv.Host, dbEnv.Port)
		if err != nil {
			log.Fatalf("mysql init failure: %v\n", err)
		}
		output, err := mysql.BackupWithCompression(ctx)
		if err != nil {
			log.Fatalf("backup and compression failure: %v\n", err)
		}
		storageFunc := lib.GetStorageLayer(dbEnv.DestType)
		if err != nil {
			log.Fatalf("storage layer failure: %v\n", err)
		}
		storageClient, err := storageFunc(ctx, dbEnv.DestName, dbEnv.DestPath)
		if err != nil {
			log.Fatalf("storage client failure: %v\n", err)
		}
		err = storageClient.UploadBackup(ctx, output)
		if err != nil {
			log.Fatalf("upload failure: %v\n", err)
		}
	case "mongodb":
		// call the mysql handler
		mongo, err := lib.InitMonogDb(dbEnv.DbName, dbEnv.CollectionName, dbEnv.User, dbEnv.Password, dbEnv.Host, dbEnv.Port)
		if err != nil {
			log.Fatalf("mysql init failure: %v\n", err)
		}
		output, err := mongo.BackupWithCompression(ctx)
		if err != nil {
			log.Fatalf("backup and compression failure: %v\n", err)
		}
		storageFunc := lib.GetStorageLayer(dbEnv.DestType)
		if err != nil {
			log.Fatalf("storage layer failure: %v\n", err)
		}
		storageClient, err := storageFunc(ctx, dbEnv.DestName, dbEnv.DestPath)
		if err != nil {
			log.Fatalf("storage client failure: %v\n", err)
		}
		err = storageClient.UploadBackup(ctx, output)
		if err != nil {
			log.Fatalf("upload failure: %v\n", err)
		}
	default:
		log.Fatalf("No DB type found in relation to the provided input: %v", dbEnv.DestName)
	}
}
