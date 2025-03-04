package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"link/internal/config"
	"link/internal/database"
	"link/internal/routes"
	"link/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize experiment environment (creates venv if needed)
	if err := utils.InitializeExperimentEnvironment(cfg.Server.PythonEnvPath); err != nil {
		log.Fatalf("Failed to initialize experiment environment: %v", err)
	}

	// Get Python environment
	pythonEnv, err := utils.GetSharedPythonEnv(cfg.Server.PythonEnvPath)
	if err != nil {
		log.Fatalf("Failed to get Python environment: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(e, db, cfg, pythonEnv)

	go func() {
		if err := e.Start(cfg.Server.Port); err != nil {
			log.Printf("Shutting down the server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	// Cleanup Python processes
	if err := pythonEnv.CleanupSuperLink(); err != nil {
		log.Printf("Error during SuperLink cleanup: %v", err)
	}

	if err := pythonEnv.CleanupFlwr(); err != nil {
		log.Printf("Error during Flwr cleanup: %v", err)
	}

	log.Println("Server stopped")
}
