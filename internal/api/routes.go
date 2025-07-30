package api

import (
	"tessellate-projects/internal/db"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, database *db.Database) {
	// Create handlers
	projectHandler := NewProjectHandler(database)
	userHandler := NewUserHandler(database)
	clientHandler := NewClientHandler(database)
	requirementHandler := NewRequirementHandler(database)
	auditTaskHandler := NewAuditTaskHandler(database)
	issueHandler := NewIssueHandler(database)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Projects
		projects := v1.Group("/projects")
		{
			projects.GET("", projectHandler.GetProjects)
			projects.POST("", projectHandler.CreateProject)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
			projects.POST("/:id/archive", projectHandler.ArchiveProject)

			// Project relationships
			projects.GET("/:id/requirements", requirementHandler.GetProjectRequirements)
			projects.POST("/:id/requirements", requirementHandler.CreateRequirement)
			projects.GET("/:id/users", userHandler.GetProjectUsers)
			projects.POST("/:id/users/:userId", userHandler.AssignUserToProject)
			projects.DELETE("/:id/users/:userId", userHandler.RemoveUserFromProject)
			projects.GET("/:id/issues", issueHandler.GetProjectIssues)
		}

		// Users
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/:id/projects", projectHandler.GetUserProjects)
		}

		// Clients
		clients := v1.Group("/clients")
		{
			clients.GET("", clientHandler.GetClients)
			clients.POST("", clientHandler.CreateClient)
			clients.GET("/:id", clientHandler.GetClient)
			clients.PUT("/:id", clientHandler.UpdateClient)
			clients.DELETE("/:id", clientHandler.DeleteClient)
			clients.GET("/:id/users", userHandler.GetClientUsers)
			clients.GET("/:id/projects", projectHandler.GetClientProjects)
		}

		// Requirements
		requirements := v1.Group("/requirements")
		{
			requirements.GET("", requirementHandler.GetRequirements)
			requirements.GET("/:id", requirementHandler.GetRequirement)
			requirements.PUT("/:id", requirementHandler.UpdateRequirement)
			requirements.DELETE("/:id", requirementHandler.DeleteRequirement)
			requirements.GET("/:id/audit-tasks", auditTaskHandler.GetRequirementAuditTasks)
			requirements.POST("/:id/audit-tasks", auditTaskHandler.CreateAuditTask)
		}

		// Audit Tasks
		auditTasks := v1.Group("/audit-tasks")
		{
			auditTasks.GET("", auditTaskHandler.GetAuditTasks)
			auditTasks.GET("/:id", auditTaskHandler.GetAuditTask)
			auditTasks.PUT("/:id", auditTaskHandler.UpdateAuditTask)
			auditTasks.DELETE("/:id", auditTaskHandler.DeleteAuditTask)
			auditTasks.GET("/:id/issues", issueHandler.GetAuditTaskIssues)
			auditTasks.POST("/:id/issues", issueHandler.CreateIssue)
		}

		// Issues
		issues := v1.Group("/issues")
		{
			issues.GET("", issueHandler.GetIssues)
			issues.GET("/:id", issueHandler.GetIssue)
			issues.PUT("/:id", issueHandler.UpdateIssue)
			issues.DELETE("/:id", issueHandler.DeleteIssue)
		}

		// File uploads
		uploads := v1.Group("/uploads")
		{
			uploads.POST("/requirements-csv/:projectId", requirementHandler.UploadRequirementsCSV)
		}

		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}
	}

	// API documentation endpoint
	v1.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Tessellate Projects API v1",
			"endpoints": gin.H{
				"projects":     "/api/v1/projects",
				"users":        "/api/v1/users",
				"clients":      "/api/v1/clients",
				"requirements": "/api/v1/requirements",
				"audit-tasks":  "/api/v1/audit-tasks",
				"issues":       "/api/v1/issues",
				"uploads":      "/api/v1/uploads",
				"auth":         "/api/v1/auth",
			},
		})
	})
}
