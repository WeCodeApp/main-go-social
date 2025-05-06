package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"gateway-api/internal/config"
	"gateway-api/internal/utils/logger"
)

// AuthMiddleware handles authentication and authorization
type AuthMiddleware struct {
	cfg    *config.Config
	logger *logger.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(cfg *config.Config, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:    cfg,
		logger: logger,
	}
}

// Authenticate verifies the JWT token in the Authorization header
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be Bearer {token}",
			})
			return
		}

		// Parse and validate the token
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(m.cfg.JWTSecret), nil
		})

		if err != nil {
			m.logger.Error("Failed to parse token", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add claims to the context
			userID, ok := claims["sub"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token claims",
				})
				return
			}

			// Set user ID in context
			c.Set("userID", userID)

			// Also set the JWT token in context
			c.Set("jwt_token", tokenString)

			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}
	}
}
