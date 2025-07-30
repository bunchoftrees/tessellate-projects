package db

import (
	"gorm.io/gorm"
)

type Role string

const (
    RoleAdmin      Role = "ADMIN"
    RoleConsultant Role = "CONSULTANT"
    RoleClient     Role = "CLIENT"
)

type RequirementStatus string

const (
	RequirementStatusMet    RequirementStatus = "MET"
	RequirementStatusNotMet RequirementStatus = "NOT_MET"
)

type User struct {
    gorm.Model
    Name     string
    Email    string
    Password string
    Role     Role
    Projects []*Project `gorm:"many2many:project_users"`
    ClientID *uint
    Client   *Client
}

type Project struct {
    gorm.Model
    Name         string
    ClientName   string
    Status       string `gorm:"type:VARCHAR(20);default:'NEW'"`
    Users        []*User `gorm:"many2many:project_users"`
    ClientID     *uint
    Client       *Client
    Requirements []*Requirement
}

type Requirement struct {
    gorm.Model
    ProjectID  uint
    Project    *Project
    Text       string
    Category   *string
    Status     RequirementStatus
    AuditTasks []*AuditTask
}

type AuditTask struct {
    gorm.Model
    RequirementID uint
    Requirement   *Requirement
    Text          string
    Status        string
    Notes         *string
    Issue         *Issue
}

type Issue struct {
    gorm.Model
    AuditTaskID  uint
    AuditTask    *AuditTask
    Title        string
    Description  *string
    Priority     *string
    Phase        *string
    EstimateHrs  *int
    Status       string
    Type         string
}

type Client struct {
    gorm.Model
    Name         string
    Industry     *string
    ContactName  *string
    ContactEmail *string
    Users        []*User
    Projects     []*Project
}

type AuthPayload struct {
    Token string
    User  *User
}