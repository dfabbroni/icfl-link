package handlers

import (
	"encoding/json"
	"fmt"
	"link/internal/config"
	"link/internal/models"
	"link/internal/store"
	"link/internal/utils"

	"crypto/sha256"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var experimentMutex sync.Mutex

type ExperimentHandler struct {
	DB        *gorm.DB
	Config    *config.Config
	PythonEnv *utils.PythonEnv
}

func (h *ExperimentHandler) CreateExperiment(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(50 << 20); err != nil { // 50 MB max
		log.Printf("Error parsing multipart form: %v\n", err)
		return utils.NewBadRequestError("Failed to parse form data")
	}

	experiment := new(models.Experiment)
	experiment.Name = c.FormValue("name")
	experiment.Description = c.FormValue("description")
	experiment.Status = c.FormValue("status")
	userID := uint(c.Get("user_id").(float64))
	experiment.UserID = userID

	if err := h.DB.Create(experiment).Error; err != nil {
		log.Printf("Error creating experiment: %v\n", err)
		return utils.NewInternalServerError("Failed to create experiment")
	}

	if err := h.handleFileUploads(c, experiment); err != nil {
		return err
	}

	if err := h.createExperimentNodes(c, experiment); err != nil {
		return err
	}

	return c.JSON(201, experiment)
}

func (h *ExperimentHandler) handleFileUploads(c echo.Context, experiment *models.Experiment) error {
	zipFile, err := c.FormFile("experimentFiles")
	if err != nil {
		return utils.NewBadRequestError("Failed to get experiment files")
	}

	// Create experiment directory
	experimentDir := filepath.Join("uploads", fmt.Sprintf("%d", experiment.ID))
	if err := os.MkdirAll(experimentDir, 0755); err != nil {
		return utils.NewInternalServerError("Failed to create experiment directory")
	}

	// Extract zip contents
	if err := utils.ExtractZipFile(zipFile, experimentDir); err != nil {
		return utils.NewInternalServerError(fmt.Sprintf("Failed to extract zip file: %v", err))
	}

	// Find the experiment name folder
	entries, err := os.ReadDir(experimentDir)
	if err != nil {
		return utils.NewInternalServerError("Failed to read experiment directory")
	}

	var experimentNameDir string
	for _, entry := range entries {
		if entry.IsDir() {
			experimentNameDir = entry.Name()
			break
		}
	}

	if experimentNameDir == "" {
		return utils.NewBadRequestError("Invalid zip structure: missing experiment folder")
	}

	// Verify the inner folder structure
	innerPath := filepath.Join(experimentDir, experimentNameDir, experimentNameDir)
	if _, err := os.Stat(innerPath); os.IsNotExist(err) {
		return utils.NewBadRequestError("Invalid zip structure: missing inner folder")
	}

	// Verify required files exist
	requiredFiles := []string{
		filepath.Join(experimentDir, experimentNameDir, "pyproject.toml"),
		filepath.Join(experimentDir, experimentNameDir, experimentNameDir, "client_app.py"),
		filepath.Join(experimentDir, experimentNameDir, experimentNameDir, "server_app.py"),
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return utils.NewBadRequestError(fmt.Sprintf("Missing required file: %s", filepath.Base(file)))
		}
	}

	// Update experiment with base path
	experiment.BasePath = experimentDir + "/" + experimentNameDir
	if err := h.DB.Save(experiment).Error; err != nil {
		return utils.NewInternalServerError("Failed to update experiment with file path")
	}

	return nil
}

func (h *ExperimentHandler) createExperimentNodes(c echo.Context, experiment *models.Experiment) error {
	selectedNodesJSON := c.FormValue("selectedNodes")
	var selectedNodes []struct {
		ID         uint `json:"id"`
		NodeID     uint `json:"node_id"`
		MetadataID uint `json:"metadata_id"`
	}
	if err := json.Unmarshal([]byte(selectedNodesJSON), &selectedNodes); err != nil {
		log.Printf("Error unmarshalling selectedNodes: %v\n", err)
		return utils.NewBadRequestError("Invalid node selection data")
	}

	experimentNodes := make([]models.ExperimentNode, len(selectedNodes))
	for i, trio := range selectedNodes {
		experimentNodes[i] = models.ExperimentNode{
			ExperimentID: experiment.ID,
			NodeID:       trio.NodeID,
			MetadataID:   trio.ID,
			Status:       models.ExperimentNodeStatusPending,
		}
	}

	if err := h.DB.Create(&experimentNodes).Error; err != nil {
		log.Printf("Error creating experiment nodes: %v\n", err)
		return utils.NewInternalServerError("Failed to create experiment nodes")
	}

	instructions := make([]store.NodeInstruction, len(selectedNodes))
	for i, trio := range selectedNodes {
		instructions[i] = store.NodeInstruction{
			NodeID: trio.NodeID,
			Instruction: models.Instruction{
				Type: models.InstructionNewExperiment,
				Payload: map[string]interface{}{
					"experiment_id": experiment.ID,
					"name":          experiment.Name,
					"description":   experiment.Description,
					"files_path":    experiment.BasePath,
					"metadata_id":   trio.MetadataID,
				},
			},
		}
	}
	store.GlobalInstructionStore.AddInstructions(instructions)

	return nil
}

func (h *ExperimentHandler) AcceptExperiment(c echo.Context) error {
	return h.updateExperimentNodeStatus(c, models.ExperimentNodeStatusAccepted)
}

func (h *ExperimentHandler) RejectExperiment(c echo.Context) error {
	return h.updateExperimentNodeStatus(c, models.ExperimentNodeStatusRejected)
}

func (h *ExperimentHandler) updateExperimentNodeStatus(c echo.Context, status models.ExperimentNodeStatus) error {
	experimentID := c.Param("experimentID")
	node := c.Get("node").(models.Node)

	var experimentNode models.ExperimentNode
	if err := h.DB.Where("experiment_id = ? AND node_id = ?", experimentID, node.ID).First(&experimentNode).Error; err != nil {
		return utils.NewNotFoundError("Experiment node not found")
	}

	experimentNode.Status = status
	if err := h.DB.Save(&experimentNode).Error; err != nil {
		return utils.NewInternalServerError("Failed to update experiment node status")
	}

	return c.JSON(200, experimentNode)
}

func (h *ExperimentHandler) ListExperiments(c echo.Context) error {
	var experiments []models.Experiment
	if err := h.DB.Preload("ExperimentNodes.Metadata").
		Preload("ExperimentNodes.Node", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Password", "PublicKey")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Password")
		}).
		Find(&experiments).Error; err != nil {
		return utils.NewInternalServerError("Failed to fetch experiments")
	}

	return c.JSON(200, experiments)
}

func (h *ExperimentHandler) StartTraining(c echo.Context) error {
	experimentMutex.Lock()
	defer experimentMutex.Unlock()

	return h.DB.Transaction(func(tx *gorm.DB) error {
		// Check if there is already an experiment in preparing or training state
		var activeCount int64
		if err := tx.Model(&models.Experiment{}).
			Where("status IN (?)", []string{string(models.ExperimentNodeStatusPreparing), string(models.ExperimentNodeStatusTraining)}).
			Count(&activeCount).Error; err != nil {
			return utils.NewInternalServerError("Failed to check active experiments")
		}

		if activeCount > 0 {
			return utils.NewBadRequestError("Another experiment is already in progress")
		}

		experimentID := c.Param("id")

		var experiment models.Experiment
		if err := tx.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Password")
		}).
			Preload("ExperimentNodes.Node", func(db *gorm.DB) *gorm.DB {
				return db.Omit("Password", "PublicKey")
			}).
			First(&experiment, experimentID).Error; err != nil {
			return utils.NewNotFoundError("Experiment not found")
		}

		// Update experiment nodes status and create instructions only for accepted nodes
		var experimentNodes []models.ExperimentNode
		if err := tx.Where("experiment_id = ? AND status = ?", experiment.ID, models.ExperimentNodeStatusAccepted).Find(&experimentNodes).Error; err != nil {
			return utils.NewInternalServerError("Failed to fetch experiment nodes")
		}

		if len(experimentNodes) == 0 {
			return utils.NewBadRequestError("No nodes have accepted this experiment")
		}

		instructions := make([]store.NodeInstruction, len(experimentNodes))
		for i, en := range experimentNodes {
			en.Status = models.ExperimentNodeStatusPreparing
			if err := tx.Save(&en).Error; err != nil {
				return utils.NewInternalServerError("Failed to update experiment node status")
			}

			instructions[i] = store.NodeInstruction{
				NodeID: en.NodeID,
				Instruction: models.Instruction{
					Type:    models.InstructionStartTraining,
					Payload: map[string]interface{}{"experiment_id": experiment.ID},
				},
			}
		}

		experiment.Status = string(models.ExperimentNodeStatusPreparing)
		if err := tx.Save(&experiment).Error; err != nil {
			return utils.NewInternalServerError("Failed to update experiment status")
		}

		if err := h.writeNodeKeysToCSV(experiment.ID); err != nil {
			return utils.NewInternalServerError(fmt.Sprintf("Failed to write node keys to CSV: %v", err))
		}

		if err := h.PythonEnv.InitializeSuperLink(); err != nil {
			log.Fatalf("Failed to initialize SuperLink: %v", err)
		}

		store.GlobalInstructionStore.AddInstructions(instructions)

		return c.JSON(200, experiment)
	})
}

func (h *ExperimentHandler) writeNodeKeysToCSV(experimentID uint) error {
	csvFilePath := "authentication/keys/client_public_keys.csv"

	var experimentNodes []models.ExperimentNode
	if err := h.DB.Preload("Node", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "public_key")
	}).Where("experiment_id = ?", experimentID).Find(&experimentNodes).Error; err != nil {
			return fmt.Errorf("failed to fetch experiment nodes: %w", err)
	}

	file, err := os.OpenFile(csvFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
			return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	for _, node := range experimentNodes {
			if _, err := file.WriteString(node.Node.PublicKey + "\n"); err != nil {
					return fmt.Errorf("failed to write to CSV file: %w", err)
			}
	}

	return nil
}

func (h *ExperimentHandler) NodeTrainingStarted(c echo.Context) error {
	experimentMutex.Lock()
	defer experimentMutex.Unlock()

	experimentID := c.Param("experimentID")
	node := c.Get("node").(models.Node)

	var experimentNode models.ExperimentNode
	if err := h.DB.Where("experiment_id = ? AND node_id = ?", experimentID, node.ID).First(&experimentNode).Error; err != nil {
		return utils.NewNotFoundError("Experiment node not found")
	}

	if experimentNode.Status == models.ExperimentNodeStatusTraining {
		return c.JSON(200, map[string]string{"status": "already acknowledged"})
	}

	experimentNode.Status = models.ExperimentNodeStatusTraining
	if err := h.DB.Save(&experimentNode).Error; err != nil {
		return utils.NewInternalServerError("Failed to update experiment node status")
	}

	// Check if all nodes have started training
	var count int64
	if err := h.DB.Model(&models.ExperimentNode{}).
		Where("experiment_id = ? AND status = ?", experimentID, models.ExperimentNodeStatusTraining).
		Count(&count).Error; err != nil {
		return utils.NewInternalServerError("Failed to count training nodes")
	}

	var totalNodes int64
	if err := h.DB.Model(&models.ExperimentNode{}).
		Where("experiment_id = ?", experimentID).
		Count(&totalNodes).Error; err != nil {
		return utils.NewInternalServerError("Failed to count total nodes")
	}

	if count == totalNodes {
		// All nodes have started training, start the server process
		h.startServerProcess(experimentID)

		// Update experiment status
		if err := h.DB.Model(&models.Experiment{}).
			Where("id = ?", experimentID).
			Update("status", models.ExperimentNodeStatusTraining).Error; err != nil {
			return utils.NewInternalServerError("Failed to update experiment status")
		}
	}

	return c.JSON(200, map[string]string{"status": "acknowledged"})
}

func (h *ExperimentHandler) startServerProcess(experimentID string) {
	var experiment models.Experiment
	if err := h.DB.First(&experiment, experimentID).Error; err != nil {
		log.Printf("Failed to find experiment %s: %v", experimentID, err)
		return
	}

	// Get the experiment name
	parts := strings.Split(experiment.BasePath, "/")
	experimentName := parts[len(parts)-1]

	if err := h.PythonEnv.InstallExperimentDependencies(experiment.BasePath); err != nil {
		log.Printf("Failed to install experiment dependencies: %v", err)
		return
	}

	experimentIDStr := fmt.Sprintf("%d", experiment.ID)

	if err := h.PythonEnv.RunFlwr(experiment.BasePath, experimentIDStr, experimentName); err != nil {
		log.Printf("Failed to run FLWR for experiment %s: %v", experimentID, err)
		return
	}

	log.Printf("Starting server process for experiment ID: %s", experimentID)
}

func (h *ExperimentHandler) StopTraining(c echo.Context) error {
	return h.DB.Transaction(func(tx *gorm.DB) error {
		experimentID := c.Param("id")

		// Stop the server process for this experiment
		h.stopServerProcess(experimentID)

		var experiment models.Experiment
		if err := tx.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Password")
		}).
			Preload("ExperimentNodes.Node", func(db *gorm.DB) *gorm.DB {
				return db.Omit("Password", "PublicKey")
			}).
			First(&experiment, experimentID).Error; err != nil {
			return utils.NewNotFoundError("Experiment not found")
		}

		if experiment.Status != string(models.ExperimentNodeStatusTraining) && experiment.Status != string(models.ExperimentNodeStatusPreparing) {
			return utils.NewBadRequestError("Experiment is not currently in training or preparing")
		}

		// Update experiment nodes status and create instructions
		var experimentNodes []models.ExperimentNode
		if err := tx.Where("experiment_id = ? AND status = ?", experiment.ID, models.ExperimentNodeStatusTraining).Find(&experimentNodes).Error; err != nil {
			return utils.NewInternalServerError("Failed to fetch experiment nodes")
		}

		instructions := make([]store.NodeInstruction, len(experimentNodes))
		for i, en := range experimentNodes {
			en.Status = models.ExperimentNodeStatusStopped
			if err := tx.Save(&en).Error; err != nil {
				return utils.NewInternalServerError("Failed to update experiment node status")
			}

			instructions[i] = store.NodeInstruction{
				NodeID: en.NodeID,
				Instruction: models.Instruction{
					Type:    models.InstructionStopTraining,
					Payload: map[string]interface{}{"experiment_id": experiment.ID},
				},
			}
		}
		store.GlobalInstructionStore.AddInstructions(instructions)

		// Update experiment status
		experiment.Status = string(models.ExperimentNodeStatusStopped)
		if err := tx.Save(&experiment).Error; err != nil {
			return utils.NewInternalServerError("Failed to update experiment status")
		}

		return c.JSON(200, experiment)
	})
}

func (h *ExperimentHandler) stopServerProcess(experimentID string) {
	h.PythonEnv.CleanupSuperLink()
	h.PythonEnv.CleanupFlwr()
	fmt.Printf("Stopping server process for experiment ID: %s\n", experimentID)
}

func (h *ExperimentHandler) UpdateExperiment(c echo.Context) error {
	experimentID := c.Param("id")

	return h.DB.Transaction(func(tx *gorm.DB) error {
		var experiment models.Experiment
		if err := tx.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Password")
		}).
			Preload("ExperimentNodes.Node", func(db *gorm.DB) *gorm.DB {
				return db.Omit("Password", "PublicKey")
			}).
			First(&experiment, experimentID).Error; err != nil {
			return utils.NewNotFoundError("Experiment not found")
		}

		// Check if the experiment is in TRAINING status
		if experiment.Status == string(models.ExperimentNodeStatusTraining) ||
			experiment.Status == string(models.ExperimentNodeStatusPreparing) {
			return utils.NewBadRequestError("Cannot update experiment while it is in training or preparing")
		}

		// Bind new experiment data
		if err := c.Bind(&experiment); err != nil {
			return utils.NewBadRequestError("Invalid request payload")
		}

		updatedFiles, err := h.handleFileUpdates(c, &experiment)
		if err != nil {
			return err
		}

		if err := tx.Save(&experiment).Error; err != nil {
			return utils.NewInternalServerError("Failed to update experiment")
		}

		if len(updatedFiles) > 0 {
			if err := h.notifyNodesAboutUpdate(tx, &experiment, updatedFiles); err != nil {
				return err
			}
		}

		return c.JSON(200, experiment)
	})
}

func (h *ExperimentHandler) handleFileUpdates(c echo.Context, experiment *models.Experiment) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, utils.NewBadRequestError("Failed to parse multipart form")
	}

	var updatedFiles []string

	// Get experiment name from BasePath
	parts := strings.Split(experiment.BasePath, "/")
	if len(parts) < 3 { // Should have: uploads/id/expname
		return nil, utils.NewInternalServerError("Invalid base path format")
	}
	experimentName := parts[len(parts)-1]

	// Handle each uploaded file
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			filename := filepath.Base(fileHeader.Filename)
			var filePath string

			// Determine if file should go in base folder or inner folder
			if filename == "pyproject.toml" {
				filePath = filepath.Join(experiment.BasePath, filename)
			} else {
				filePath = filepath.Join(experiment.BasePath, experimentName, filename)
			}

			// Open the uploaded file
			src, err := fileHeader.Open()
			if err != nil {
				return nil, utils.NewInternalServerError(fmt.Sprintf("Failed to open uploaded file %s", filename))
			}
			defer src.Close()

			// Create or truncate the destination file
			dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return nil, utils.NewInternalServerError(fmt.Sprintf("Failed to create destination file %s", filename))
			}
			defer dst.Close()

			// Copy the contents
			if _, err = io.Copy(dst, src); err != nil {
				return nil, utils.NewInternalServerError(fmt.Sprintf("Failed to copy file contents for %s", filename))
			}

			updatedFiles = append(updatedFiles, filename)
		}
	}

	return updatedFiles, nil
}

func (h *ExperimentHandler) notifyNodesAboutUpdate(tx *gorm.DB, experiment *models.Experiment, updatedFiles []string) error {
	var experimentNodes []models.ExperimentNode
	if err := tx.Where("experiment_id = ?", experiment.ID).Find(&experimentNodes).Error; err != nil {
		return utils.NewInternalServerError("Failed to fetch experiment nodes")
	}

	// Create instructions for each node
	for _, en := range experimentNodes {
		instruction := store.NodeInstruction{
			NodeID: en.NodeID,
			Instruction: models.Instruction{
				Type: models.InstructionUpdateExperiment,
				Payload: map[string]interface{}{
					"experiment_id": experiment.ID,
					"files_path":    experiment.BasePath,
					"updated_files": updatedFiles,
				},
			},
		}
		store.GlobalInstructionStore.AddInstructions([]store.NodeInstruction{instruction})

		en.Status = models.ExperimentNodeStatusPending
		if err := tx.Save(&en).Error; err != nil {
			return utils.NewInternalServerError("Failed to update experiment node status")
		}
	}

	return nil
}

func (h *ExperimentHandler) ReceiveChecksum(c echo.Context) error {
	experimentID := c.Param("experimentID")
	nodeID := c.Get("node").(models.Node).ID
	clientAppChecksum := c.FormValue("client_app_checksum")
	pyprojectChecksum := c.FormValue("pyproject_checksum")

	fmt.Printf("Received checksum for experiment %s: client_app_checksum=%s, pyproject_checksum=%s\n",
		experimentID, clientAppChecksum, pyprojectChecksum)

	var experiment models.Experiment
	if err := h.DB.First(&experiment, experimentID).Error; err != nil {
		return utils.NewNotFoundError("Experiment not found")
	}

	// Get file paths
	parts := strings.Split(experiment.BasePath, "/")
	experimentName := parts[len(parts)-1]
	pyprojectPath := filepath.Join(experiment.BasePath, "pyproject.toml")
	clientAppPath := filepath.Join(experiment.BasePath, experimentName, "client_app.py")

	// Calculate SHA-256 for pyproject.toml
	pyprojectFile, err := os.Open(pyprojectPath)
	if err != nil {
		return utils.NewInternalServerError(fmt.Sprintf("Failed to open pyproject.toml: %v", err))
	}
	defer pyprojectFile.Close()

	pyprojectHasher := sha256.New()
	if _, err := io.Copy(pyprojectHasher, pyprojectFile); err != nil {
		return utils.NewInternalServerError(fmt.Sprintf("Failed to calculate pyproject.toml hash: %v", err))
	}
	calculatedPyprojectHash := fmt.Sprintf("%x", pyprojectHasher.Sum(nil))

	// Calculate SHA-256 for client_app.py
	clientAppFile, err := os.Open(clientAppPath)
	if err != nil {
		return utils.NewInternalServerError(fmt.Sprintf("Failed to open client_app.py: %v", err))
	}
	defer clientAppFile.Close()

	clientAppHasher := sha256.New()
	if _, err := io.Copy(clientAppHasher, clientAppFile); err != nil {
		return utils.NewInternalServerError(fmt.Sprintf("Failed to calculate client_app.py hash: %v", err))
	}
	calculatedClientAppHash := fmt.Sprintf("%x", clientAppHasher.Sum(nil))

	// Compare checksums
	if calculatedPyprojectHash != pyprojectChecksum || calculatedClientAppHash != clientAppChecksum {
		if err := h.DB.Model(&models.ExperimentNode{}).
			Where("experiment_id = ? AND node_id = ?", experimentID, nodeID).
			Update("status", models.ExperimentNodeStatusChecksumMismatch).Error; err != nil {
			return utils.NewInternalServerError("Failed to update experiment node status")
		}

		return utils.NewBadRequestError("checksum mismatch")
	}

	return c.JSON(200, map[string]string{
		"status": "verified",
	})
}

func (h *ExperimentHandler) UpdateFiles(c echo.Context) error {
	experimentID := c.Param("experimentID")
	updatedFiles := []string{"pyproject.toml", "client_app.py"}

	// Grab the nodes that have a checksum mismatch
	var experimentNodes []models.ExperimentNode
	if err := h.DB.Where("experiment_id = ? AND status = ?", experimentID, models.ExperimentNodeStatusChecksumMismatch).Find(&experimentNodes).Error; err != nil {
		return utils.NewInternalServerError("Failed to fetch experiment nodes")
	}

	for _, en := range experimentNodes {
		instruction := store.NodeInstruction{
			NodeID: en.NodeID,
			Instruction: models.Instruction{
				Type: models.InstructionUpdateExperiment,
				Payload: map[string]interface{}{
					"experiment_id": en.ExperimentID,
					"updated_files": updatedFiles,
				},
			},
		}
		store.GlobalInstructionStore.AddInstructions([]store.NodeInstruction{instruction})
	}

	return c.JSON(200, map[string]string{"status": "acknowledged"})
}
