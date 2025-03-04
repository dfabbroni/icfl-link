package routes

import (
	"link/internal/config"
	"link/internal/handlers"
	"link/internal/middleware"
	"link/internal/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB, config *config.Config, pythonEnv *utils.PythonEnv) {
	e.Use(middleware.ErrorHandler)

	nodeHandler := &handlers.NodeHandler{DB: db}
	userHandler := &handlers.UserHandler{DB: db, Config: config}
	experimentHandler := &handlers.ExperimentHandler{DB: db, Config: config, PythonEnv: pythonEnv}
	metadataHandler := &handlers.MetadataHandler{DB: db}
	fileHandler := &handlers.FileHandler{}

	// Public routes
	e.POST("/nodes", nodeHandler.RegisterNode)
	e.POST("/nodes/login", nodeHandler.LoginNode)
	e.POST("/users", userHandler.RegisterUser)
	e.POST("/users/login", userHandler.Login)

	// Protected routes
	r := e.Group("/api")
	r.Use(middleware.CombinedAuthMiddleware(db, config.Auth.SecretKey))

	// Node routes
	r.PUT("/nodes/:id/accept", nodeHandler.AcceptNode)
	r.DELETE("/nodes/:id", nodeHandler.RejectNode)
	r.PUT("/nodes/status", nodeHandler.UpdateNodeStatus)
	r.GET("/nodes", nodeHandler.ListNodes)
	r.GET("/node/instructions", nodeHandler.PollInstructions)

	// Experiment routes
	r.POST("/experiments", experimentHandler.CreateExperiment)
	r.PUT("/experiments/:experimentID/accept", experimentHandler.AcceptExperiment)
	r.PUT("/experiments/:experimentID/reject", experimentHandler.RejectExperiment)
	r.POST("/experiments/:id/start", experimentHandler.StartTraining)
	r.POST("/experiments/:id/stop", experimentHandler.StopTraining)
	r.GET("/experiments", experimentHandler.ListExperiments)
	r.PUT("/experiments/:id", experimentHandler.UpdateExperiment)
	r.POST("/experiments/:experimentID/node-start", experimentHandler.NodeTrainingStarted)
	r.POST("/experiments/:experimentID/checksum", experimentHandler.ReceiveChecksum)
	r.POST("/experiments/:experimentID/update-files", experimentHandler.UpdateFiles)

	// Metadata routes
	r.POST("/metadata", metadataHandler.RegisterMetadata)
	r.GET("/metadata", metadataHandler.FetchMetadata)

	// File routes
	r.GET("/download", fileHandler.DownloadFile)
}
