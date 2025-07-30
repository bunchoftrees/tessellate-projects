package db

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

)

// Replace the global DB variable
var DB *Database

// Add the wrapper struct
type Database struct {
    *gorm.DB
}

func ptr[T any](v T) *T {
    return &v
}

// Update InitDB to initialize the custom Database
func InitDB() *Database {
    var err error
    gormDB, err := gorm.Open(sqlite.Open("tessellate_projects.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    DB = &Database{gormDB}

    log.Println("Database connected!")

    // Auto-migrate schema for all models
    err = DB.AutoMigrate(
        &User{},
        &Project{},
        &Requirement{},
        &AuditTask{},
        &Issue{},
        &Client{},
    )
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    // Seed sample User data if DB is empty
    var userCount int64
    DB.Model(&User{}).Count(&userCount)
    if userCount == 0 {
        DB.Create(&User{Name: "Alice", Email: "alice@example.com", Role: "CONSULTANT"})
    }

    // Seed sample Project data if DB is empty
    var count int64
    DB.Model(&Project{}).Count(&count)
    if count == 0 {
        DB.Create(&Project{Name: "Demo Project", ClientName: "Demo Client"})
    }

    // Seed sample Client data if DB is empty
    DB.Model(&Client{}).Count(&count)
    if count == 0 {
        DB.Create(&Client{Name: "Demo Client Org"})
    }

    // Seed sample Requirement data if DB is empty
    DB.Model(&Requirement{}).Count(&count)
    if count == 0 {
        DB.Create(&Requirement{
            Text:     "Must support single sign-on",
            Category: ptr("Authentication"),
            Status:   "DRAFT",
        })
    }

    // Seed sample AuditTask data if DB is empty
    DB.Model(&AuditTask{}).Count(&count)
    if count == 0 {
        DB.Create(&AuditTask{
            Text:   "Check login audit",
            Notes:  ptr("Review all login-related requirements"),
            Status: "PENDING",
        })
    }

    // Seed sample Issue data if DB is empty
    DB.Model(&Issue{}).Count(&count)
    if count == 0 {
        DB.Create(&Issue{
            Title:       "Login fails on Safari",
            Description: ptr("Users report login page broken in Safari"),
        })
    }
    return DB
}

// Add custom query methods
func (db *Database) GetRequirementsByClientID(clientID string) ([]Requirement, error) {
    var reqs []Requirement
    err := db.Where("client_id = ?", clientID).Find(&reqs).Error
    return reqs, err
}

func (db *Database) GetAllRequirements() ([]Requirement, error) {
    var reqs []Requirement
    err := db.Find(&reqs).Error
    return reqs, err
}

// GetRequirementsByStatus returns all requirements with the given status.
func (db *Database) GetRequirementsByStatus(clientID string, status string) ([]Requirement, error) {
    var reqs []Requirement
    err := db.Where("client_id = ? AND status = ?", clientID, status).Find(&reqs).Error
    return reqs, err
}

// GetUserByID retrieves a User by ID.
func (db *Database) GetUserByID(userID string) (*User, error) {
    var user User
    err := db.First(&user, "id = ?", userID).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (db *Database) GetAllUsers() ([]User, error) {
    var users []User
    err := db.Find(&users).Error
    return users, err
}
// GetUsersByClient retrieves all users associated with a given client ID.
func (db *Database) GetUsersByClient(clientID string) ([]User, error) {
    var users []User
    err := db.Where("client_id = ?", clientID).Find(&users).Error
    return users, err
}

// AuthPayload represents the authentication payload containing a token and user.

