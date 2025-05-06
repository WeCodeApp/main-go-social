# Users API

This is the Users API microservice for the social media application. It provides user authentication and profile management functionality.

## Features

- User registration and login using OAuth (Google and Microsoft)
- JWT token generation for authentication
- User profile management
- Secure API endpoints with authentication middleware

## Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher
- Protocol Buffers compiler (protoc)
- Go plugins for Protocol Buffers

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/social-media.git
   cd social-media/users-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate Protocol Buffers code:
   ```bash
   # From the root of the repository
   protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       common/proto/users/users.proto
   ```

4. Create the database:
   ```bash
   mysql -u root -p -e "CREATE DATABASE users_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   ```

5. Run database migrations:
   ```bash
   # Using golang-migrate
   migrate -path database/migrations -database "mysql://root:your-db-password@tcp(localhost:3306)/users_db" up
   
   # Or manually
   mysql -u root -p users_db < database/migrations/000001_create_users_table.up.sql
   ```

6. (Optional) Seed the database with sample data:
   ```bash
   mysql -u root -p users_db < database/seeds/users.sql
   ```

## Configuration

Update the `config/config.yaml` file with your settings:

```yaml
# Server settings
server:
  port: 50051
  host: 0.0.0.0

# Database settings
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: your-password-here
  name: users_db
  charset: utf8mb4
  parseTime: true
  loc: Local
  maxOpenConns: 10
  maxIdleConns: 5
  connMaxLifetime: 1h

# JWT settings
jwt:
  secret: your-secret-key-here
  expiration: 24h # 24 hours

# OAuth settings
oauth:
  google:
    clientID: your-google-client-id
    clientSecret: your-google-client-secret
    redirectURL: http://localhost:8000/api/v1/users/oauth/google/callback
  microsoft:
    clientID: your-microsoft-client-id
    clientSecret: your-microsoft-client-secret
    redirectURL: http://localhost:8000/api/v1/users/oauth/microsoft/callback

# Logging settings
logging:
  level: info # debug, info, warn, error, fatal, panic
  format: json # json or console
  output: stdout # stdout or file
  file: logs/users-api.log # only used if output is file
```

## Running the Service

1. Start the service:
   ```bash
   go run cmd/main.go
   ```

2. Or build and run:
   ```bash
   go build -o users-api cmd/main.go
   ./users-api
   ```

The service will start on the configured port (default: 50051).

## API Endpoints

The Users API provides the following gRPC endpoints:

- `Register`: Register a new user with OAuth provider
- `Login`: Authenticate a user with OAuth provider
- `GetProfile`: Retrieve a user's profile
- `UpdateProfile`: Update a user's profile

## Testing

Run the tests:
```bash
go test ./...
```

## Logging

Logs are written to stdout by default. You can configure the logging level, format, and output in the config file.

## Monitoring

The service includes basic logging for monitoring. In a production environment, you would integrate with a monitoring system like Prometheus and Grafana.

## Troubleshooting

- If you encounter database connection issues, check your database settings in the config file.
- If authentication fails, verify your OAuth settings and JWT secret.
- For detailed logs, set the logging level to "debug" in the config file.

## License

This project is licensed under the MIT License - see the LICENSE file for details.