package lib

import (
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DbBackup defines the interface for database backup operations.
type DbBackup interface {
	GetDbName() string
	MakeBackup(outputPath string) error
	CompressBackup(backupPath string) (string, error)
	BackupWithCompression(outputPath string) (string, error)
}

const dumpPath string = "/tmp/dump.sql"

// MySQLBackup is an implementation of DbBackup for MySQL databases.
type MySQLBackup struct {
	BackupFile string `json:"backup_path,omitempty"`
	SqlClient  *sql.DB
	DbName     string `json:"db_name,omitempty"`
}

func InitMySql(db, user, pass, host, port string) (*MySQLBackup, error) {
	path := MakeBackupName()
	// Open connection to MySQL database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, db)
	db_conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &MySQLBackup{
		DbName:     db,
		BackupFile: path,
		SqlClient:  db_conn,
	}, nil
}

// GetDbName returns the database name.
func (m *MySQLBackup) GetDbName() string {
	return m.DbName
}

// MakeBackup connects to MySQL and writes database content as CSV.
func (m *MySQLBackup) MakeBackup(c context.Context) error {
	// Query all tables in the database
	defer m.SqlClient.Close()
	tables, err := m.getTables(c)
	if err != nil {
		return fmt.Errorf("failed to get tables: %v", err)
	}

	// Create output file
	outFile, err := os.Create(dumpPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Iterate through tables and dump data
	for _, table := range tables {
		if err := m.dumpTable(c, table, writer); err != nil {
			return fmt.Errorf("failed to dump table %s: %v", table, err)
		}
	}

	log.Printf("Backup completed: %s", dumpPath)
	return nil
}

// CompressBackup compresses a file using gzip.
func (m *MySQLBackup) CompressBackup() (string, error) {

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

// BackupWithCompression performs backup and compression in one step.
func (m *MySQLBackup) BackupWithCompression(c context.Context) (string, error) {
	if err := m.MakeBackup(c); err != nil {
		return "", err
	}
	return m.CompressBackup()
}

// getTables retrieves all table names from the MySQL database.
func (m *MySQLBackup) getTables(c context.Context) ([]string, error) {
	db := m.SqlClient
	ctx, cancel := context.WithTimeout(c, 50*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

// dumpTable writes table data to CSV format.
func (m *MySQLBackup) dumpTable(c context.Context, table string, writer *csv.Writer) error {
	db := m.SqlClient
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s", table))
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Write column names as header
	if err := writer.Write(columns); err != nil {
		return err
	}

	// Write data rows
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		rowData := make([]string, len(columns))
		for i, val := range values {
			if val != nil {
				rowData[i] = fmt.Sprintf("%v", val)
			} else {
				rowData[i] = "NULL"
			}
		}

		if err := writer.Write(rowData); err != nil {
			return err
		}
	}
	return nil
}
