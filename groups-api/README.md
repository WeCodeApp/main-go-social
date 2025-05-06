# Groups API Microservice

This is a gRPC microservice for managing groups in a social media application.

## Features

- Create, read, update, and delete groups
- Join and leave groups
- Create posts in groups
- Like and comment on group posts

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher

## Configuration

The service is configured using the `config.yaml` file in the `config` directory. You can modify the following settings:

- Server host and port
- Database connection details
- JWT secret and expiration
- Logging settings

## Database Setup

1. Create a MySQL database:

```sql
CREATE DATABASE groups_db;
```

2. Run the migrations:

```bash
# Using golang-migrate
migrate -path database/migrations -database "mysql://root:your-db-password@tcp(localhost:3306)/groups_db" up
```

3. (Optional) Seed the database with sample data:

```bash
cd groups-api
mysql -u root -p groups_db < database/seeds/groups.sql
```

## Running the Service

1. Install dependencies:

```bash
cd groups-api
go mod download
```

2. Build the service:

```bash
go build -o groups-api ./cmd
```

3. Run the service:

```bash
./groups-api
```

Alternatively, you can run the service directly:

```bash
go run cmd/main.go
```

## API Documentation

The Groups API provides the following gRPC methods:

- `CreateGroup`: Creates a new group
- `GetGroup`: Retrieves a group by ID
- `GetGroups`: Retrieves groups with pagination and filtering
- `UpdateGroup`: Updates a group
- `DeleteGroup`: Deletes a group
- `JoinGroup`: Adds a user to a group
- `LeaveGroup`: Removes a user from a group
- `GetGroupMembers`: Retrieves members of a group
- `CreateGroupPost`: Creates a post in a group
- `GetGroupPosts`: Retrieves posts in a group

## Authentication

The service uses JWT tokens for authentication. All endpoints except `GetGroups` require authentication.

## Error Handling

The service returns standard gRPC error codes:

- `UNAUTHENTICATED`: The user is not authenticated
- `NOT_FOUND`: The requested resource was not found
- `INVALID_ARGUMENT`: The request contains invalid arguments
- `INTERNAL`: An internal server error occurred
- `ALREADY_EXISTS`: The resource already exists
- `PERMISSION_DENIED`: The user does not have permission to perform the action