package domain

// Este archivo define las estructuras de datos y interfaces del dominio

// ImageFormat representa los formatos de imagen soportados
type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
	WEBP ImageFormat = "webp"
)

// CompressionRequest representa una solicitud de compresión
type CompressionRequest struct {
	Quality int         `json:"quality" validate:"min=1,max=100"`
	Format  ImageFormat `json:"format,omitempty"`
}

// BatchCompressionRequest representa una solicitud de compresión en lote
type BatchCompressionRequest struct {
	Images  []ImageData `json:"images" validate:"required,min=1,max=10"`
	Quality int         `json:"quality" validate:"min=1,max=100"`
	Format  ImageFormat `json:"format,omitempty"`
}

// ImageData representa los datos de una imagen
type ImageData struct {
	Filename string `json:"filename" validate:"required"`
	Data     []byte `json:"data" validate:"required"`
}

// CompressionResult representa el resultado de una compresión
type CompressionResult struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
	Size     int64  `json:"size"`
}

// BatchCompressionResult representa el resultado de una compresión en lote
type BatchCompressionResult struct {
	ZipData []byte `json:"zip_data"`
	Size    int64  `json:"size"`
}

// ImageProcessor define la interfaz para el procesamiento de imágenes
type ImageProcessor interface {
	CompressImage(imageData []byte, quality int, format ImageFormat) ([]byte, error)
	ValidateImage(imageData []byte) error
	GetImageInfo(imageData []byte) (width, height int, format ImageFormat, err error)
}

// ZipService define la interfaz para la creación de archivos ZIP
type ZipService interface {
	CreateZip(files map[string][]byte) ([]byte, error)
}
