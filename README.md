# Tessellate Projects API

A comprehensive audit management system backend built with Go, Gin, and GORM. This API provides a complete solution for managing audit projects, requirements, tasks, and issues across multiple clients and users.

## Overview

Tessellate Projects API is designed to streamline the audit process by providing structured management of:
- **Projects**: Client audit engagements
- **Requirements**: Specific audit criteria that need to be evaluated
- **Audit Tasks**: Individual tasks for evaluating requirements
- **Issues**: Problems or defects discovered during audits
- **Users & Clients**: Role-based access and client management

## Architecture

### Tech Stack
- **Go 1.23** - Backend runtime
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **SQLite** - Database (easily swappable for PostgreSQL/MySQL)
- **bcrypt** - Password hashing

### Project Structure
```
├── cmd/server/          # Application entry point
├── internal/
│   ├── api/            # HTTP handlers and routing
│   ├── auth/           # Authentication utilities
│   ├── db/             # Database models and connection
│   └── models/         # API request/response models
├── go.mod              # Go module dependencies
└── README.md
```

## Data Model

The system follows a hierarchical structure:

```
Client
├── Users (multiple)
└── Projects (multiple)
    └── Requirements (multiple)
        └── Audit Tasks (multiple)
            └── Issues (0 or 1 per task)
```

### Core Entities

- **Client**: Organizations being audited
- **User**: System users with roles (ADMIN, CONSULTANT, CLIENT)
- **Project**: Individual audit engagements
- **Requirement**: Specific criteria to be evaluated
- **Audit Task**: Tasks for checking requirements
- **Issue**: Problems discovered during audits

## API Endpoints

### Base URL: `/api/v1`

#### Projects
- `GET /projects` - List all projects (with status filter)
- `POST /projects` - Create new project
- `GET /projects/:id` - Get project details
- `PUT /projects/:id` - Update project
- `DELETE /projects/:id` - Delete project
- `POST /projects/:id/archive` - Archive project

#### Users
- `GET /users` - List all users
- `POST /users` - Create new user
- `GET /users/:id` - Get user details
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

#### Clients
- `GET /clients` - List all clients
- `POST /clients` - Create new client
- `GET /clients/:id` - Get client details
- `PUT /clients/:id` - Update client
- `DELETE /clients/:id` - Delete client

#### Requirements
- `GET /requirements` - List requirements (with project filter)
- `GET /requirements/:id` - Get requirement details
- `PUT /requirements/:id` - Update requirement
- `DELETE /requirements/:id` - Delete requirement
- `POST /projects/:id/requirements` - Create requirement for project

#### Audit Tasks
- `GET /audit-tasks` - List audit tasks (with requirement filter)
- `GET /audit-tasks/:id` - Get audit task details
- `PUT /audit-tasks/:id` - Update audit task
- `DELETE /audit-tasks/:id` - Delete audit task
- `POST /requirements/:id/audit-tasks` - Create audit task for requirement

#### Issues
- `GET /issues` - List issues (with audit task filter)
- `GET /issues/:id` - Get issue details
- `PUT /issues/:id` - Update issue
- `DELETE /issues/:id` - Delete issue
- `POST /audit-tasks/:id/issues` - Create issue for audit task

#### File Uploads
- `POST /uploads/requirements-csv/:projectId` - Bulk upload requirements via CSV

#### Authentication
- `POST /auth/login` - User login

## Getting Started

### Prerequisites
- Go 1.23 or later
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd tessellate-projects
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on port 8080 (configurable via `PORT` environment variable).

### Environment Variables

Create a `.env` file in the root directory:

```env
PORT=8080
GIN_MODE=debug
# Add database connection string if using external DB
```

### Database

The application uses SQLite by default with the database file `tessellate-projects.db`. On first run, it will:
- Create the database schema automatically
- Seed sample data for testing

## API Usage Examples

### Create a Project
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Security Audit 2024",
    "clientName": "Acme Corp",
    "clientId": 1
  }'
```

### Upload Requirements via CSV
```bash
curl -X POST http://localhost:8080/api/v1/uploads/requirements-csv/1 \
  -F "file=@requirements.csv"
```

CSV format:
```csv
Requirement Text,Category
"Must support two-factor authentication",Security
"System must log all user actions",Compliance
```

### Create an Audit Task
```bash
curl -X POST http://localhost:8080/api/v1/requirements/1/audit-tasks \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Verify 2FA implementation",
    "status": "PENDING",
    "notes": "Check all login flows"
  }'
```

## Features

### Role-Based Access
- **ADMIN**: Full system access
- **CONSULTANT**: Manage projects and audits
- **CLIENT**: View own projects and data

### Bulk Operations
- CSV upload for requirements
- Batch operations support

### Flexible Status Tracking
- Project status: NEW, IN_PROGRESS, COMPLETED, ARCHIVED
- Requirement status: MET, NOT_MET
- Audit task status: PENDING, IN_PROGRESS, COMPLETED
- Issue status: OPEN, IN_PROGRESS, RESOLVED, CLOSED

### CORS Support
Development-friendly CORS configuration for frontend integration.

## Development

### Database Schema Updates
The application uses GORM's auto-migration feature. Schema changes are automatically applied on startup.

### Adding New Endpoints
1. Create handler in `internal/api/`
2. Add route in `internal/api/routes.go`
3. Define request/response models in `internal/models/`

### Testing
```bash
# Run tests
go test ./...

# Test API endpoints
curl http://localhost:8080/api/v1/
```

## Production Deployment

### Docker Support
A Dockerfile will be implemented (TODO) provided for containerized deployment:

```bash
docker build -t tessellate-projects-api .
docker run -p 8080:8080 tessellate-projects-api
```

### Environment Configuration
Set `GIN_MODE=release` for production deployment.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See [LICENSE.txt](LICENSE.txt) for details.

## API Documentation

Visit `http://localhost:8080/api/v1` for a basic endpoint overview. The API returns JSON responses with consistent error handling and proper HTTP status codes.

### Response Format
```json
{
  "id": 1,
  "name": "Example Project",
  "status": "NEW",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Error Format
```json
{
  "error": "Resource not found",
  "message": "Project with ID 123 does not exist",
  "code": 404
}
```