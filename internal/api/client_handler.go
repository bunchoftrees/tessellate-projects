package api

import (
	"net/http"
	"strconv"

	"tessellate-projects/internal/db"
	"tessellate-projects/internal/models"

	"github.com/gin-gonic/gin"
)

// ClientHandler
type ClientHandler struct {
	db *db.Database
}

func NewClientHandler(database *db.Database) *ClientHandler {
	return &ClientHandler{db: database}
}

func (h *ClientHandler) GetClients(c *gin.Context) {
	var clients []db.Client
	if err := h.db.Preload("Users").Preload("Projects").Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch clients", Code: 500})
		return
	}

	response := make([]models.ClientResponse, len(clients))
	for i, client := range clients {
		response[i] = h.convertToClientResponse(&client)
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	var req models.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	client := db.Client{
		Name:         req.Name,
		Industry:     req.Industry,
		ContactName:  req.ContactName,
		ContactEmail: req.ContactEmail,
	}

	if err := h.db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create client",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToClientResponse(&client)
	c.JSON(http.StatusCreated, response)
}

func (h *ClientHandler) GetClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid client ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var client db.Client
	if err := h.db.Preload("Users").Preload("Projects").First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Client not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	response := h.convertToClientResponse(&client)
	c.JSON(http.StatusOK, response)
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid client ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	var client db.Client
	if err := h.db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Client not found",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Update fields if provided
	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.Industry != nil {
		client.Industry = req.Industry
	}
	if req.ContactName != nil {
		client.ContactName = req.ContactName
	}
	if req.ContactEmail != nil {
		client.ContactEmail = req.ContactEmail
	}

	if err := h.db.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update client",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	response := h.convertToClientResponse(&client)
	c.JSON(http.StatusOK, response)
}

func (h *ClientHandler) DeleteClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid client ID",
			Code:  http.StatusBadRequest,
		})
		return
	}

	if err := h.db.Delete(&db.Client{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete client",
			Code:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}

// Helper function to convert db.Client to models.ClientResponse
func (h *ClientHandler) convertToClientResponse(client *db.Client) models.ClientResponse {
	response := models.ClientResponse{
		ID:           client.ID,
		Name:         client.Name,
		Industry:     client.Industry,
		ContactName:  client.ContactName,
		ContactEmail: client.ContactEmail,
		CreatedAt:    client.CreatedAt,
		UpdatedAt:    client.UpdatedAt,
	}

	// Convert users if loaded
	if len(client.Users) > 0 {
		response.Users = make([]models.UserResponse, len(client.Users))
		for i, user := range client.Users {
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

	// Convert projects if loaded
	if len(client.Projects) > 0 {
		response.Projects = make([]models.ProjectResponse, len(client.Projects))
		for i, project := range client.Projects {
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
