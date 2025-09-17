package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/miguelmoralesr13/image-compress/internal/domain"
)

// CompressionHandler maneja las solicitudes de compresión de imágenes
type CompressionHandler struct {
	imageProcessor domain.ImageProcessor
	storage        domain.ImageStorage
	zipService     domain.ZipService
	validator      *validator.Validate
}

// NewCompressionHandler crea una nueva instancia del handler
func NewCompressionHandler(
	imageProcessor domain.ImageProcessor,
	storage domain.ImageStorage,
	zipService domain.ZipService,
) *CompressionHandler {
	return &CompressionHandler{
		imageProcessor: imageProcessor,
		storage:        storage,
		zipService:     zipService,
		validator:      validator.New(),
	}
}

// CompressImage maneja la compresión de una sola imagen
func (h *CompressionHandler) CompressImage(w http.ResponseWriter, r *http.Request) {
	// Parsear multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		http.Error(w, "Error parseando formulario multipart", http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error obteniendo archivo de imagen", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Leer datos del archivo
	imageData := make([]byte, header.Size)
	if _, err := file.Read(imageData); err != nil {
		http.Error(w, "Error leyendo datos de imagen", http.StatusBadRequest)
		return
	}

	// Obtener parámetros
	qualityStr := r.FormValue("quality")
	quality := 80 // Calidad por defecto
	if qualityStr != "" {
		if q, err := strconv.Atoi(qualityStr); err == nil {
			quality = q
		}
	}

	formatStr := r.FormValue("format")
	format := domain.JPEG // Formato por defecto
	if formatStr != "" {
		format = domain.ImageFormat(formatStr)
	}

	// Validar imagen
	if err := h.imageProcessor.ValidateImage(imageData); err != nil {
		http.Error(w, fmt.Sprintf("Imagen inválida: %v", err), http.StatusBadRequest)
		return
	}

	// Comprimir imagen
	compressedData, err := h.imageProcessor.CompressImage(imageData, quality, format)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error comprimiendo imagen: %v", err), http.StatusInternalServerError)
		return
	}

	// Configurar headers para descarga
	filename := fmt.Sprintf("compressed_%d.%s", time.Now().Unix(), format)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(compressedData)))

	// Escribir datos comprimidos
	if _, err := w.Write(compressedData); err != nil {
		http.Error(w, "Error escribiendo respuesta", http.StatusInternalServerError)
		return
	}
}

// CompressBatch maneja la compresión de múltiples imágenes
func (h *CompressionHandler) CompressBatch(w http.ResponseWriter, r *http.Request) {
	var req domain.BatchCompressionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
		return
	}

	// Validar request
	if err := h.validator.Struct(req); err != nil {
		http.Error(w, fmt.Sprintf("Datos inválidos: %v", err), http.StatusBadRequest)
		return
	}

	// Validar límite de imágenes
	if len(req.Images) > 10 {
		http.Error(w, "Máximo 10 imágenes por lote", http.StatusBadRequest)
		return
	}

	// Procesar cada imagen
	files := make(map[string][]byte)
	for i, imgData := range req.Images {
		// Validar imagen
		if err := h.imageProcessor.ValidateImage(imgData.Data); err != nil {
			http.Error(w, fmt.Sprintf("Imagen %d inválida: %v", i+1, err), http.StatusBadRequest)
			return
		}

		// Comprimir imagen
		compressedData, err := h.imageProcessor.CompressImage(imgData.Data, req.Quality, req.Format)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error comprimiendo imagen %d: %v", i+1, err), http.StatusInternalServerError)
			return
		}

		// Agregar al mapa de archivos
		filename := imgData.Filename
		if filename == "" {
			filename = fmt.Sprintf("image_%d.%s", i+1, req.Format)
		}
		files[filename] = compressedData
	}

	// Crear archivo ZIP
	zipData, err := h.zipService.CreateZip(files)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creando ZIP: %v", err), http.StatusInternalServerError)
		return
	}

	// Configurar headers para descarga
	zipFilename := fmt.Sprintf("compressed_batch_%d.zip", time.Now().Unix())
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilename))
	w.Header().Set("Content-Length", strconv.Itoa(len(zipData)))

	// Escribir datos ZIP
	if _, err := w.Write(zipData); err != nil {
		http.Error(w, "Error escribiendo respuesta", http.StatusInternalServerError)
		return
	}
}

// GetImageInfo obtiene información de una imagen
func (h *CompressionHandler) GetImageInfo(w http.ResponseWriter, r *http.Request) {
	// Parsear multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Error parseando formulario multipart", http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error obteniendo archivo de imagen", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Leer datos del archivo
	imageData := make([]byte, header.Size)
	if _, err := file.Read(imageData); err != nil {
		http.Error(w, "Error leyendo datos de imagen", http.StatusBadRequest)
		return
	}

	// Validar imagen
	if err := h.imageProcessor.ValidateImage(imageData); err != nil {
		http.Error(w, fmt.Sprintf("Imagen inválida: %v", err), http.StatusBadRequest)
		return
	}

	// Obtener información
	width, height, format, err := h.imageProcessor.GetImageInfo(imageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obteniendo información: %v", err), http.StatusInternalServerError)
		return
	}

	// Crear respuesta
	response := map[string]interface{}{
		"filename": header.Filename,
		"width":    width,
		"height":   height,
		"format":   format,
		"size":     len(imageData),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registra las rutas del handler
func (h *CompressionHandler) RegisterRoutes(r chi.Router) {
	r.Route("/compress", func(r chi.Router) {
		r.Post("/", h.CompressImage)
		r.Post("/batch", h.CompressBatch)
		r.Post("/info", h.GetImageInfo)
	})
}
