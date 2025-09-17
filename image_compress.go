// @title Image Compress API
// @version 1.0.0
// @description API REST para comprimir imágenes individuales y en lote, desarrollada en Go siguiendo los principios SOLID.
// @contact.name Miguel Morales
// @contact.email miguel@example.com
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/miguelmoralesr13/image-compress/internal/domain"
	"github.com/miguelmoralesr13/image-compress/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Configuración desde variables de entorno
	port := getEnv("PORT", "8080")
	maxImageSizeStr := getEnv("MAX_IMAGE_SIZE", "33554432") // 32MB por defecto
	maxBatchSizeStr := getEnv("MAX_BATCH_SIZE", "10")
	requestTimeoutStr := getEnv("REQUEST_TIMEOUT", "60")

	// Convertir valores numéricos
	maxImageSize, err := strconv.ParseInt(maxImageSizeStr, 10, 64)
	if err != nil {
		log.Fatalf("Error parseando MAX_IMAGE_SIZE: %v", err)
	}

	maxBatchSize, err := strconv.Atoi(maxBatchSizeStr)
	if err != nil {
		log.Fatalf("Error parseando MAX_BATCH_SIZE: %v", err)
	}

	requestTimeout, err := strconv.Atoi(requestTimeoutStr)
	if err != nil {
		log.Fatalf("Error parseando REQUEST_TIMEOUT: %v", err)
	}

	// Inicializar servicios (Inyección de dependencias)
	imageProcessor := services.NewImageProcessorService(maxImageSize)
	zipService := services.NewZipService()

	// Configurar router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(time.Duration(requestTimeout) * time.Second))

	// Rutas
	r.Get("/health", healthCheck)
	r.Post("/compress", compressImage(imageProcessor))
	r.Post("/compress/batch", compressBatch(imageProcessor, zipService, maxBatchSize))
	r.Post("/compress/info", getImageInfo(imageProcessor))

	// Swagger UI
	swaggerHost := getEnv("SWAGGER_HOST", "localhost:"+port)
	swaggerScheme := getEnv("SWAGGER_SCHEME", "http")
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(swaggerScheme+"://"+swaggerHost+"/swagger/doc.json"),
	))

	// Ruta para el archivo JSON de Swagger
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	// Iniciar servidor
	log.Printf("Servidor iniciado en puerto %s", port)
	log.Printf("Tamaño máximo de imagen: %d MB", maxImageSize/(1024*1024))
	log.Printf("Tamaño máximo de lote: %d imágenes", maxBatchSize)
	log.Printf("Timeout de peticiones: %d segundos", requestTimeout)
	log.Printf("NOTA: Todo el procesamiento se hace en memoria, sin archivos temporales")

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}

// healthCheck responde con el estado de la API
// @Summary Health Check
// @Description Verifica el estado de la API
// @Tags General
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Estado de la API"
// @Router /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "image-compress-api",
		"version":   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// compressImage maneja la compresión de una sola imagen
// @Summary Comprimir una imagen
// @Description Comprime una imagen individual con la calidad y formato especificados
// @Tags Compression
// @Accept multipart/form-data
// @Produce application/octet-stream
// @Param image formData file true "Archivo de imagen a comprimir"
// @Param quality formData int false "Calidad de compresión (1-100)" default(80)
// @Param format formData string false "Formato de salida (jpeg, png, webp)" Enums(jpeg, png, webp) default(jpeg)
// @Success 200 {file} file "Imagen comprimida"
// @Failure 400 {string} string "Error en la solicitud"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /compress [post]
func compressImage(processor domain.ImageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parsear multipart form
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			http.Error(w, "Error parseando formulario multipart", http.StatusBadRequest)
			return
		}

		// Obtener archivo
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error obteniendo archivo de imagen", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Leer datos del archivo
		imageData, err := io.ReadAll(file)
		if err != nil {
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
		if err := processor.ValidateImage(imageData); err != nil {
			http.Error(w, fmt.Sprintf("Imagen inválida: %v", err), http.StatusBadRequest)
			return
		}

		// Comprimir imagen
		compressedData, err := processor.CompressImage(imageData, quality, format)
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
}

// compressBatch maneja la compresión de múltiples imágenes
// @Summary Comprimir múltiples imágenes
// @Description Comprime múltiples imágenes y las devuelve en un archivo ZIP
// @Tags Compression
// @Accept json
// @Produce application/zip
// @Param request body domain.BatchCompressionRequest true "Datos de compresión en lote"
// @Success 200 {file} file "Archivo ZIP con imágenes comprimidas"
// @Failure 400 {string} string "Error en la solicitud"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /compress/batch [post]
func compressBatch(processor domain.ImageProcessor, zipService domain.ZipService, maxBatchSize int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.BatchCompressionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}

		// Validar límite de imágenes
		if len(req.Images) > maxBatchSize {
			http.Error(w, fmt.Sprintf("Máximo %d imágenes por lote", maxBatchSize), http.StatusBadRequest)
			return
		}

		// Procesar cada imagen
		files := make(map[string][]byte)
		for i, imgData := range req.Images {
			// Validar imagen
			if err := processor.ValidateImage(imgData.Data); err != nil {
				http.Error(w, fmt.Sprintf("Imagen %d inválida: %v", i+1, err), http.StatusBadRequest)
				return
			}

			// Comprimir imagen
			compressedData, err := processor.CompressImage(imgData.Data, req.Quality, req.Format)
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
		zipData, err := zipService.CreateZip(files)
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
}

// getImageInfo obtiene información de una imagen
// @Summary Obtener información de imagen
// @Description Obtiene información detallada de una imagen
// @Tags Compression
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Archivo de imagen a analizar"
// @Success 200 {object} map[string]interface{} "Información de la imagen"
// @Failure 400 {string} string "Error en la solicitud"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /compress/info [post]
func getImageInfo(processor domain.ImageProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		imageData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error leyendo datos de imagen", http.StatusBadRequest)
			return
		}

		// Validar imagen
		if err := processor.ValidateImage(imageData); err != nil {
			http.Error(w, fmt.Sprintf("Imagen inválida: %v", err), http.StatusBadRequest)
			return
		}

		// Obtener información
		width, height, format, err := processor.GetImageInfo(imageData)
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
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
