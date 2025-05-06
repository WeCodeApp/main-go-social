# Post API Microservice

The Post API is a gRPC microservice that handles post-related operations for the social media application. It provides functionality for creating, retrieving, updating, and deleting posts, as well as adding, retrieving, and deleting comments, and liking and unliking posts.

## Features

- Create, retrieve, update, and delete posts
- Add, retrieve, and delete comments
- Like and unlike posts
- Visibility rules for posts (public, private, group)
- Authentication and authorization
- Pagination for posts and comments
- MySQL database for persistence

## Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher
- Docker (optional)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/social-media.git
cd social-media/post-api
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up the database:

```bash
# Create the database
mysql -u root -p -e "CREATE DATABASE post_db;"

# Run migrations
migrate -path database/migrations -database "mysql://root:password@tcp(localhost:3306)/post_db" up

# Seed the database (optional)
mysql -u root -p post_db < database/seeds/posts.sql
```

## Configuration

The Post API is configured using a YAML file located at `config/config.yaml`. You can modify this file to change the server, database, JWT, and logging settings.

```yaml
# Server settings
server:
  port: 50052
  host: 0.0.0.0

# Database settings
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  name: post_db
  charset: utf8mb4
  parseTime: true
  loc: Local
  maxOpenConns: 10
  maxIdleConns: 5
  connMaxLifetime: 1h

# JWT settings
jwt:
  secret: your-jwt-secret
  expiration: 24h # 24 hours

# Logging settings
logging:
  level: info # debug, info, warn, error, fatal, panic
  format: json # json or console
  output: stdout # stdout or file
  file: logs/post-api.log # only used if output is file
```

## Running the Service

### Locally

```bash
# Build the service
go build -o post-api cmd/main.go

# Run the service
./post-api
```

### Using Docker

```bash
# Build the Docker image
docker build -t post-api .

# Run the Docker container
docker run -p 50052:50052 post-api
```

## API Documentation

The Post API provides the following gRPC methods:

- `CreatePost`: Creates a new post
- `GetPost`: Retrieves a post by ID
- `GetPosts`: Retrieves posts with pagination and filtering
- `UpdatePost`: Updates a post
- `DeletePost`: Deletes a post
- `AddComment`: Adds a comment to a post
- `GetComments`: Retrieves comments for a post
- `DeleteComment`: Deletes a comment
- `LikePost`: Likes a post
- `UnlikePost`: Unlikes a post

For detailed information about the request and response messages, see the `posts.proto` file in the `common/proto/posts` directory.

## Authentication

The Post API uses JWT tokens for authentication. The token should be included in the `authorization` header of the gRPC request with the format `Bearer <token>`.

The following methods are public and do not require authentication:

- `GetPost`: Retrieves a post by ID (only public posts)
- `GetPosts`: Retrieves posts with pagination and filtering (only public posts)
- `GetComments`: Retrieves comments for a post (only for public posts)

All other methods require authentication.

## Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test ./... -cover
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -am 'Add my feature'`
4. Push to the branch: `git push origin feature/my-feature`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.