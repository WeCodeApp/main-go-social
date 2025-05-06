package routes

import (
	"github.com/gin-gonic/gin"

	"gateway-api/internal/config"
	"gateway-api/internal/controllers"
	"gateway-api/internal/middleware"
	"gateway-api/internal/services"
	"gateway-api/internal/utils/logger"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(router *gin.RouterGroup, cfg *config.Config, logger *logger.Logger) {
	// Create controllers
	userController := controllers.NewUserController(cfg, logger)
	postController := controllers.NewPostController(cfg, logger)
	friendController := controllers.NewFriendController(cfg, logger)
	groupController := controllers.NewGroupController(cfg, logger)

	// Create auth service and controller
	userService := services.NewUserService(cfg, logger)
	authService := services.NewAuthService(cfg, logger, userService)
	authController := controllers.NewAuthController(cfg, logger, authService, userService)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg, logger)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Auth routes
	authRoutes := router.Group("/auth")
	{
		// OAuth routes
		authRoutes.GET("/google", authController.GoogleLogin)
		authRoutes.GET("/google/callback", authController.GoogleCallback)
		authRoutes.GET("/microsoft", authController.MicrosoftLogin)
		authRoutes.GET("/microsoft/callback", authController.MicrosoftCallback)
		authRoutes.POST("/signout", authMiddleware.Authenticate(), authController.Signout)
	}

	// User routes
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
		userRoutes.GET("/me", authMiddleware.Authenticate(), userController.GetProfile)
		userRoutes.PUT("/me", authMiddleware.Authenticate(), userController.UpdateProfile)
	}

	// Post routes
	postRoutes := router.Group("/posts")
	{
		postRoutes.GET("", postController.GetPosts)
		postRoutes.GET("/:id", postController.GetPost)
		postRoutes.POST("", authMiddleware.Authenticate(), postController.CreatePost)
		postRoutes.PUT("/:id", authMiddleware.Authenticate(), postController.UpdatePost)
		postRoutes.DELETE("/:id", authMiddleware.Authenticate(), postController.DeletePost)

		// Comments
		postRoutes.GET("/:id/comments", postController.GetComments)
		postRoutes.POST("/:id/comments", authMiddleware.Authenticate(), postController.AddComment)
		postRoutes.DELETE("/:id/comments/:commentId", authMiddleware.Authenticate(), postController.DeleteComment)

		// Likes
		postRoutes.POST("/:id/like", authMiddleware.Authenticate(), postController.LikePost)
		postRoutes.DELETE("/:id/like", authMiddleware.Authenticate(), postController.UnlikePost)
	}

	// Friend routes
	friendRoutes := router.Group("/friends")
	{
		friendRoutes.GET("", authMiddleware.Authenticate(), friendController.GetFriends)
		friendRoutes.POST("/requests", authMiddleware.Authenticate(), friendController.SendFriendRequest)
		friendRoutes.GET("/requests", authMiddleware.Authenticate(), friendController.GetFriendRequests)
		friendRoutes.PUT("/requests/:id/accept", authMiddleware.Authenticate(), friendController.AcceptFriendRequest)
		friendRoutes.PUT("/requests/:id/reject", authMiddleware.Authenticate(), friendController.RejectFriendRequest)
		friendRoutes.DELETE("/:id", authMiddleware.Authenticate(), friendController.RemoveFriend)
		friendRoutes.POST("/block/:id", authMiddleware.Authenticate(), friendController.BlockUser)
		friendRoutes.DELETE("/block/:id", authMiddleware.Authenticate(), friendController.UnblockUser)
	}

	// Group routes
	groupRoutes := router.Group("/groups")
	{
		groupRoutes.GET("", groupController.GetGroups)
		groupRoutes.GET("/:id", groupController.GetGroup)
		groupRoutes.POST("", authMiddleware.Authenticate(), groupController.CreateGroup)
		groupRoutes.PUT("/:id", authMiddleware.Authenticate(), groupController.UpdateGroup)
		groupRoutes.DELETE("/:id", authMiddleware.Authenticate(), groupController.DeleteGroup)

		// Group membership
		groupRoutes.POST("/:id/members", authMiddleware.Authenticate(), groupController.JoinGroup)
		groupRoutes.DELETE("/:id/members", authMiddleware.Authenticate(), groupController.LeaveGroup)

		// Group posts
		groupRoutes.GET("/:id/posts", groupController.GetGroupPosts)
		groupRoutes.POST("/:id/posts", authMiddleware.Authenticate(), groupController.CreateGroupPost)
	}
}
