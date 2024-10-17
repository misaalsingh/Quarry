package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"regexp"
	"github.com/jackc/pgx/v4"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func maskSensitiveInfo(dsn string) string {
    // This is a simple implementation. Adjust based on your DSN format
    return regexp.MustCompile(`:[^:@]+@`).ReplaceAllString(dsn, ":***@")
}

// SQLQueryResponse struct to handle the response from Python text-to-SQL service
type SQLQueryResponse struct {
	SQLQuery string `json:"sql_query"`
}

// ConnectDB connects to CockroachDB and prints the current time.
func ConnectDB() (*gorm.DB, error) {
	// Example DSN, adjust with your CockroachDB connection details
	dsn := "postgresql://Missy:kuCLZql3L7g6lwBVj_kBvg@spector-cluster-2997.j77.aws-us-east-1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	ctx := context.Background()
	
	// Connect using pgx
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err // Return the error
	}
	defer conn.Close(ctx) // Close connection after checking for errors

	var now time.Time
	err = conn.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
		return nil, err // Return the error
	}

	fmt.Println("Current time from DB:", now)

	// Connect using GORM for further operations
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to open GORM connection:", err)
		return nil, err
	}

	return db, nil
}

// ConnectUserDatabase dynamically connects to the specified database using GORM.
func ConnectUserDatabase(dsn string, driver string) (*gorm.DB, error) {
    log.Printf("=== Attempting to connect to %s database ===", driver)
    
    var db *gorm.DB
    var err error
    
    maskedDSN := maskSensitiveInfo(dsn)
    log.Printf("Using DSN: %s", maskedDSN)
    
    switch driver {
    case "postgres":
        log.Println("Using PostgreSQL driver")
        db, err = gorm.Open(postgres.New(postgres.Config{
            DSN:                  dsn,
            PreferSimpleProtocol: true,
        }), &gorm.Config{})
    case "mysql":
        log.Println("Using MySQL driver")
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    case "sqlite":
        log.Println("Using SQLite driver")
        db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
    case "sqlserver":
        log.Println("Using SQL Server driver")
        db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
    default:
        log.Printf("❌ Unsupported database driver: %s", driver)
        return nil, fmt.Errorf("unsupported database driver: %s", driver)
    }
    
    if err != nil {
        log.Printf("❌ Failed to connect to database: %v", err)
        return nil, fmt.Errorf("failed to connect to %s: %v", driver, err)
    }
    
    // Test the connection
    sqlDB, err := db.DB()
    if err != nil {
        log.Printf("❌ Failed to get underlying *sql.DB: %v", err)
        return nil, fmt.Errorf("failed to get database instance: %v", err)
    }
    
    err = sqlDB.Ping()
    if err != nil {
        log.Printf("❌ Failed to ping database: %v", err)
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }
    
    log.Println("✅ Successfully connected to database")
    
    // Log some connection pool settings
    log.Printf("Connection pool settings:")
    log.Printf("  Max idle connections: %d", sqlDB.Stats().Idle)
    log.Printf("  Max open connections: %d", sqlDB.Stats().OpenConnections)
    
    return db, nil
}


// Call the Python text-to-SQL service to convert text query into SQL.
func convertTextToSQL(textQuery string) (string, error) {
	// Prepare the request payload
	data := map[string]string{"text_query": textQuery}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	// Call the Python service
	resp, err := http.Post("http://localhost:5000/api/query", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response SQLQueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.SQLQuery, nil
}

// Execute the SQL query against the connected database.
func ExecuteSQLQuery(db *gorm.DB, sqlQuery string) (*sql.Rows, error) {
	// Use GORM's Raw method to execute the query and get the result as *sql.Rows
	rows, err := db.Raw(sqlQuery).Rows()
	if err != nil {
		log.Println("Error executing SQL query:", err)
		return nil, err
	}

	// Print the columns as an example
	columns, err := rows.Columns()
	if err != nil {
		log.Println("Error fetching columns:", err)
		rows.Close() // Make sure to close the rows in case of error
		return nil, err
	}
	fmt.Println("Query Result Columns:", columns)

	return rows, nil
}
