package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	lib "github.com/shasan101/db-backup-jobs/lib"
)

type DbEnv struct {
	DbName     string
	ConnString string
	DbFilePath string
	DestName   string
	DestPath   string
}

func main() {
	// programLevel := new(slog.LevelVar)
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))

	var dbEnv DbEnv
	err := envconfig.Process("BACKUP", &dbEnv)
	if err != nil {
		log.Fatal(err.Error())
	}
	switch dbEnv.DestPath {
	case "mysql":
		// call the mysql handler
		mysql := lib.InitMySql()
		backup := mysql.BackupWithCompression()
		err := mysql.dumpBackup()
	default:
		log.Fatalf("No DB type found in relation to the provided input: %v", dbEnv.DestName)
	}

}
