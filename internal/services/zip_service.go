package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
)

// ZipService implementa la interfaz ZipService
type ZipService struct{}

// NewZipService crea una nueva instancia del servicio ZIP
func NewZipService() *ZipService {
	return &ZipService{}
}

// CreateZip crea un archivo ZIP con los archivos proporcionados
func (s *ZipService) CreateZip(files map[string][]byte) ([]byte, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no hay archivos para comprimir")
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Agregar cada archivo al ZIP
	for filename, data := range files {
		// Sanitizar el nombre del archivo
		sanitizedFilename := s.sanitizeFilename(filename)

		// Crear entrada en el ZIP
		writer, err := zipWriter.Create(sanitizedFilename)
		if err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("error creando entrada ZIP para %s: %w", filename, err)
		}

		// Escribir datos del archivo
		if _, err := writer.Write(data); err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("error escribiendo datos para %s: %w", filename, err)
		}
	}

	// Cerrar el writer del ZIP
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error cerrando ZIP: %w", err)
	}

	return buf.Bytes(), nil
}

// sanitizeFilename limpia el nombre del archivo para evitar problemas de seguridad
func (s *ZipService) sanitizeFilename(filename string) string {
	// Obtener solo el nombre del archivo sin la ruta
	baseName := filepath.Base(filename)

	// Reemplazar caracteres problemáticos
	baseName = strings.ReplaceAll(baseName, "..", "")
	baseName = strings.ReplaceAll(baseName, "/", "_")
	baseName = strings.ReplaceAll(baseName, "\\", "_")
	baseName = strings.ReplaceAll(baseName, ":", "_")
	baseName = strings.ReplaceAll(baseName, "*", "_")
	baseName = strings.ReplaceAll(baseName, "?", "_")
	baseName = strings.ReplaceAll(baseName, "\"", "_")
	baseName = strings.ReplaceAll(baseName, "<", "_")
	baseName = strings.ReplaceAll(baseName, ">", "_")
	baseName = strings.ReplaceAll(baseName, "|", "_")

	// Si el nombre está vacío después de la limpieza, usar un nombre por defecto
	if baseName == "" || baseName == "." {
		baseName = "image"
	}

	return baseName
}
