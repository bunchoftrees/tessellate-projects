package api

import (
	"net/http"
	"strconv"

	"tessellate-projects/internal/db"
	"tessellate-projects/internal/models"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	db *db.Database
}

func NewProjectHandler(database *db.Database) *ProjectHandler {
	return &ProjectHandler{db: database}
}

// GetProjects handles GET /api/v1/projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	var projects []db.Project

	// Optional status filter
	status := c.Query("status")
	query := h.db.DB

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch projects",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	// Convert to response models
	response := make([]models.ProjectResponse, len(projects))
	for i, project := range projects {
		response[i] = models.ProjectResponse{
			ID:         project.ID,
			Name:       project.Name,
			ClientName: project.ClientName,
			Status:     project.Status,
			ClientID:   project.ClientID,
			CreatedAt:  project.CreatedAt,
			UpdatedAt:  project.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// GetProject handles GET /api/v1/projects/:id
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var project db.Project
	if err := h.db.Preload("Client").Preload("Users").Preload("Requirements").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Convert to response model
	response := h.convertToProjectResponse(&project)
	c.JSON(http.StatusOK, response)
}

// CreateProject handles POST /api/v1/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	project := db.Project{
		Name:       req.Name,
		ClientName: req.ClientName,
		Status:     "NEW",
		ClientID:   req.ClientID,
	}

	if err := h.db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create project",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToProjectResponse(&project)
	c.JSON(http.StatusCreated, response)
}

// UpdateProject handles PUT /api/v1/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var project db.Project
	if err := h.db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Update fields if provided
	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.ClientName != nil {
		project.ClientName = *req.ClientName
	}
	if req.Status != nil {
		project.Status = *req.Status
	}
	if req.ClientID != nil {
		project.ClientID = req.ClientID
	}

	if err := h.db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update project",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToProjectResponse(&project)
	c.JSON(http.StatusOK, response)
}

// DeleteProject handles DELETE /api/v1/projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err := h.db.Delete(&db.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete project",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// ArchiveProject handles POST /api/v1/projects/:id/archive
func (h *ProjectHandler) ArchiveProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var project db.Project
	if err := h.db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	project.Status = "ARCHIVED"
	if err := h.db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to archive project",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToProjectResponse(&project)
	c.JSON(http.StatusOK, response)
}

// GetUserProjects handles GET /api/v1/users/:id/projects
func (h *ProjectHandler) GetUserProjects(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var user db.User
	if err := h.db.Preload("Projects").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	response := make([]models.ProjectResponse, len(user.Projects))
	for i, project := range user.Projects {
		response[i] = h.convertToProjectResponse(project)
	}

	c.JSON(http.StatusOK, response)
}

// GetClientProjects handles GET /api/v1/clients/:id/projects
func (h *ProjectHandler) GetClientProjects(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid client ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var projects []db.Project
	if err := h.db.Where("client_id = ?", clientID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch client projects",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.ProjectResponse, len(projects))
	for i, project := range projects {
		response[i] = h.convertToProjectResponse(&project)
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to convert db.Project to models.ProjectResponse
func (h *ProjectHandler) convertToProjectResponse(project *db.Project) models.ProjectResponse {
	response := models.ProjectResponse{
		ID:         project.ID,
		Name:       project.Name,
		ClientName: project.ClientName,
		Status:     project.Status,
		ClientID:   project.ClientID,
		CreatedAt:  project.CreatedAt,
		UpdatedAt:  project.UpdatedAt,
	}

	// Convert related entities if loaded
	if project.Client != nil {
		response.Client = &models.ClientResponse{
			ID:           project.Client.ID,
			Name:         project.Client.Name,
			Industry:     project.Client.Industry,
			ContactName:  project.Client.ContactName,
			ContactEmail: project.Client.ContactEmail,
			CreatedAt:    project.Client.CreatedAt,
			UpdatedAt:    project.Client.UpdatedAt,
		}
	}

	if len(project.Users) > 0 {
		response.Users = make([]models.UserResponse, len(project.Users))
		for i, user := range project.Users {
			response.Users[i] = models.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				Role:      string(user.Role),
				ClientID:  user.ClientID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
		}
	}

	return response
}
