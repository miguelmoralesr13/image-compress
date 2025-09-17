package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthHandler maneja las solicitudes de health check
type HealthHandler struct{}

// NewHealthHandler crea una nueva instancia del handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck responde con el estado de la API
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "image-compress-api",
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAPIInfo devuelve información sobre la API
func (h *HealthHandler) GetAPIInfo(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"name":        "Image Compress API",
		"version":     "1.0.0",
		"description": "API para comprimir imágenes individuales y en lote",
		"endpoints": map[string]interface{}{
			"POST /compress": map[string]interface{}{
				"description": "Comprime una sola imagen",
				"parameters": map[string]interface{}{
					"image":   "Archivo de imagen (multipart/form-data)",
					"quality": "Calidad de compresión (1-100, opcional, default: 80)",
					"format":  "Formato de salida (jpeg, png, webp, opcional, default: jpeg)",
				},
			},
			"POST /compress/batch": map[string]interface{}{
				"description": "Comprime múltiples imágenes y las devuelve en un ZIP",
				"parameters": map[string]interface{}{
					"images":  "Array de objetos con filename y data (JSON)",
					"quality": "Calidad de compresión (1-100)",
					"format":  "Formato de salida (jpeg, png, webp, opcional, default: jpeg)",
				},
			},
			"POST /compress/info": map[string]interface{}{
				"description": "Obtiene información de una imagen",
				"parameters": map[string]interface{}{
					"image": "Archivo de imagen (multipart/form-data)",
				},
			},
			"GET /health": map[string]interface{}{
				"description": "Health check de la API",
			},
			"GET /": map[string]interface{}{
				"description": "Información de la API",
			},
		},
		"supported_formats": []string{"jpeg", "png", "webp"},
		"max_image_size":    "32MB",
		"max_batch_size":    10,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registra las rutas del handler
func (h *HealthHandler) RegisterRoutes(r interface{}) {
	// Implementación simplificada para evitar dependencia circular
	// En una implementación real, se usaría chi.Router
}
