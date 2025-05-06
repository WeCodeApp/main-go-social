package middleware

import (
	"context"
	"strings"
	"time"
	"post-api/internal/utils/logger"

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
			"/posts.PostService/GetPost":     true,
			"/posts.PostService/GetPosts":    true,
			"/posts.PostService/GetComments": true,
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
		userID, friendIDs, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		// Add user ID and friend IDs to context
		ctx = context.WithValue(ctx, "user_id", userID)
		if len(friendIDs) > 0 {
			ctx = context.WithValue(ctx, "friend_ids", friendIDs)
		}

		// Proceed with the request
		return handler(ctx, req)
	}
}

// authenticate authenticates the request
func (i *AuthInterceptor) authenticate(ctx context.Context) (string, []string, error) {
	// Get metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// Get authorization header
	values := md.Get("authorization")
	if len(values) == 0 {
		return "", nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	// Extract token from authorization header
	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", nil, status.Error(codes.Unauthenticated, "invalid authorization format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate token with claims
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unauthenticated, "unexpected signing method")
		}
		return []byte(i.jwtSecret), nil
	})
	if err != nil {
		i.logger.Error("Failed to parse token", err)
		return "", nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
	}

	// Check if token is valid
	if !token.Valid {
		return "", nil, status.Error(codes.Unauthenticated, "invalid token claims")
	}

	// Check token expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", nil, status.Error(codes.Unauthenticated, "invalid token expiration")
	}
	if time.Now().Unix() > int64(exp) {
		return "", nil, status.Error(codes.Unauthenticated, "token expired")
	}

	// Get user ID from claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", nil, status.Error(codes.Unauthenticated, "invalid user ID in token")
	}

	// Get friend IDs from metadata
	friendIDs := md.Get("friend_ids")

	return userID, friendIDs, nil
}
