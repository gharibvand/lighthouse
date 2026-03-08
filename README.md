# Lighthouse - Video Streaming Platform Backend

A Netflix-like video streaming platform backend built with Go.

## Features

- User Authentication & Authorization (JWT)
- Multiple User Profiles per Account
- Content Management (Movies & TV Shows)
- Video Streaming (HLS)
- Watch History & Continue Watching
- Recommendations Engine
- Search & Discovery
- Watchlist & Ratings
- User Interactions
- Swagger API Documentation

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 18+ (for native UUID v7 support)
- Redis (optional, for caching)
- Docker & Docker Compose (for easy setup)

## Setup

### 1. Start Docker Containers

```powershell
docker-compose up -d
```

### 2. Run Database Migrations

```powershell
# Install migrate tool (if not installed)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database "postgres://lighthouse_user:lighthouse_password@localhost:5432/lighthouse?sslmode=disable" up
```

**Or if migrate is in your PATH:**
```powershell
migrate -path ./migrations -database "postgres://lighthouse_user:lighthouse_password@localhost:5432/lighthouse?sslmode=disable" up
```

**Rollback migrations:**
```powershell
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database "postgres://lighthouse_user:lighthouse_password@localhost:5432/lighthouse?sslmode=disable" down
```

### 3. Run Database Seeders (Optional)

After running migrations, you can seed the database with initial data:

```powershell
go run cmd/seed/main.go
```

This will execute all SQL files in the `seeders/` directory. Currently includes:
- Default subscription plans (Basic, Standard, Premium)

**Note:** Seeders are idempotent (safe to run multiple times) using `ON CONFLICT DO NOTHING`.

### 4. Environment Variables

Create a `.env` file in the project root:

```env
PORT=8080
DATABASE_URL=postgres://lighthouse_user:lighthouse_password@localhost:5432/lighthouse?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production-make-it-long-and-random
STORAGE_PATH=./storage
REDIS_URL=redis://localhost:6379
```

**Note:** The `.env` file is optional. If not found, the server will use default values that match the docker-compose configuration.

**Generate a random JWT secret:**
```powershell
-join ((48..57) + (65..90) + (97..122) | Get-Random -Count 64 | ForEach-Object {[char]$_})
```

### 5. Install Dependencies

```powershell
go mod download
```

### 6. Generate Swagger Documentation (Optional)

```powershell
# Install swag (if not installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
& "$env:USERPROFILE\go\bin\swag.exe" init -g cmd/server/main.go -o ./docs
```

**Or if swag is in your PATH:**
```powershell
swag init -g cmd/server/main.go -o ./docs
```

**Note:** Run this command whenever you add new endpoints or change Swagger annotations.

### 7. Run the Server

```powershell
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

**Stop any process using port 8080 (if needed):**
```powershell
Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue | 
    Select-Object -ExpandProperty OwningProcess | 
    ForEach-Object { Stop-Process -Id $_ -Force }
```

## Access Points

- **API Base URL:** `http://localhost:8080/api/v1`
- **Swagger UI:** `http://localhost:8080/swagger/index.html`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh token
- `GET /api/v1/auth/me` - Get current user (protected)

### Content
- `GET /api/v1/content/movies` - Get movies list
- `GET /api/v1/content/tv-shows` - Get TV shows list
- `GET /api/v1/content/:id` - Get content by ID
- `GET /api/v1/content/search?q=query` - Search content
- `GET /api/v1/content/trending` - Get trending content
- `GET /api/v1/content/recommendations` - Get recommendations (protected)

### Streaming
- `GET /api/v1/stream/:contentId/playlist.m3u8` - Get HLS playlist
- `GET /api/v1/stream/:contentId/:quality/:segment` - Get video segment
- `POST /api/v1/stream/:contentId/watch-progress` - Update watch progress (protected)

### User
- `GET /api/v1/user/profiles` - Get user profiles (protected)
- `POST /api/v1/user/profiles` - Create profile (protected)
- `PUT /api/v1/user/profiles/:id` - Update profile (protected)
- `DELETE /api/v1/user/profiles/:id` - Delete profile (protected)
- `GET /api/v1/user/watchlist` - Get watchlist (protected)
- `POST /api/v1/user/watchlist` - Add to watchlist (protected)
- `DELETE /api/v1/user/watchlist/:contentId` - Remove from watchlist (protected)
- `GET /api/v1/user/history` - Get watch history (protected)
- `POST /api/v1/user/ratings` - Add/Update rating (protected)

## Project Structure

```
lighthouse/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── domain/          # Domain models
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   ├── handler/         # HTTP handlers
│   └── middleware/      # Middleware
├── pkg/                 # Shared packages
├── migrations/          # Database migrations
├── docs/                # Swagger documentation (generated)
└── storage/             # Video files storage
```

## Development

### Running Tests
```powershell
go test ./...
```

### Building the Server
```powershell
go build -o lighthouse.exe ./cmd/server
```

### Database Migrations

**Create new migration:**
```powershell
& "$env:USERPROFILE\go\bin\migrate.exe" create -ext sql -dir ./migrations -seq migration_name
```

**Check migration version:**
```powershell
& "$env:USERPROFILE\go\bin\migrate.exe" -path ./migrations -database "postgres://lighthouse_user:lighthouse_password@localhost:5432/lighthouse?sslmode=disable" version
```

### Adding Swagger Documentation

Add annotations to your handlers:

```go
// GetExample godoc
// @Summary      Example endpoint
// @Description  This is an example endpoint
// @Tags         example
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Content ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /example/{id} [get]
func (h *Handler) GetExample(c *gin.Context) {
    // ...
}
```

**For protected endpoints, add:**
```go
// @Security     BearerAuth
```

**Then regenerate documentation:**
```powershell
& "$env:USERPROFILE\go\bin\swag.exe" init -g cmd/server/main.go -o ./docs
```

## Environment Variables Reference

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://lighthouse_user:lighthouse_password@localhost:5432/lighthouse?sslmode=disable` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key-change-in-production` |
| `STORAGE_PATH` | Path for video file storage | `./storage` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379` |

## License

MIT
