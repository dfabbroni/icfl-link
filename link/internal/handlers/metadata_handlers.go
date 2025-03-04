package handlers

import (
	"fmt"
	"link/internal/models"
	"link/internal/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MetadataHandler struct {
	DB *gorm.DB
}

func (h *MetadataHandler) RegisterMetadata(c echo.Context) error {
	node := c.Get("node").(models.Node)

	metadata := new(models.Metadata)
	if err := c.Bind(metadata); err != nil {
		return utils.NewBadRequestError("Invalid request payload")
	}

	metadata.NodeID = node.ID

	fmt.Println("metadata:", metadata)

	if err := h.DB.Create(metadata).Error; err != nil {
		return utils.NewInternalServerError("Failed to register metadata")
	}

	return c.JSON(200, metadata)
}

func (h *MetadataHandler) FetchMetadata(c echo.Context) error {
	var metadata []models.Metadata
	if err := h.DB.Find(&metadata).Error; err != nil {
		return utils.NewInternalServerError("Failed to fetch metadata")
	}

	return c.JSON(200, metadata)
}
