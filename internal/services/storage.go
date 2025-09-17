package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/miguelmoralesr13/image-compress/internal/domain"
)

// FileStorage implementa la interfaz ImageStorage usando el sistema de archivos
type FileStorage struct {
	tempDir string
}

// NewFileStorage crea una nueva instancia del almacenamiento
func NewFileStorage(tempDir string) *FileStorage {
	return &FileStorage{
		tempDir: tempDir,
	}
}

// SaveImage guarda una imagen en el directorio temporal
func (s *FileStorage) SaveImage(data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", domain.ErrEmptyImageData
	}

	// Crear directorio temporal si no existe
	if err := os.MkdirAll(s.tempDir, 0755); err != nil {
		return "", fmt.Errorf("error creando directorio temporal: %w", err)
	}

	// Generar nombre único para el archivo
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg" // Extensión por defecto
	}

	uniqueFilename := fmt.Sprintf("%s_%d%s",
		uuid.New().String()[:8],
		time.Now().Unix(),
		ext)

	filePath := filepath.Join(s.tempDir, uniqueFilename)

	// Escribir archivo
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return filePath, nil
}

// GetImage lee una imagen del almacenamiento
func (s *FileStorage) GetImage(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %w", err)
	}
	return data, nil
}

// DeleteImage elimina una imagen del almacenamiento
func (s *FileStorage) DeleteImage(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error eliminando archivo: %w", err)
	}
	return nil
}

// CleanupOldFiles elimina archivos antiguos del directorio temporal
func (s *FileStorage) CleanupOldFiles(maxAge time.Duration) error {
	files, err := os.ReadDir(s.tempDir)
	if err != nil {
		return fmt.Errorf("error leyendo directorio temporal: %w", err)
	}

	now := time.Now()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			filePath := filepath.Join(s.tempDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				// Log error pero continuar con otros archivos
				continue
			}
		}
	}

	return nil
}
