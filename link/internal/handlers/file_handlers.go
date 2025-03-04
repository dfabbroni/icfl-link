package handlers

import (
	"path/filepath"
	"strings"

	"link/internal/utils"

	"github.com/labstack/echo/v4"
)

type FileHandler struct{}

func (h *FileHandler) DownloadFile(c echo.Context) error {
	path := c.QueryParam("path")
	if path == "" {
		return utils.NewBadRequestError("File path is required")
	}

	if !strings.HasPrefix(path, "uploads/") {
		return utils.NewBadRequestError("Invalid file path")
	}

	fullPath := filepath.Clean(path)

	return c.File(fullPath)
}
