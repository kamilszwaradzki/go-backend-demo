# Go Backend Demo

![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)
![License](https://img.shields.io/badge/License-MIT-blue)

A simple and efficient REST API for task management built with Go and Gorilla Mux. Features in-memory storage, CORS support, and Docker deployment.

## Features

- Complete task CRUD operations
- Fast and lightweight
- CORS enabled for frontend integration
- In-memory data storage (no database required)
- Docker multi-stage build
- Clean JSON API responses

## Prerequisites

- Go 1.21+
- (Optional) Docker

## Installation

### Local Development

1. **Clone the repository**
```bash
git clone https://github.com/kamilszwaradzki/go-backend-demo.git
cd go-backend-demo
```

2. **Download dependencies**
```bash
go mod download
```

3. **Run the application**
```bash
go run main.go
```

Server will start on `http://localhost:8080`

4. **Test the API**
```bash
curl http://localhost:8080/api/tasks
```

### Docker

1. **Build the image**
```bash
docker build -t go-task-api .
```

2. **Run the container**
```bash
docker run -p 8080:8080 go-task-api
```

3. **Verify it's running**
```bash
curl http://localhost:8080/api/tasks
```

## Project Structure

```
go-backend-demo/
├── main.go           # Main application code
├── go.mod            # Go module definition
├── Dockerfile        # Multi-stage Docker build
└── README.md         # This file
```

## API Documentation

Base URL: `http://localhost:8080/api`

### Endpoints

#### Get All Tasks
```http
GET /api/tasks
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Complete project",
      "description": "Finish the task management API",
      "status": "in_progress",
      "created_at": "2024-01-04T10:00:00Z",
      "updated_at": "2024-01-04T10:00:00Z"
    }
  ]
}
```

#### Get Single Task
```http
GET /api/tasks/{id}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Complete project",
    "description": "Finish the task management API",
    "status": "in_progress",
    "created_at": "2024-01-04T10:00:00Z",
    "updated_at": "2024-01-04T10:00:00Z"
  }
}
```

#### Create Task
```http
POST /api/tasks
Content-Type: application/json

{
  "title": "New task",
  "description": "Task description",
  "status": "pending"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Task created successfully",
  "data": {
    "id": 2,
    "title": "New task",
    "description": "Task description",
    "status": "pending",
    "created_at": "2024-01-04T10:30:00Z",
    "updated_at": "2024-01-04T10:30:00Z"
  }
}
```

Status options: `pending`, `in_progress`, `completed`

#### Update Task
```http
PUT /api/tasks/{id}
Content-Type: application/json

{
  "title": "Updated title",
  "description": "Updated description",
  "status": "completed"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Task updated successfully",
  "data": {
    "id": 1,
    "title": "Updated title",
    "description": "Updated description",
    "status": "completed",
    "created_at": "2024-01-04T10:00:00Z",
    "updated_at": "2024-01-04T10:35:00Z"
  }
}
```

#### Delete Task
```http
DELETE /api/tasks/{id}
```

**Response:**
```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

### Error Responses

**404 Not Found:**
```json
{
  "success": false,
  "error": "Task not found"
}
```

**400 Bad Request:**
```json
{
  "success": false,
  "error": "Invalid task ID"
}
```

**422 Validation Error:**
```json
{
  "success": false,
  "error": "Title is required"
}
```

## Testing with cURL

### Create a task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Master Go programming","status":"pending"}'
```

### Get all tasks
```bash
curl http://localhost:8080/api/tasks
```

### Get single task
```bash
curl http://localhost:8080/api/tasks/1
```

### Update task
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go Advanced","status":"in_progress"}'
```

### Delete task
```bash
curl -X DELETE http://localhost:8080/api/tasks/1
```

## Code Overview

### Main Components

**Task Struct**
```go
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Response Struct**
```go
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

### CORS Configuration

The API includes CORS middleware that allows:
- All origins (`*`)
- Methods: GET, POST, PUT, DELETE, OPTIONS
- Headers: Content-Type, Authorization, Accept

For production, restrict `allowed_origins` to your frontend domain.

## Building for Production

### Compile binary
```bash
go build -o task-api main.go
```

### Run binary
```bash
./task-api
```

### Cross-compilation

**For Linux:**
```bash
GOOS=linux GOARCH=amd64 go build -o task-api-linux main.go
```

**For Windows:**
```bash
GOOS=windows GOARCH=amd64 go build -o task-api.exe main.go
```

**For macOS:**
```bash
GOOS=darwin GOARCH=amd64 go build -o task-api-mac main.go
```

## Docker Details

### Multi-stage Build

The Dockerfile uses a two-stage build:
1. **Builder stage**: Compiles the Go binary
2. **Runtime stage**: Creates minimal Alpine image with just the binary

This results in a small (~15MB) production image.

### Build Arguments
```bash
docker build \
  --build-arg GO_VERSION=1.21 \
  -t go-task-api:latest .
```

### Run with custom port
```bash
docker run -p 9000:8080 go-task-api
```

## Extending the Application

### Add Database Support

To add SQLite persistence:

1. **Install SQLite driver**
```bash
go get github.com/mattn/go-sqlite3
```

2. **Initialize database**
```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

db, err := sql.Open("sqlite3", "./tasks.db")
```

3. **Replace in-memory storage with database queries**

### Add Authentication

To add JWT authentication:

1. **Install JWT library**
```bash
go get github.com/golang-jwt/jwt/v5
```

2. **Create authentication middleware**
```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate JWT token
        // Call next.ServeHTTP(w, r) if valid
    })
}
```

### Add Logging

```go
import "log"

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

// Apply to router
r.Use(loggingMiddleware)
```

## Performance

- Lightweight: ~15MB Docker image
- Fast startup: <100ms
- Low memory footprint: ~10MB RAM
- Concurrent request handling with goroutines

## Limitations

- **In-memory storage**: Data is lost on restart
- **No authentication**: Anyone can access the API
- **No persistence**: Tasks are not saved to disk
- **Single instance**: No clustering support

For production use, consider adding:
- Database (PostgreSQL, MySQL, MongoDB)
- Authentication (JWT, OAuth)
- Rate limiting
- Request logging
- Health check endpoints

## Troubleshooting

### Port Already in Use
```bash
# Find process on port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Module Issues
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
```

### Docker Build Fails
```bash
# Clean Docker cache
docker system prune -a

# Rebuild without cache
docker build --no-cache -t go-task-api .
```

## Related Projects

- [React Frontend Demo](../react-frontend-demo) - React dashboard UI
- [Laravel Backend Demo](../php-backend-demo) - PHP alternative backend

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Author

Kamil Szwaradzki

## Resources

- [Gorilla Mux Documentation](https://github.com/gorilla/mux)
- [Go Documentation](https://golang.org/doc/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)