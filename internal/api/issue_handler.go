package api

import (
	"net/http"
	"strconv"

	"tessellate-projects/internal/db"
	"tessellate-projects/internal/models"

	"github.com/gin-gonic/gin"
)

// IssueHandler
type IssueHandler struct {
	db *db.Database
}

func NewIssueHandler(database *db.Database) *IssueHandler {
	return &IssueHandler{db: database}
}

func (h *IssueHandler) GetIssues(c *gin.Context) {
	var issues []db.Issue

	// Optional audit task filter
	auditTaskID := c.Query("auditTaskId")
	query := h.db.DB

	if auditTaskID != "" {
		query = query.Where("audit_task_id = ?", auditTaskID)
	}

	if err := query.Find(&issues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch issues",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.IssueResponse, len(issues))
	for i, issue := range issues {
		response[i] = h.convertToIssueResponse(&issue)
	}

	c.JSON(http.StatusOK, response)
}

func (h *IssueHandler) CreateIssue(c *gin.Context) {
	auditTaskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid audit task ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	type CreateIssueRequest struct {
		Title       string  `json:"title" binding:"required"`
		Description *string `json:"description,omitempty"`
		Priority    *string `json:"priority,omitempty"`
		Phase       *string `json:"phase,omitempty"`
		EstimateHrs *int    `json:"estimateHrs,omitempty"`
		Status      *string `json:"status,omitempty"`
		Type        *string `json:"type,omitempty"`
	}

	var req CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Verify audit task exists
	var auditTask db.AuditTask
	if err := h.db.First(&auditTask, auditTaskID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Audit task not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	status := "OPEN"
	if req.Status != nil {
		status = *req.Status
	}

	issueType := "DEFECT"
	if req.Type != nil {
		issueType = *req.Type
	}

	issue := db.Issue{
		AuditTaskID: uint(auditTaskID),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Phase:       req.Phase,
		EstimateHrs: req.EstimateHrs,
		Status:      status,
		Type:        issueType,
	}

	if err := h.db.Create(&issue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create issue",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToIssueResponse(&issue)
	c.JSON(http.StatusCreated, response)
}

func (h *IssueHandler) GetIssue(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid issue ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var issue db.Issue
	if err := h.db.First(&issue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Issue not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	response := h.convertToIssueResponse(&issue)
	c.JSON(http.StatusOK, response)
}

func (h *IssueHandler) UpdateIssue(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid issue ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	type UpdateIssueRequest struct {
		Title       *string `json:"title,omitempty"`
		Description *string `json:"description,omitempty"`
		Priority    *string `json:"priority,omitempty"`
		Phase       *string `json:"phase,omitempty"`
		EstimateHrs *int    `json:"estimateHrs,omitempty"`
		Status      *string `json:"status,omitempty"`
		Type        *string `json:"type,omitempty"`
	}

	var req UpdateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var issue db.Issue
	if err := h.db.First(&issue, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Issue not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Update fields if provided
	if req.Title != nil {
		issue.Title = *req.Title
	}
	if req.Description != nil {
		issue.Description = req.Description
	}
	if req.Priority != nil {
		issue.Priority = req.Priority
	}
	if req.Phase != nil {
		issue.Phase = req.Phase
	}
	if req.EstimateHrs != nil {
		issue.EstimateHrs = req.EstimateHrs
	}
	if req.Status != nil {
		issue.Status = *req.Status
	}
	if req.Type != nil {
		issue.Type = *req.Type
	}

	if err := h.db.Save(&issue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update issue",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToIssueResponse(&issue)
	c.JSON(http.StatusOK, response)
}

func (h *IssueHandler) DeleteIssue(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid issue ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err := h.db.Delete(&db.Issue{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete issue",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue deleted successfully"})
}

func (h *IssueHandler) GetProjectIssues(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var issues []db.Issue
	if err := h.db.Joins("JOIN audit_tasks ON issues.audit_task_id = audit_tasks.id").
		Joins("JOIN requirements ON audit_tasks.requirement_id = requirements.id").
		Where("requirements.project_id = ?", projectID).
		Find(&issues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch project issues",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.IssueResponse, len(issues))
	for i, issue := range issues {
		response[i] = h.convertToIssueResponse(&issue)
	}

	c.JSON(http.StatusOK, response)
}

func (h *IssueHandler) GetAuditTaskIssues(c *gin.Context) {
	auditTaskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid audit task ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var issues []db.Issue
	if err := h.db.Where("audit_task_id = ?", auditTaskID).Find(&issues).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch audit task issues",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.IssueResponse, len(issues))
	for i, issue := range issues {
		response[i] = h.convertToIssueResponse(&issue)
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to convert db.Issue to models.IssueResponse
func (h *IssueHandler) convertToIssueResponse(issue *db.Issue) models.IssueResponse {
	return models.IssueResponse{
		ID:          issue.ID,
		AuditTaskID: issue.AuditTaskID,
		Title:       issue.Title,
		Description: issue.Description,
		Priority:    issue.Priority,
		Phase:       issue.Phase,
		EstimateHrs: issue.EstimateHrs,
		Status:      issue.Status,
		Type:        issue.Type,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   issue.UpdatedAt,
	}
}
