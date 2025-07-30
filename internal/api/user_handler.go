package api

import (
	"net/http"
	"strconv"

	"tessellate-projects/internal/auth"
	"tessellate-projects/internal/db"
	"tessellate-projects/internal/models"

	"github.com/gin-gonic/gin"
)

// UserHandler
type UserHandler struct {
	db *db.Database
}

func NewUserHandler(database *db.Database) *UserHandler {
	return &UserHandler{db: database}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var users []db.User
	if err := h.db.Preload("Client").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch users", Code: 500})
		return
	}

	response := make([]models.UserResponse, len(users))
	for i, user := range users {
		response[i] = h.convertToUserResponse(&user)
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user := db.User{
		Name:     req.Name,
		Email:    req.Email,
		Role:     db.Role(req.Role),
		ClientID: req.ClientID,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create user",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToUserResponse(&user)
	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var user db.User
	if err := h.db.Preload("Client").Preload("Projects").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	response := h.convertToUserResponse(&user)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var user db.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = db.Role(*req.Role)
	}
	if req.ClientID != nil {
		user.ClientID = req.ClientID
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update user",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToUserResponse(&user)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err := h.db.Delete(&db.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete user",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetProjectUsers(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var project db.Project
	if err := h.db.Preload("Users").Preload("Users.Client").First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	response := make([]models.UserResponse, len(project.Users))
	for i, user := range project.Users {
		response[i] = h.convertToUserResponse(user)
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) AssignUserToProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var project db.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	var user db.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Add user to project
	if err := h.db.Model(&project).Association("Users").Append(&user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to assign user to project",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User assigned to project successfully"})
}

func (h *UserHandler) RemoveUserFromProject(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid project ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid user ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var project db.Project
	if err := h.db.First(&project, projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Project not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	var user db.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "User not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Remove user from project
	if err := h.db.Model(&project).Association("Users").Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to remove user from project",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from project successfully"})
}

func (h *UserHandler) GetClientUsers(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid client ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var users []db.User
	if err := h.db.Where("client_id = ?", clientID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch client users",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := make([]models.UserResponse, len(users))
	for i, user := range users {
		response[i] = h.convertToUserResponse(&user)
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var user db.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid credentials",
			Code:  http.StatusUnauthorized,
		})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid credentials",
			Code:  http.StatusUnauthorized,
		})
		return
	}

	// For now, return a simple token (in production, use JWT)
	response := gin.H{
		"token": "mock-token-" + strconv.Itoa(int(user.ID)),
		"user":  h.convertToUserResponse(&user),
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to convert db.User to models.UserResponse
func (h *UserHandler) convertToUserResponse(user *db.User) models.UserResponse {
	response := models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		ClientID:  user.ClientID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Convert client if loaded
	if user.Client != nil {
		response.Client = &models.ClientResponse{
			ID:           user.Client.ID,
			Name:         user.Client.Name,
			Industry:     user.Client.Industry,
			ContactName:  user.Client.ContactName,
			ContactEmail: user.Client.ContactEmail,
			CreatedAt:    user.Client.CreatedAt,
			UpdatedAt:    user.Client.UpdatedAt,
		}
	}

	// Convert projects if loaded
	if len(user.Projects) > 0 {
		response.Projects = make([]models.ProjectResponse, len(user.Projects))
		for i, project := range user.Projects {
			response.Projects[i] = models.ProjectResponse{
				ID:         project.ID,
				Name:       project.Name,
				ClientName: project.ClientName,
				Status:     project.Status,
				ClientID:   project.ClientID,
				CreatedAt:  project.CreatedAt,
				UpdatedAt:  project.UpdatedAt,
			}
		}
	}

	return response
}
