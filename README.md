# Social Media Application

This is a microservices-based social media application built with Go. It consists of 5 separate services:

1. **Gateway API**: RESTful API that provides all endpoints to the frontend and communicates with other services via gRPC.
2. **Users API**: Handles user authentication and management.
3. **Post API**: Manages posts, comments, and likes.
4. **Friends API**: Manages friend relationships, requests, and blocking.
5. **Groups API**: Manages groups, memberships, and group posts.

## Features

- User registration and authentication using Microsoft or Google OAuth
- User profile management
- Friend management (add, remove, block)
- Group management (create, join, leave)
- Post creation with privacy settings (public, private, group)
- Comments and likes on posts
- Swagger documentation for the Gateway API

## Architecture

The application follows a microservices architecture with the following components:

- **Gateway API**: RESTful API built with Gin framework
- **Service APIs**: gRPC services built with Go
- **Database**: MySQL for data storage
- **Authentication**: JWT-based authentication

Each microservice follows the MVC (Model-View-Controller) architecture:

```
service-api/
├── cmd/                  # Application entry points
├── internal/             # Private application code
│   ├── models/           # Data models
│   ├── controllers/      # Request handlers
│   ├── services/         # Business logic
│   ├── repository/       # Data access layer
│   ├── middleware/       # Middleware components
│   ├── config/           # Configuration
│   └── utils/            # Utility functions
├── database/             # Database migrations and seeds
│   ├── migrations/       # Database schema migrations
│   └── seeds/            # Seed data
├── docs/                 # Documentation
└── scripts/              # Utility scripts
```

## Prerequisites

- Go 1.16 or higher
- MySQL 8.0 or higher
- Protocol Buffers compiler (protoc)
- Docker (optional)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/social-media.git
cd social-media
```

2. Install dependencies for each service:

```bash
# Gateway API
cd gateway-api
go mod tidy

# Users API
cd ../users-api
go mod tidy

# Post API
cd ../post-api
go mod tidy

# Friends API
cd ../friends-api
go mod tidy

# Groups API
cd ../groups-api
go mod tidy
```

3. Set up the database:

```bash
# Create databases
mysql -u root -p -e "CREATE DATABASE users_db;"
mysql -u root -p -e "CREATE DATABASE posts_db;"
mysql -u root -p -e "CREATE DATABASE friends_db;"
mysql -u root -p -e "CREATE DATABASE groups_db;"

# Run migrations for each service
cd users-api
go run cmd/migrate/main.go

cd ../post-api
go run cmd/migrate/main.go

cd ../friends-api
go run cmd/migrate/main.go

cd ../groups-api
go run cmd/migrate/main.go
```

## Running the Services

### Users API

```bash
cd users-api
go run cmd/main.go
```

The Users API will be available at `localhost:50051`.

### Post API

```bash
cd post-api
go run cmd/main.go
```

The Post API will be available at `localhost:50052`.

### Friends API

```bash
cd friends-api
go run cmd/main.go
```

The Friends API will be available at `localhost:50053`.

### Groups API

```bash
cd groups-api
go run cmd/main.go
```

The Groups API will be available at `localhost:50054`.

### Gateway API

```bash
cd gateway-api
go run cmd/main.go
```

The Gateway API will be available at `localhost:8000`.

Swagger documentation will be available at `localhost:8000/swagger/index.html`.

## Configuration

Each service can be configured using environment variables or a configuration file. The default configuration file is located at `config/config.yaml` in each service directory.

Example configuration for the Gateway API:

```yaml
environment: development
port: 8000
users_service_url: localhost:50051
posts_service_url: localhost:50052
friends_service_url: localhost:50053
groups_service_url: localhost:50054
jwt_secret: your-secret-key
log_level: info
```

## API Documentation

The Gateway API is documented using Swagger. You can access the Swagger UI at `http://localhost:8000/swagger/index.html` when the Gateway API is running.

## Testing

To run the tests for each service:

```bash
# Users API
cd users-api
go test ./...

# Post API
cd ../post-api
go test ./...

# Friends API
cd ../friends-api
go test ./...

# Groups API
cd ../groups-api
go test ./...

# Gateway API
cd ../gateway-api
go test ./...
```

## Logging

Each service uses structured logging with different log levels (debug, info, warn, error, fatal). The log level can be configured in the configuration file or using environment variables.

Logs are output to stdout in JSON format for easy parsing and analysis.

## Deployment

The application can be deployed using Docker and Docker Compose. A `docker-compose.yml` file is provided in the root directory.

```bash
docker-compose up -d
```

This will start all services and the MySQL database.

```aiignore
protoc --proto_path=./proto \
  --go_out=./pb \
  --go-grpc_out=./pb \
  ./proto/*/*.proto
```
```aiignore
swag init -g cmd/main.go
```

## Setting ENV vars
Change the following variables according to your spec
```aiignore
your-jwt-secret
your-google-client-id
your-google-client-secret
your-microsoft-client-id
your-microsoft-client-secret
your-db-password
```