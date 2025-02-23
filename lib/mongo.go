package lib

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDbBackup struct {
	BackupFile  string `json:"backup_path,omitempty"`
	Collection  string `json:"collection_name,omitempty"`
	DbName      string `json:"db_name,omitempty"`
	MongoClient *mongo.Client
}

func InitMonogDb(db, collection, user, pass, host, port string) (*MongoDbBackup, error) {
	path := MakeBackupName()

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?compressors=snappy,zlib,zstd", user, pass, host, port)

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}

	return &MongoDbBackup{
		Collection:  collection,
		DbName:      db,
		MongoClient: client,
		BackupFile:  path,
	}, nil
}

func (m *MongoDbBackup) GetCollection() string {
	return m.Collection
}

func (m *MongoDbBackup) MakeBackup(c context.Context) error {
	collection := m.MongoClient.Database(m.DbName).Collection(m.Collection)

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(ctx)
	var results interface{}
	cur.All(ctx, &results)
	if err != nil {
		return err
	}
	err = WriteToFile(results.([]byte), "/tmp/"+m.BackupFile)
	return err
}

func (m *MongoDbBackup) CompressBackup() (string, error) {
	// Read the original file
	inFileData, err := os.ReadFile(dumpPath)
	if err != nil {
		return "", fmt.Errorf("failed to open backup file: %v", err)
	}

	// Create compressed file
	outFile, err := os.Create("/tmp/" + m.BackupFile)
	if err != nil {
		return "", fmt.Errorf("failed to create compressed file: %v", err)
	}
	defer outFile.Close()

	// Compress data using gzip
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	_, err = gzipWriter.Write(inFileData)
	if err != nil {
		return "", fmt.Errorf("failed to compress backup: %v", err)
	}

	log.Println("Compression completed:", m.BackupFile)
	return m.BackupFile, nil
}

func (m *MongoDbBackup) BackupWithCompression(c context.Context) (string, error) {
	if err := m.MakeBackup(c); err != nil {
		return "", err
	}
	return m.CompressBackup()
}

func WriteToFile(data []byte, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	// close f on exit and check for its returned error
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	w := bufio.NewWriter(f)
	_, err = w.Write(data)
	if err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}
	return nil
}
