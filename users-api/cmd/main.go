package main

import (
	pb "common/pb/common/proto/users"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"users-api/internal/config"
	"users-api/internal/controllers"
	"users-api/internal/middleware"
	"users-api/internal/repository"
	"users-api/internal/services"
	"users-api/internal/utils/logger"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output, cfg.Logging.File)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Logger.Sync()

	log.Info("Starting Users API")

	// Connect to database
	db, err := gorm.Open(mysql.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database connection", err)
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(
		userRepo,
		log,
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Microsoft.ClientID,
		cfg.OAuth.Microsoft.ClientSecret,
	)

	// Initialize auth service
	authService := services.NewAuthService(
		userRepo,
		log,
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Microsoft.ClientID,
		cfg.OAuth.Microsoft.ClientSecret,
	)

	// Initialize controllers
	userController := controllers.NewUserController(userService, log)
	authController := controllers.NewAuthController(authService, userController, log)

	// Initialize auth interceptor
	authInterceptor := middleware.NewAuthInterceptor(cfg.JWT.Secret, log)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	// Register services
	pb.RegisterUserServiceServer(grpcServer, authController)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		log.Fatal("Failed to listen", err)
	}

	go func() {
		log.Info(fmt.Sprintf("Starting gRPC server on %s:%s", cfg.Server.Host, cfg.Server.Port))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	grpcServer.GracefulStop()
	log.Info("Server exited properly")
}
