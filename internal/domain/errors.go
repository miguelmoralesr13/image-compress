package domain

import "errors"

// Errores del dominio
var (
	ErrInvalidImageFormat = errors.New("formato de imagen no válido")
	ErrImageTooLarge      = errors.New("imagen demasiado grande")
	ErrInvalidQuality     = errors.New("calidad de compresión inválida")
	ErrEmptyImageData     = errors.New("datos de imagen vacíos")
	ErrUnsupportedFormat  = errors.New("formato de imagen no soportado")
	ErrBatchSizeExceeded  = errors.New("tamaño del lote excedido")
	ErrInvalidImageData   = errors.New("datos de imagen inválidos")
)
