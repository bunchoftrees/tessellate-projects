package models

import "time"

// API Response models - these are what we return to clients
// They're separate from database models for better API design

type ProjectResponse struct {
	ID           uint                    `json:"id"`
	Name         string                  `json:"name"`
	ClientName   string                  `json:"clientName"`
	Status       string                  `json:"status"`
	ClientID     *uint                   `json:"clientId,omitempty"`
	Client       *ClientResponse         `json:"client,omitempty"`
	Users        []UserResponse          `json:"users,omitempty"`
	Requirements []RequirementResponse   `json:"requirements,omitempty"`
	Issues       []IssueResponse         `json:"issues,omitempty"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
}

type UserResponse struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Role      string           `json:"role"`
	ClientID  *uint            `json:"clientId,omitempty"`
	Client    *ClientResponse  `json:"client,omitempty"`
	Projects  []ProjectResponse `json:"projects,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

type ClientResponse struct {
	ID           uint              `json:"id"`
	Name         string            `json:"name"`
	Industry     *string           `json:"industry,omitempty"`
	ContactName  *string           `json:"contactName,omitempty"`
	ContactEmail *string           `json:"contactEmail,omitempty"`
	Users        []UserResponse    `json:"users,omitempty"`
	Projects     []ProjectResponse `json:"projects,omitempty"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

type RequirementResponse struct {
	ID         uint                  `json:"id"`
	ProjectID  uint                  `json:"projectId"`
	Text       string                `json:"text"`
	Category   *string               `json:"category,omitempty"`
	Status     string                `json:"status"`
	AuditTasks []AuditTaskResponse   `json:"auditTasks,omitempty"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

type AuditTaskResponse struct {
	ID            uint            `json:"id"`
	RequirementID uint            `json:"requirementId"`
	Text          string          `json:"text"`
	Status        string          `json:"status"`
	Notes         *string         `json:"notes,omitempty"`
	Issue         *IssueResponse  `json:"issue,omitempty"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type IssueResponse struct {
	ID          uint      `json:"id"`
	AuditTaskID uint      `json:"auditTaskId"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	Priority    *string   `json:"priority,omitempty"`
	Phase       *string   `json:"phase,omitempty"`
	EstimateHrs *int      `json:"estimateHrs,omitempty"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Request models for creating/updating
type CreateProjectRequest struct {
	Name       string `json:"name" binding:"required"`
	ClientName string `json:"clientName" binding:"required"`
	ClientID   *uint  `json:"clientId,omitempty"`
}

type UpdateProjectRequest struct {
	Name       *string `json:"name,omitempty"`
	ClientName *string `json:"clientName,omitempty"`
	Status     *string `json:"status,omitempty"`
	ClientID   *uint   `json:"clientId,omitempty"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=ADMIN CONSULTANT CLIENT"`
	ClientID *uint  `json:"clientId,omitempty"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Role     *string `json:"role,omitempty"`
	ClientID *uint   `json:"clientId,omitempty"`
}

type CreateClientRequest struct {
	Name         string  `json:"name" binding:"required"`
	Industry     *string `json:"industry,omitempty"`
	ContactName  *string `json:"contactName,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
}

type UpdateClientRequest struct {
	Name         *string `json:"name,omitempty"`
	Industry     *string `json:"industry,omitempty"`
	ContactName  *string `json:"contactName,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
}

type CreateRequirementRequest struct {
	ProjectID uint    `json:"projectId" binding:"required"`
	Text      string  `json:"text" binding:"required"`
	Category  *string `json:"category,omitempty"`
	Status    *string `json:"status,omitempty"`
}

type UpdateRequirementRequest struct {
	Text     *string `json:"text,omitempty"`
	Category *string `json:"category,omitempty"`
	Status   *string `json:"status,omitempty"`
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}