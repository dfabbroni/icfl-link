package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func SaveUploadedFile(file *multipart.FileHeader, directory string) (string, error) {
	filename := fmt.Sprintf("%s", filepath.Base(file.Filename))
	path := filepath.Join("uploads", directory, filename)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file contents: %w", err)
	}

	return path, nil
}

func ExtractZipFile(zipFile *multipart.FileHeader, destinationDir string) error {
	src, err := zipFile.Open()
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer src.Close()

	tempZipPath := filepath.Join(destinationDir, "temp.zip")
	dst, err := os.Create(tempZipPath)
	if err != nil {
		return fmt.Errorf("failed to create temporary zip file: %w", err)
	}
	defer os.Remove(tempZipPath)
	defer dst.Close()

	// Copy zip content to temp file
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to save zip file: %w", err)
	}
	dst.Close()

	reader, err := zip.OpenReader(tempZipPath)
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}
	defer reader.Close()

	// Extract files
	for _, file := range reader.File {
		if strings.HasPrefix(filepath.Base(file.Name), ".") {
			continue
		}

		// Create full path for extracted file
		path := filepath.Join(destinationDir, file.Name)

		if file.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}
			continue
		}

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("failed to create directory structure: %w", err)
		}

		// Extract the file
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}

		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return fmt.Errorf("failed to open source file: %w", err)
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			dstFile.Close()
			srcFile.Close()
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

		dstFile.Close()
		srcFile.Close()
	}

	return nil
}
