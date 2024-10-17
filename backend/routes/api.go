package routes

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "strconv"
    "backend/db"
    "gorm.io/gorm"
    "log"
	"sync"
	"fmt"
	"regexp"
)

func maskSensitiveInfo(dsn string) string {
    // This is a simple implementation. Adjust based on your DSN format
    return regexp.MustCompile(`:[^:@]+@`).ReplaceAllString(dsn, ":***@")
}

// SetupRoutes sets up all the routes for the API.
func SetupRoutes(router *mux.Router, dbConn *gorm.DB) {
    // User endpoints
    setupUserRoutes(router, dbConn)

    // Transaction endpoints
    setupTransactionRoutes(router, dbConn)

    // Session endpoints
    setupSessionRoutes(router, dbConn)

    // Database connection endpoints
    setupDatabaseRoutes(router)

	setUpTestRoute(router)
}

// setupUserRoutes defines the user-related API routes.
func setupUserRoutes(router *mux.Router, dbConn *gorm.DB) {
    // Route to create a new user
    router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        var user db.User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        if err := db.CreateUser(dbConn, &user); err != nil {
            log.Println("Error creating user:", err)
            http.Error(w, "Failed to create user", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }).Methods("POST")

    // Route to get user details
    router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        user, err := db.GetUser(dbConn, userID)
        if err != nil {
            log.Println("Error fetching user:", err)
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(user)
    }).Methods("GET")

    // Route to update user details
    router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        user, err := db.GetUser(dbConn, userID)
        if err != nil {
            log.Println("Error fetching user for update:", err)
            http.Error(w, "User not found", http.StatusNotFound)
            return
        }

        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        if err := db.UpdateUser(dbConn, user); err != nil {
            log.Println("Error updating user:", err)
            http.Error(w, "Failed to update user", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(user)
    }).Methods("PUT")

    // Route to delete a user
    router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        userID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }

        if err := db.DeleteUser(dbConn, userID); err != nil {
            log.Println("Error deleting user:", err)
            http.Error(w, "Failed to delete user", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }).Methods("DELETE")

    // Route to get all users
    router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        users, err := db.GetUsers(dbConn)
        if err != nil {
            log.Println("Error fetching users:", err)
            http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(users)
    }).Methods("GET")
}

// setupTransactionRoutes defines the transaction-related API routes.
func setupTransactionRoutes(router *mux.Router, dbConn *gorm.DB) {
    // Route to create a new transaction
    router.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
        var transaction db.Transaction
        if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        if err := db.CreateTransaction(dbConn, &transaction); err != nil {
            log.Println("Error creating transaction:", err)
            http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(transaction)
    }).Methods("POST")

    // Route to get a transaction by ID
    router.HandleFunc("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        transactionID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
            return
        }

        transaction, err := db.GetTransaction(dbConn, transactionID)
        if err != nil {
            log.Println("Error fetching transaction:", err)
            http.Error(w, "Transaction not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(transaction)
    }).Methods("GET")

    // Route to update a transaction
    router.HandleFunc("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        transactionID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
            return
        }

        transaction, err := db.GetTransaction(dbConn, transactionID)
        if err != nil {
            log.Println("Error fetching transaction for update:", err)
            http.Error(w, "Transaction not found", http.StatusNotFound)
            return
        }

        if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        if err := db.UpdateTransaction(dbConn, transaction); err != nil {
            log.Println("Error updating transaction:", err)
            http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(transaction)
    }).Methods("PUT")

    // Route to delete a transaction
    router.HandleFunc("/transactions/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        transactionID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
            return
        }

        if err := db.DeleteTransaction(dbConn, transactionID); err != nil {
            log.Println("Error deleting transaction:", err)
            http.Error(w, "Failed to delete transaction", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }).Methods("DELETE")
}

// setupSessionRoutes defines the session-related API routes.
func setupSessionRoutes(router *mux.Router, dbConn *gorm.DB) {
    // Route to create a new session
    router.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
        var session db.Session
        if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        if err := db.CreateSession(dbConn, &session); err != nil {
            log.Println("Error creating session:", err)
            http.Error(w, "Failed to create session", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(session)
    }).Methods("POST")

    // Route to get a session by ID
    router.HandleFunc("/sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        sessionID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid session ID", http.StatusBadRequest)
            return
        }

        session, err := db.GetSession(dbConn, sessionID)
        if err != nil {
            log.Println("Error fetching session:", err)
            http.Error(w, "Session not found", http.StatusNotFound)
            return
        }

        json.NewEncoder(w).Encode(session)
    }).Methods("GET")

    // Route to update a session
    router.HandleFunc("/sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        sessionID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid session ID", http.StatusBadRequest)
            return
        }

        session, err := db.GetSession(dbConn, sessionID)
        if err != nil {
            log.Println("Error fetching session for update:", err)
            http.Error(w, "Session not found", http.StatusNotFound)
            return
        }

        if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        if err := db.UpdateSession(dbConn, session); err != nil {
            log.Println("Error updating session:", err)
            http.Error(w, "Failed to update session", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(session)
    }).Methods("PUT")

    // Route to delete a session
    router.HandleFunc("/sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        sessionID, err := strconv.Atoi(vars["id"])
        if err != nil {
            http.Error(w, "Invalid session ID", http.StatusBadRequest)
            return
        }

        if err := db.DeleteSession(dbConn, sessionID); err != nil {
            log.Println("Error deleting session:", err)
            http.Error(w, "Failed to delete session", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusNoContent)
    }).Methods("DELETE")

    // Example: Additional session-related routes
    
}

var (
    userDB *gorm.DB  // Global variable to hold the database connection
    mu     sync.Mutex // Mutex to prevent race conditions when accessing the global variable
)

// ConnectDatabaseRequest is the struct for the database connection request
type ConnectDatabaseRequest struct {
    DSN    string `json:"dsn"`    // Data Source Name (connection string)
    Driver string `json:"driver"` // Database driver (e.g., postgres, mysql, sqlite, sqlserver)
    Port   string `json:"port"`
}

// ConnectDatabaseResponse is the struct for the response after connecting to the database
type ConnectDatabaseResponse struct {
    Message string `json:"message"`
    Success bool `json:"success"` // Success or error message
}

// QueryRequest is the struct for the SQL query request
type QueryRequest struct {
    SQLQuery string `json:"sql_query"` // SQL query to be executed
}

// QueryResponse is the struct for the SQL query response
type QueryResponse struct {
    Result []map[string]interface{} `json:"result"` // Query result as a slice of maps
    Error  string                   `json:"error,omitempty"` // Optional error message
}


type PreviewRequest struct {
    TableName string `json:"table_name"`
}

type PreviewResponse struct {
    Success bool            `json:"success"`
    Message string          `json:"message"`
    Columns []string        `json:"columns"`
    Rows    [][]interface{} `json:"rows"`
}

// setupDatabaseRoutes defines the database connection-related API routes.
func setupDatabaseRoutes(router *mux.Router) {
    router.HandleFunc("/database/connect", func(w http.ResponseWriter, r *http.Request) {
        log.Println("=== Starting database connection request ===")

        var req ConnectDatabaseRequest
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validate input
        if req.DSN == "" || req.Driver == "" {
            http.Error(w, "DSN and Driver are required", http.StatusBadRequest)
            return
        }

        // Construct the full DSN if port is provided


        new_DB, err := db.ConnectUserDatabase(req.DSN, req.Driver)
        if err != nil {
            log.Printf("Error connecting to database: %v", err)
            http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
            return
        }
        userDB = new_DB

        response := ConnectDatabaseResponse{
            Success: true,
            Message: "Successfully connected to the database",
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)

        log.Println("=== Database connection request completed successfully ===")
    }).Methods("POST")

    router.HandleFunc("/database/preview", func(w http.ResponseWriter, r *http.Request) {
        log.Println("=== Starting database preview request ===")

        var req PreviewRequest
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        if req.TableName == "" {
            http.Error(w, "Table name is required", http.StatusBadRequest)
            return
        }

        mu.Lock()
        defer mu.Unlock()

        if userDB == nil {
            http.Error(w, "No active database connection", http.StatusInternalServerError)
            return
        }

        // Construct the query to get the first 10 rows
        query := fmt.Sprintf("SELECT * FROM %s LIMIT 10", req.TableName)
         
        rows, err := db.ExecuteSQLQuery(userDB, query) 
        if err != nil {
            log.Printf("Error querying database: %v", err)
            http.Error(w, "Failed to query database", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        columns, err := rows.Columns()
        if err != nil {
            log.Printf("Error getting columns: %v", err)
            http.Error(w, "Failed to get column information", http.StatusInternalServerError)
            return
        }

        var result [][]interface{}
        for rows.Next() {
            row := make([]interface{}, len(columns))
            rowPointers := make([]interface{}, len(columns))
            for i := range row {
                rowPointers[i] = &row[i]
            }

            err := rows.Scan(rowPointers...)
            if err != nil {
                log.Printf("Error scanning row: %v", err)
                continue
            }

            result = append(result, row)
        }

        response := PreviewResponse{
            Success: true,
            Message: "Successfully retrieved preview data",
            Columns: columns,
            Rows:    result,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)

        log.Println("=== Database preview request completed successfully ===")
    }).Methods("GET")

    // Route to allow users to execute SQL queries on their connected database
    router.HandleFunc("/database/query", func(w http.ResponseWriter, r *http.Request) {
        var req QueryRequest

        // Decode the JSON request body into the struct
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            log.Println("Failed to decode request body for query:", err)
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        // Log the received query request
        log.Printf("Received SQL query request: %+v\n", req)

        // Ensure the SQL query is not empty
        if req.SQLQuery == "" {
            log.Println("SQL query is empty")
            http.Error(w, "SQL query cannot be empty", http.StatusBadRequest)
            return
        }

        // Ensure the database connection has been established
        mu.Lock()
        db := userDB
        mu.Unlock()

        if db == nil {
            log.Println("No database connection found")
            http.Error(w, "No database connection found. Please connect to a database first.", http.StatusInternalServerError)
            return
        }

        // Execute the SQL query
        var result []map[string]interface{}
        err := db.Raw(req.SQLQuery).Scan(&result).Error
        if err != nil {
            log.Println("Failed to execute SQL query:", err)
            response := QueryResponse{
                Error: err.Error(),
            }
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(response)
            return
        }

        // Send back the result
        response := QueryResponse{
            Result: result,
        }
        log.Println("SQL query executed successfully, returning result")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(response)
    }).Methods("POST")
}

func setUpTestRoute(router *mux.Router) {
    router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "message": "API is working",
        })
    }).Methods("GET")

    router.HandleFunc("/test_db", func(w http.ResponseWriter, r *http.Request) {
        dsn := "postgresql://misaalsingh:Jupiter_52@localhost:5432/test_database"
        driver := "postgres"
        new_db, err := db.ConnectUserDatabase(dsn, driver)
        if err != nil {
            http.Error(w, fmt.Sprintf("Failed to connect to database: %v", err), http.StatusInternalServerError)
            return
        }
        json.NewEncoder(w).Encode(map[string]string{
            "message": "Connection is successful",
        })

        sqlDB, _ := new_db.DB()
        defer sqlDB.Close()
    }).Methods("GET")
}