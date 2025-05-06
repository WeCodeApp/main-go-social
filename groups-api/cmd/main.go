package main

import (
	pb "common/pb/common/proto/groups"
	"fmt"
	"groups-api/internal/config"
	"groups-api/internal/controllers"
	"groups-api/internal/middleware"
	"groups-api/internal/repository"
	"groups-api/internal/services"
	"groups-api/internal/utils/logger"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	log.Info("Starting Groups API")

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
	groupRepo := repository.NewGroupRepository(db)

	// Initialize services
	groupService := services.NewGroupService(groupRepo, log)

	// Initialize controllers
	groupController := controllers.NewGroupController(groupService, log)

	// Initialize auth interceptor
	authInterceptor := middleware.NewAuthInterceptor(cfg.JWT.Secret, log)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	// Register services
	pb.RegisterGroupServiceServer(grpcServer, groupController)

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
