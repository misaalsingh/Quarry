package db

import (
    "encoding/csv"
    "io"
    "gorm.io/gorm"
	"log"
	"encoding/json"
)


// ImportCSV reads a CSV file and returns its contents as a JSON object

func ImportCSV(r io.Reader) ([]byte, error) {
	reader := csv.NewReader(r)
	// Read the header row
	headers, err := reader.Read()
	if err != nil {
		log.Println("Error reading header row:", err)
		return nil, err
	}

	var records []map[string]string

	// Iterate over CSV rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading CSV:", err)
			continue
		}

		// Create a map for the record
		recordMap := make(map[string]string)
		for i, value := range record {
			recordMap[headers[i]] = value // Map header to value
		}
		records = append(records, recordMap) // Add the map to the records slice
	}

	// Convert records to JSON
	jsonData, err := json.Marshal(records)
	if err != nil {
		log.Println("Error marshaling records to JSON:", err)
		return nil, err
	}

	return jsonData, nil // Return the JSON data
}


// CreateUser creates a new user in the database.
func CreateUser(db *gorm.DB, user *User) error {
    if err := db.Create(user).Error; err != nil {
        log.Println("Error creating user:", err)
        return err
    }
    return nil
}

// GetUser retrieves a user by their ID.
func GetUser(db *gorm.DB, userID int) (*User, error) {
    var user User
    if err := db.Preload("Transactions").Preload("Sessions").First(&user, userID).Error; err != nil {
        log.Println("Error fetching user:", err)
        return nil, err
    }
    return &user, nil
}

// UpdateUser updates an existing user's details.
func UpdateUser(db *gorm.DB, user *User) error {
    if err := db.Save(user).Error; err != nil {
        log.Println("Error updating user:", err)
        return err
    }
    return nil
}

// DeleteUser removes a user and their related data from the database.
func DeleteUser(db *gorm.DB, userID int) error {
    if err := db.Delete(&User{}, userID).Error; err != nil {
        log.Println("Error deleting user:", err)
        return err
    }
    return nil
}

func GetTransactions(dbConn *gorm.DB) ([]Transaction, error) {
	var transactions []Transaction
	if err := dbConn.Find(&transactions).Error; err != nil {
        return nil, err
    }
    return transactions, nil
}

// CreateTransaction creates a new transaction for a user.
func CreateTransaction(db *gorm.DB, transaction *Transaction) error {
    if err := db.Create(transaction).Error; err != nil {
        log.Println("Error creating transaction:", err)
        return err
    }
    return nil
}

func GetTransaction(db *gorm.DB, transactionID int) (*Transaction, error) {
    var transaction Transaction
    if err := db.First(&transaction, transactionID).Error; err != nil {
        log.Println("Error fetching session:", err)
        return nil, err
    }
    return &transaction, nil
}

// UpdateTransaction updates an existing transaction's details.
func UpdateTransaction(db *gorm.DB, transaction *Transaction) error { // Use *Transaction instead of transaction
    if err := db.Save(transaction).Error; err != nil {
        log.Println("Error updating transaction:", err)
        return err
    }
    return nil
}

// DeleteTransaction removes a transaction by its ID.
func DeleteTransaction(db *gorm.DB, transactionID int) error { // Use transactionID instead of transaction
    if err := db.Delete(&Transaction{}, transactionID).Error; err != nil {
        log.Println("Error deleting transaction:", err)
        return err
    }
    return nil
}


// GetUserTransactions retrieves all transactions for a specific user.
func GetUserTransactions(db *gorm.DB, userID int) ([]Transaction, error) {
    var transactions []Transaction
    if err := db.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
        log.Println("Error fetching transactions:", err)
        return nil, err
    }
    return transactions, nil
}

// CreateSession creates a new session for a user.
func CreateSession(db *gorm.DB, session *Session) error {
    if err := db.Create(session).Error; err != nil {
        log.Println("Error creating session:", err)
        return err
    }
    return nil
}

// GetSession retrieves a session by its ID.
func GetSession(db *gorm.DB, sessionID int) (*Session, error) {
    var session Session
    if err := db.First(&session, sessionID).Error; err != nil {
        log.Println("Error fetching session:", err)
        return nil, err
    }
    return &session, nil
}

// DeleteSession removes a session by its ID.
func DeleteSession(db *gorm.DB, sessionID int) error {
    if err := db.Delete(&Session{}, sessionID).Error; err != nil {
        log.Println("Error deleting session:", err)
        return err
    }
    return nil
}

// CreateOrUpdateSession creates a new session or updates an existing session if it exists.
func CreateOrUpdateSession(db *gorm.DB, session *Session) error {
    existingSession := Session{}
    if err := db.Where("user_id = ? AND token = ?", session.UserID, session.Token).First(&existingSession).Error; err == nil {
        // Session exists, update it
        session.ID = existingSession.ID
        return UpdateSession(db, session)
    }
    // Create new session
    return CreateSession(db, session)
}

// UpdateSession updates an existing session's details.
func UpdateSession(db *gorm.DB, session *Session) error {
    if err := db.Save(session).Error; err != nil {
        log.Println("Error updating session:", err)
        return err
    }
    return nil
}

// GetUsers retrieves all users from the database.
func GetUsers(db *gorm.DB) ([]User, error) {
    var users []User
    if err := db.Preload("Transactions").Preload("Sessions").Find(&users).Error; err != nil {
        log.Println("Error fetching users:", err)
        return nil, err
    }
    return users, nil
}
