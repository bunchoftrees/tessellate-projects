package api

import (
	"net/http"
	"strconv"

	"tessellate-projects/internal/db"
	"tessellate-projects/internal/models"

	"github.com/gin-gonic/gin"
)

// AuditTaskHandler
type AuditTaskHandler struct {
	db *db.Database
}

func NewAuditTaskHandler(database *db.Database) *AuditTaskHandler {
	return &AuditTaskHandler{db: database}
}

func (h *AuditTaskHandler) GetAuditTasks(c *gin.Context) {
	var tasks []db.AuditTask

	// Optional requirement filter
	requirementID := c.Query("requirementId")
	query := h.db.Preload("Issue")

	if requirementID != "" {
		query = query.Where("requirement_id = ?", requirementID)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch audit tasks",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.AuditTaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = h.convertToAuditTaskResponse(&task)
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuditTaskHandler) CreateAuditTask(c *gin.Context) {
	requirementID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid requirement ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	type CreateAuditTaskRequest struct {
		Text   string  `json:"text" binding:"required"`
		Status *string `json:"status,omitempty"`
		Notes  *string `json:"notes,omitempty"`
	}

	var req CreateAuditTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Verify requirement exists
	var requirement db.Requirement
	if err := h.db.First(&requirement, requirementID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Requirement not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	status := "PENDING"
	if req.Status != nil {
		status = *req.Status
	}

	task := db.AuditTask{
		RequirementID: uint(requirementID),
		Text:          req.Text,
		Status:        status,
		Notes:         req.Notes,
	}

	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create audit task",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToAuditTaskResponse(&task)
	c.JSON(http.StatusCreated, response)
}

func (h *AuditTaskHandler) GetAuditTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid audit task ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var task db.AuditTask
	if err := h.db.Preload("Issue").First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Audit task not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	response := h.convertToAuditTaskResponse(&task)
	c.JSON(http.StatusOK, response)
}

func (h *AuditTaskHandler) UpdateAuditTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid audit task ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	type UpdateAuditTaskRequest struct {
		Text   *string `json:"text,omitempty"`
		Status *string `json:"status,omitempty"`
		Notes  *string `json:"notes,omitempty"`
	}

	var req UpdateAuditTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var task db.AuditTask
	if err := h.db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Audit task not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Update fields if provided
	if req.Text != nil {
		task.Text = *req.Text
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Notes != nil {
		task.Notes = req.Notes
	}

	if err := h.db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update audit task",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToAuditTaskResponse(&task)
	c.JSON(http.StatusOK, response)
}

func (h *AuditTaskHandler) DeleteAuditTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid audit task ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err := h.db.Delete(&db.AuditTask{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete audit task",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Audit task deleted successfully"})
}

func (h *AuditTaskHandler) GetRequirementAuditTasks(c *gin.Context) {
	requirementID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid requirement ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var tasks []db.AuditTask
	if err := h.db.Where("requirement_id = ?", requirementID).Preload("Issue").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch requirement audit tasks",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.AuditTaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = h.convertToAuditTaskResponse(&task)
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to convert db.AuditTask to models.AuditTaskResponse
func (h *AuditTaskHandler) convertToAuditTaskResponse(task *db.AuditTask) models.AuditTaskResponse {
	response := models.AuditTaskResponse{
		ID:            task.ID,
		RequirementID: task.RequirementID,
		Text:          task.Text,
		Status:        task.Status,
		Notes:         task.Notes,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
	}

	// Convert issue if loaded
	if task.Issue != nil {
		response.Issue = &models.IssueResponse{
			ID:          task.Issue.ID,
			AuditTaskID: task.Issue.AuditTaskID,
			Title:       task.Issue.Title,
			Description: task.Issue.Description,
			Priority:    task.Issue.Priority,
			Phase:       task.Issue.Phase,
			EstimateHrs: task.Issue.EstimateHrs,
			Status:      task.Issue.Status,
			Type:        task.Issue.Type,
			CreatedAt:   task.Issue.CreatedAt,
			UpdatedAt:   task.Issue.UpdatedAt,
		}
	}

	return response
}
