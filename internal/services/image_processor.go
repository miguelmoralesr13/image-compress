package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/miguelmoralesr13/image-compress/internal/domain"
)

// ImageProcessorService implementa la interfaz ImageProcessor
type ImageProcessorService struct {
	maxImageSize int64
}

// NewImageProcessorService crea una nueva instancia del servicio
func NewImageProcessorService(maxImageSize int64) *ImageProcessorService {
	return &ImageProcessorService{
		maxImageSize: maxImageSize,
	}
}

// CompressImage comprime una imagen con la calidad especificada
func (s *ImageProcessorService) CompressImage(imageData []byte, quality int, format domain.ImageFormat) ([]byte, error) {
	if len(imageData) == 0 {
		return nil, domain.ErrEmptyImageData
	}

	if quality < 1 || quality > 100 {
		return nil, domain.ErrInvalidQuality
	}

	// Decodificar la imagen usando solo librerías estándar
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("error decodificando imagen: %w", err)
	}

	// Crear buffer para la imagen comprimida
	var buf bytes.Buffer

	// Comprimir según el formato
	switch format {
	case domain.JPEG:
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	case domain.PNG:
		// PNG no tiene calidad, pero podemos optimizar el tamaño
		err = png.Encode(&buf, img)
	case domain.WEBP:
		// Para WEBP necesitaríamos una librería adicional como go-webp
		// Por ahora usamos JPEG como fallback
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	default:
		// Formato por defecto: JPEG
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	}

	if err != nil {
		return nil, fmt.Errorf("error codificando imagen: %w", err)
	}

	return buf.Bytes(), nil
}

// ValidateImage valida que los datos de imagen sean válidos
func (s *ImageProcessorService) ValidateImage(imageData []byte) error {
	if len(imageData) == 0 {
		return domain.ErrEmptyImageData
	}

	if int64(len(imageData)) > s.maxImageSize {
		return domain.ErrImageTooLarge
	}

	// Intentar decodificar la imagen para validar
	_, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return domain.ErrInvalidImageData
	}

	return nil
}

// GetImageInfo obtiene información de la imagen
func (s *ImageProcessorService) GetImageInfo(imageData []byte) (width, height int, format domain.ImageFormat, err error) {
	if len(imageData) == 0 {
		return 0, 0, "", domain.ErrEmptyImageData
	}

	// Decodificar la imagen
	img, formatStr, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return 0, 0, "", err
	}

	// Obtener dimensiones
	bounds := img.Bounds()
	width = bounds.Dx()
	height = bounds.Dy()

	// Convertir formato
	format = s.convertFormat(formatStr)

	return width, height, format, nil
}

// convertFormat convierte el formato de string a ImageFormat
func (s *ImageProcessorService) convertFormat(formatStr string) domain.ImageFormat {
	switch formatStr {
	case "jpeg", "jpg":
		return domain.JPEG
	case "png":
		return domain.PNG
	case "webp":
		return domain.WEBP
	default:
		return domain.JPEG // Formato por defecto
	}
}
