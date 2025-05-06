package middleware

import (
	"context"
	"errors"
	"groups-api/internal/utils/logger"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor is a gRPC interceptor for authentication
type AuthInterceptor struct {
	jwtSecret     string
	logger        *logger.Logger
	publicMethods map[string]bool
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(jwtSecret string, logger *logger.Logger) *AuthInterceptor {
	return &AuthInterceptor{
		jwtSecret: jwtSecret,
		logger:    logger,
		publicMethods: map[string]bool{
			"/groups.GroupService/GetGroups": true,
		},
	}
}

// Unary returns a unary server interceptor for authentication
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for public methods
		if i.publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Authenticate the request
		userID, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		// Add user ID to context
		ctx = context.WithValue(ctx, "userID", userID)

		// Proceed with the request
		return handler(ctx, req)
	}
}

// authenticate authenticates the request
func (i *AuthInterceptor) authenticate(ctx context.Context) (string, error) {
	// Get metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	// Get authorization header
	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	// Extract token from authorization header
	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(i.jwtSecret), nil
	})
	if err != nil {
		i.logger.Error("Failed to parse token", err)
		return "", status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	// Get claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", status.Errorf(codes.Unauthenticated, "invalid token claims")
	}

	// Check token expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "invalid token expiration")
	}
	if time.Now().Unix() > int64(exp) {
		return "", status.Errorf(codes.Unauthenticated, "token expired")
	}

	// Get user ID from claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "invalid user ID in token")
	}

	return userID, nil
}