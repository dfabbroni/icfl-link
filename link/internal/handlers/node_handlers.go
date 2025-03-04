package handlers

import (
	"time"

	"link/internal/models"
	"link/internal/store"
	"link/internal/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type NodeHandler struct {
	DB *gorm.DB
}

func (h *NodeHandler) RegisterNode(c echo.Context) error {
	node := new(models.Node)
	if err := c.Bind(node); err != nil {
		return utils.NewBadRequestError("Invalid request payload")
	}

	if node.Username == "" || node.Password == "" || node.PublicKey == "" {
		return utils.NewBadRequestError("Username, password, and public key are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(node.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewInternalServerError("Failed to hash password")
	}
	node.Password = string(hashedPassword)

	node.Approved = false
	node.LastSeen = time.Now()

	if err := h.DB.Create(node).Error; err != nil {
		return utils.NewInternalServerError("Failed to register node")
	}

	node.Password = ""
	return c.JSON(201, node)
}

func (h *NodeHandler) LoginNode(c echo.Context) error {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&credentials); err != nil {
		return utils.NewBadRequestError("invalid request body")
	}

	var node models.Node
	if err := h.DB.Where("username = ?", credentials.Username).First(&node).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewUnauthorizedError("invalid credentials")
		}
		return utils.NewInternalServerError("database error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(node.Password), []byte(credentials.Password)); err != nil {
		return utils.NewUnauthorizedError("invalid credentials")
	}

	return c.JSON(201, node)
}

func (h *NodeHandler) AcceptNode(c echo.Context) error {
	id := c.Param("id")

	var node models.Node
	if err := h.DB.First(&node, id).Error; err != nil {
		return utils.NewNotFoundError("Node not found")
	}

	node.Approved = true
	if err := h.DB.Save(&node).Error; err != nil {
		return utils.NewInternalServerError("Failed to accept node")
	}

	return c.JSON(200, node)
}

func (h *NodeHandler) RejectNode(c echo.Context) error {
	id := c.Param("id")

	if err := h.DB.Delete(&models.Node{}, id).Error; err != nil {
		return utils.NewInternalServerError("Failed to reject node")
	}

	return c.NoContent(204)
}

func (h *NodeHandler) UpdateNodeStatus(c echo.Context) error {
	node := c.Get("node").(models.Node)

	node.LastSeen = time.Now()
	if err := h.DB.Save(&node).Error; err != nil {
		return utils.NewInternalServerError("Failed to update node status")
	}

	return c.JSON(200, node)
}

func (h *NodeHandler) ListNodes(c echo.Context) error {
	var nodes []models.Node
	if err := h.DB.Find(&nodes).Error; err != nil {
		return utils.NewInternalServerError("Failed to fetch nodes")
	}

	return c.JSON(200, nodes)
}

func (h *NodeHandler) PollInstructions(c echo.Context) error {
	node := c.Get("node").(models.Node)
	instructions := store.GlobalInstructionStore.GetInstructions(node.ID)
	return c.JSON(200, instructions)
}
