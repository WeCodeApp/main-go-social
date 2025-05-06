# Friends API Microservice

This is a gRPC microservice for managing friend relationships in a social media application.

## Features

- Send, accept, and reject friend requests
- Get a list of friends
- Remove friends
- Block and unblock users
- Check friendship status between users

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
CREATE DATABASE friends_db;
```

2. Run the migrations:

```bash
# Using golang-migrate
migrate -path database/migrations -database "mysql://root:your-db-password@tcp(localhost:3306)/friends_db" up
```

3. (Optional) Seed the database with sample data:

```bash
cd friends-api
mysql -u root -p friends_db < database/seeds/friends.sql
```

## Running the Service

1. Install dependencies:

```bash
cd friends-api
go mod download
```

2. Build the service:

```bash
go build -o friends-api ./cmd
```

3. Run the service:

```bash
./friends-api
```

Alternatively, you can run the service directly:

```bash
go run cmd/main.go
```

## API Documentation

The Friends API provides the following gRPC methods:

- `SendFriendRequest`: Sends a friend request to another user
- `GetFriendRequests`: Retrieves friend requests for a user
- `AcceptFriendRequest`: Accepts a friend request
- `RejectFriendRequest`: Rejects a friend request
- `GetFriends`: Retrieves friends for a user
- `RemoveFriend`: Removes a friend
- `BlockUser`: Blocks a user
- `UnblockUser`: Unblocks a user
- `GetBlockedUsers`: Retrieves blocked users for a user
- `CheckFriendship`: Checks if two users are friends

## Authentication

The service uses JWT tokens for authentication. All endpoints except `CheckFriendship` require authentication.

## Error Handling

The service returns standard gRPC error codes:

- `UNAUTHENTICATED`: The user is not authenticated
- `NOT_FOUND`: The requested resource was not found
- `INVALID_ARGUMENT`: The request contains invalid arguments
- `INTERNAL`: An internal server error occurred
- `ALREADY_EXISTS`: The resource already exists
- `PERMISSION_DENIED`: The user does not have permission to perform the action