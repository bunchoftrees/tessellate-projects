package api

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"

	"tessellate-projects/internal/db"
	"tessellate-projects/internal/models"

	"github.com/gin-gonic/gin"
)

// RequirementHandler
type RequirementHandler struct {
	db *db.Database
}

func NewRequirementHandler(database *db.Database) *RequirementHandler {
	return &RequirementHandler{db: database}
}

func (h *RequirementHandler) GetRequirements(c *gin.Context) {
	var requirements []db.Requirement

	// Optional project filter
	projectID := c.Query("projectId")
	query := h.db.Preload("AuditTasks")

	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	if err := query.Find(&requirements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch requirements",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.RequirementResponse, len(requirements))
	for i, requirement := range requirements {
		response[i] = h.convertToRequirementResponse(&requirement)
	}

	c.JSON(http.StatusOK, response)
}

func (h *RequirementHandler) CreateRequirement(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var req models.CreateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Verify project exists
	var project db.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	status := db.RequirementStatusNotMet
	if req.Status != nil {
		status = db.RequirementStatus(*req.Status)
	}

	requirement := db.Requirement{
		ProjectID: uint(projectID),
		Text:      req.Text,
		Category:  req.Category,
		Status:    status,
	}

	if err := h.db.Create(&requirement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create requirement",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToRequirementResponse(&requirement)
	c.JSON(http.StatusCreated, response)
}

func (h *RequirementHandler) GetRequirement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid requirement ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var requirement db.Requirement
	if err := h.db.Preload("AuditTasks").First(&requirement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Requirement not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	response := h.convertToRequirementResponse(&requirement)
	c.JSON(http.StatusOK, response)
}

func (h *RequirementHandler) UpdateRequirement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid requirement ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var requirement db.Requirement
	if err := h.db.First(&requirement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Requirement not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Update fields if provided
	if req.Text != nil {
		requirement.Text = *req.Text
	}
	if req.Category != nil {
		requirement.Category = req.Category
	}
	if req.Status != nil {
		requirement.Status = db.RequirementStatus(*req.Status)
	}

	if err := h.db.Save(&requirement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update requirement",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToRequirementResponse(&requirement)
	c.JSON(http.StatusOK, response)
}

func (h *RequirementHandler) DeleteRequirement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid requirement ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err := h.db.Delete(&db.Requirement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete requirement",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Requirement deleted successfully"})
}

func (h *RequirementHandler) GetProjectRequirements(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var requirements []db.Requirement
	if err := h.db.Where("project_id = ?", projectID).Preload("AuditTasks").Find(&requirements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch project requirements",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.RequirementResponse, len(requirements))
	for i, requirement := range requirements {
		response[i] = h.convertToRequirementResponse(&requirement)
	}

	c.JSON(http.StatusOK, response)
}

func (h *RequirementHandler) UploadRequirementsCSV(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	// Verify project exists
	var project db.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "No file uploaded",
			Code:  http.StatusBadRequest,
		})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var createdRequirements []db.Requirement

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Error reading CSV",
				Message: err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}

		if len(record) == 0 {
			continue
		}

		text := record[0]
		var category *string
		if len(record) > 1 && record[1] != "" {
			category = &record[1]
		}

		requirement := db.Requirement{
			ProjectID: uint(projectID),
			Text:      text,
			Category:  category,
			Status:    db.RequirementStatusNotMet,
		}

		if err := h.db.Create(&requirement).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to create requirement from CSV",
				Code:  http.StatusInternalServerError,
			})
			return
		}

		createdRequirements = append(createdRequirements, requirement)
	}

	response := make([]models.RequirementResponse, len(createdRequirements))
	for i, requirement := range createdRequirements {
		response[i] = h.convertToRequirementResponse(&requirement)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Requirements uploaded successfully",
		"count":        len(createdRequirements),
		"requirements": response,
	})
}

// Helper function to convert db.Requirement to models.RequirementResponse
func (h *RequirementHandler) convertToRequirementResponse(requirement *db.Requirement) models.RequirementResponse {
	response := models.RequirementResponse{
		ID:        requirement.ID,
		ProjectID: requirement.ProjectID,
		Text:      requirement.Text,
		Category:  requirement.Category,
		Status:    string(requirement.Status),
		CreatedAt: requirement.CreatedAt,
		UpdatedAt: requirement.UpdatedAt,
	}

	// Convert audit tasks if loaded
	if len(requirement.AuditTasks) > 0 {
		response.AuditTasks = make([]models.AuditTaskResponse, len(requirement.AuditTasks))
		for i, task := range requirement.AuditTasks {
			response.AuditTasks[i] = models.AuditTaskResponse{
				ID:            task.ID,
				RequirementID: task.RequirementID,
				Text:          task.Text,
				Status:        task.Status,
				Notes:         task.Notes,
				CreatedAt:     task.CreatedAt,
				UpdatedAt:     task.UpdatedAt,
			}
		}
	}

	return response
}
