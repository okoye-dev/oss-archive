package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/okoye-dev/oss-archive/internal/config"
)

// Server wraps the HTTP server with configuration
type Server struct {
	httpServer *http.Server
	config     *config.Config
}

// New creates a new server instance with the given configuration
func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// SetupRoutes configures all the application routes
func (s *Server) SetupRoutes() *gin.Engine {
	gin.SetMode(s.config.Logging.Mode)
	
	router := gin.Default()
	SetupRoutes(router)

	return router
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
	router := s.SetupRoutes()

	// Create HTTP server with timeouts from config
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.config.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", s.config.Server.Port)
		log.Printf("Timeouts - Read: %ds,   Write: %ds,   Idle: %ds", 
			s.config.Server.ReadTimeout, 
			s.config.Server.WriteTimeout, 
			s.config.Server.IdleTimeout)
		log.Printf("📝 Logging mode: %s", s.config.Logging.Mode)
		
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	return s.waitForShutdown()
}

// waitForShutdown handles graceful shutdown
func (s *Server) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("🛑 Shutdown signal received...")

	// Create shutdown context with timeout from config
	shutdownTimeout := time.Duration(s.config.Server.ShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Printf("...Graceful shutdown timeout: %ds...", s.config.Server.ShutdownTimeout)
	
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("✅ Server shutting down... come back next time :)")
	return nil
}

// Stop stops the server immediately
func (s *Server) Stop() error {
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	return nil
}
