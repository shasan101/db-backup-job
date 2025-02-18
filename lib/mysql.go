package lib

import (
	"compress/gzip"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DbBackup defines the interface for database backup operations.
type DbBackup interface {
	GetDbName() string
	MakeBackup(outputPath string) error
	CompressBackup(backupPath string) (string, error)
	BackupWithCompression(outputPath string) (string, error)
}

// MySQLBackup is an implementation of DbBackup for MySQL databases.
type MySQLBackup struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     string
}

func InitMySql(db, user, pass, host, port string) *MySQLBackup {
	return &MySQLBackup{
		DBName:   db,
		User:     user,
		Password: pass,
		Host:     host,
		Port:     port,
	}
}

// GetDbName returns the database name.
func (m *MySQLBackup) GetDbName() string {
	return m.DBName
}

// MakeBackup connects to MySQL and writes database content as CSV.
func (m *MySQLBackup) MakeBackup(outputPath string) error {
	// Open connection to MySQL database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.User, m.Password, m.Host, m.Port, m.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Query all tables in the database
	tables, err := m.getTables(db)
	if err != nil {
		return fmt.Errorf("failed to get tables: %v", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Iterate through tables and dump data
	for _, table := range tables {
		if err := m.dumpTable(db, table, writer); err != nil {
			return fmt.Errorf("failed to dump table %s: %v", table, err)
		}
	}

	fmt.Println("Backup completed:", outputPath)
	return nil
}

// CompressBackup compresses a file using gzip.
func (m *MySQLBackup) CompressBackup(backupPath string) (string, error) {
	compressedPath := backupPath + ".gz"

	// Open the original file
	inFile, err := os.Open(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to open backup file: %v", err)
	}
	defer inFile.Close()

	// Create compressed file
	outFile, err := os.Create(compressedPath)
	if err != nil {
		return "", fmt.Errorf("failed to create compressed file: %v", err)
	}
	defer outFile.Close()

	// Compress data using gzip
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	_, err = gzipWriter.ReadFrom(inFile)
	if err != nil {
		return "", fmt.Errorf("failed to compress backup: %v", err)
	}

	fmt.Println("Compression completed:", compressedPath)
	return compressedPath, nil
}

// BackupWithCompression performs backup and compression in one step.
func (m *MySQLBackup) BackupWithCompression(outputPath string) (string, error) {
	if err := m.MakeBackup(outputPath); err != nil {
		return "", err
	}
	return m.CompressBackup(outputPath)
}

// getTables retrieves all table names from the MySQL database.
func (m *MySQLBackup) getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
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
func (m *MySQLBackup) dumpTable(db *sql.DB, table string, writer *csv.Writer) error {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", table))
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
