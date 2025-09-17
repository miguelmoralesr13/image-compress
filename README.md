# Image Compress API

[![Go Version](https://img.shields.io/badge/Go-1.25.1-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![API](https://img.shields.io/badge/API-REST-orange.svg)](http://localhost:8080/swagger/)

API REST para comprimir imágenes individuales y en lote, desarrollada en Go siguiendo los principios SOLID.

> 🚀 **Docker Ready**: Incluye configuración completa de Docker para desarrollo y producción
> 📚 **Documentación**: Swagger UI integrado para documentación interactiva
> ⚡ **Rápido**: Procesamiento eficiente de imágenes con Go
> 🔒 **Sin almacenamiento**: Las imágenes se procesan y devuelven inmediatamente, no se almacenan

## 📑 Índice

- [Características](#características)
- [Instalación y Ejecución](#-instalación-y-ejecución)
- [Documentación Swagger](#-documentación-swagger)
- [Uso de la API](#uso-de-la-api)
- [Configuración](#️-configuración)
- [Arquitectura](#arquitectura)
- [Comandos Docker](#-comandos-docker)
- [Desarrollo](#️-desarrollo)
- [Troubleshooting](#-troubleshooting)
- [Próximas mejoras](#-próximas-mejoras)

## Características

- ✅ Compresión de imágenes individuales (devuelve inmediatamente)
- ✅ Compresión en lote (múltiples imágenes en ZIP, devuelve inmediatamente)
- ✅ Soporte para formatos: JPEG, PNG, WEBP
- ✅ Validación de archivos y parámetros
- ✅ API REST con documentación integrada
- ✅ Arquitectura limpia siguiendo principios SOLID
- ✅ Manejo de errores robusto
- ✅ CORS habilitado
- ✅ **Sin almacenamiento**: Las imágenes se procesan en memoria y se devuelven al cliente
- ✅ **Procesamiento temporal**: Solo usa directorio temporal para archivos ZIP

## 🚀 Instalación y Ejecución

### Opción 1: Docker (Recomendado)

#### Desarrollo
```bash
# Clonar el repositorio
git clone <repository-url>
cd image-compress

# Ejecutar con Docker Compose (desarrollo)
docker-compose up --build

# O construir y ejecutar manualmente
docker build -t image-compress:dev .
docker run -p 8080:8080 image-compress:dev
```

#### Producción
```bash
# Ejecutar con Docker Compose (producción con Swagger)
docker-compose -f docker-compose.prod.yml up --build -d

# O construir y ejecutar manualmente
docker build -f Dockerfile.prod -t image-compress:prod .
docker run -p 8080:8080 image-compress:prod
```

### Opción 2: Ejecución Local

1. **Clona el repositorio:**
```bash
git clone <repository-url>
cd image-compress
```

2. **Instala las dependencias:**
```bash
go mod tidy
go mod download
```

3. **Instala Swagger CLI (opcional):**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. **Genera documentación Swagger (opcional):**
```bash
# Solo si quieres generar Swagger localmente
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g image_compress.go -o ./docs
```

5. **Ejecuta el servidor:**
```bash
go run image_compress.go
```

El servidor se iniciará en `http://localhost:8080`

## 📚 Documentación Swagger

La API incluye documentación interactiva con Swagger UI:

- **Swagger UI**: `http://localhost:8080/swagger/`
- **Documentación JSON**: `http://localhost:8080/swagger/doc.json`

> **Nota**: La documentación Swagger solo está disponible en la versión de producción.

### Generar documentación Swagger

La documentación Swagger se genera automáticamente durante el build de Docker. Si necesitas generarla manualmente:

```bash
# Instalar swag CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Generar documentación
swag init -g image_compress.go -o ./docs
```

> **Nota**: Los archivos en `./docs/` se generan automáticamente y están en `.gitignore`

## Uso de la API

### 1. Comprimir una imagen individual

**Endpoint:** `POST /compress`

**Descripción:** Comprime una imagen y la devuelve inmediatamente como descarga. **No se almacena en el servidor.**

**Parámetros:**
- `image`: Archivo de imagen (multipart/form-data)
- `quality`: Calidad de compresión (1-100, opcional, default: 80)
- `format`: Formato de salida (jpeg, png, webp, opcional, default: jpeg)

**Respuesta:** Archivo de imagen comprimida (descarga directa)

**Ejemplo con curl:**
```bash
curl -X POST \
  -F "image=@/path/to/image.jpg" \
  -F "quality=70" \
  -F "format=jpeg" \
  http://localhost:8080/compress \
  --output compressed_image.jpg
```

### 2. Comprimir múltiples imágenes (lote)

**Endpoint:** `POST /compress/batch`

**Body (JSON):**
```json
{
  "images": [
    {
      "filename": "image1.jpg",
      "data": "base64_encoded_data"
    },
    {
      "filename": "image2.png", 
      "data": "base64_encoded_data"
    }
  ],
  "quality": 80,
  "format": "jpeg"
}
```

**Ejemplo con curl:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "images": [
      {
        "filename": "test1.jpg",
        "data": "'$(base64 -w 0 /path/to/image1.jpg)'"
      }
    ],
    "quality": 70,
    "format": "jpeg"
  }' \
  http://localhost:8080/compress/batch \
  --output batch_compressed.zip
```

### 3. Obtener información de una imagen

**Endpoint:** `POST /compress/info`

**Parámetros:**
- `image`: Archivo de imagen (multipart/form-data)

**Ejemplo con curl:**
```bash
curl -X POST \
  -F "image=@/path/to/image.jpg" \
  http://localhost:8080/compress/info
```

**Respuesta:**
```json
{
  "filename": "image.jpg",
  "width": 1920,
  "height": 1080,
  "format": "jpeg",
  "size": 2048576
}
```

### 4. Health Check

**Endpoint:** `GET /health`

**Respuesta:**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "image-compress-api",
  "version": "1.0.0"
}
```

### 5. Información de la API

**Endpoint:** `GET /`

Devuelve documentación completa de la API con todos los endpoints disponibles.

## ⚙️ Configuración

### Variables de entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `PORT` | Puerto del servidor | `8080` |
| `TEMP_DIR` | Directorio temporal para archivos | `/tmp/image-compress` |
| `MAX_IMAGE_SIZE` | Tamaño máximo de imagen en bytes | `33554432` (32MB) |
| `MAX_BATCH_SIZE` | Número máximo de imágenes por lote | `10` |
| `LOG_LEVEL` | Nivel de logging | `info` |
| `SWAGGER_HOST` | Host para Swagger UI | `localhost:8080` |
| `SWAGGER_SCHEME` | Esquema para Swagger UI | `http` |
| `DEFAULT_QUALITY` | Calidad por defecto | `80` |
| `DEFAULT_FORMAT` | Formato por defecto | `jpeg` |
| `REQUEST_TIMEOUT` | Timeout de peticiones en segundos | `60` |
| `UPLOAD_TIMEOUT` | Timeout de subida en segundos | `300` |

### Configuración con archivo .env

Crea un archivo `.env` basado en `env.example`:

```bash
cp env.example .env
# Edita .env con tus valores
```

### Límites por defecto

- Tamaño máximo de imagen: 32MB
- Máximo de imágenes por lote: 10
- Timeout de request: 60 segundos

## Arquitectura

El proyecto sigue los principios SOLID:

- **S** - Single Responsibility: Cada clase tiene una responsabilidad específica
- **O** - Open/Closed: Abierto para extensión, cerrado para modificación
- **L** - Liskov Substitution: Las implementaciones son intercambiables
- **I** - Interface Segregation: Interfaces específicas y cohesivas
- **D** - Dependency Inversion: Dependencias a través de abstracciones

### Estructura del proyecto

```
├── internal/
│   ├── domain/           # Entidades y interfaces del dominio
│   │   ├── image.go      # Estructuras de datos y interfaces
│   │   └── errors.go     # Errores del dominio
│   ├── services/         # Implementaciones de servicios
│   │   ├── image_processor.go  # Procesamiento de imágenes
│   │   ├── storage.go          # Almacenamiento temporal
│   │   └── zip_service.go      # Creación de archivos ZIP
│   └── handlers/         # Handlers HTTP
│       ├── compression_handler.go  # Endpoints de compresión
│       └── health_handler.go       # Health check y documentación
├── image_compress.go     # Punto de entrada principal
├── go.mod               # Dependencias
└── README.md           # Este archivo
```

## Librerías utilizadas

- **Chi**: Router HTTP ligero y rápido
- **imaging**: Procesamiento de imágenes
- **validator**: Validación de datos
- **uuid**: Generación de identificadores únicos
- **cors**: Manejo de CORS

## 🐳 Comandos Docker

### Desarrollo

```bash
# Construir imagen de desarrollo
docker build -t image-compress:dev .

# Ejecutar contenedor de desarrollo
docker run -p 8080:8080 image-compress:dev

# Ejecutar con variables de entorno personalizadas
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e MAX_IMAGE_SIZE=67108864 \
  -e MAX_BATCH_SIZE=20 \
  image-compress:dev

# Ejecutar con archivo .env
docker run -p 8080:8080 --env-file .env image-compress:dev
```

### Producción

```bash
# Construir imagen de producción (con Swagger)
docker build -f Dockerfile.prod -t image-compress:prod .

# Ejecutar contenedor de producción
docker run -p 8080:8080 image-compress:prod

# Ejecutar en modo detached
docker run -d -p 8080:8080 --name image-compress image-compress:prod
```

### Docker Compose

```bash
# Desarrollo
docker-compose up --build

# Producción
docker-compose -f docker-compose.prod.yml up --build -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down

# Limpiar contenedores e imágenes
docker-compose down --rmi all --volumes --remove-orphans
```

### Comandos de desarrollo local

```bash
# Instalar dependencias
go mod tidy
go mod download

# Generar documentación Swagger
swag init -g image_compress.go -o ./docs

# Ejecutar la aplicación
go run image_compress.go

# Construir binario
go build -o image-compress image_compress.go

# Ejecutar tests
go test -v ./...

# Verificar código
go vet ./...
go fmt ./...
```

## 🛠️ Desarrollo

### Flujo de desarrollo recomendado

1. **Usar Docker para desarrollo:**
```bash
# Clonar y ejecutar
git clone <repository-url>
cd image-compress
docker-compose up --build
```

2. **Desarrollo local (si prefieres):**
```bash
# Instalar dependencias
go mod tidy

# Ejecutar la aplicación
go run image_compress.go
```

### Scripts de inicio

El proyecto incluye scripts de inicio para diferentes sistemas operativos:

**Linux/macOS:**
```bash
./scripts/start.sh
```

**Windows PowerShell:**
```powershell
.\scripts\start.ps1
```

### Estructura de archivos Docker

- `Dockerfile` - Imagen para desarrollo (sin Swagger)
- `Dockerfile.prod` - Imagen para producción (con Swagger)
- `docker-compose.yml` - Configuración de desarrollo
- `docker-compose.prod.yml` - Configuración de producción
- `env.example` - Variables de entorno de ejemplo

## 🔧 Troubleshooting

### Problemas comunes

**Error al construir imagen Docker:**
```bash
# Limpiar caché de Docker
docker system prune -f
docker-compose build --no-cache
```

**Error de permisos en directorio temporal:**
```bash
# Verificar permisos del directorio temporal
ls -la /tmp/image-compress
# O cambiar el directorio temporal
export TEMP_DIR=/path/to/your/temp/dir
```

**API no responde:**
```bash
# Verificar que el contenedor esté ejecutándose
docker ps
docker logs <container-name>

# Verificar health check
curl http://localhost:8080/health
```

### Ejemplos de uso completos

**Compresión de imagen con diferentes calidades:**
```bash
# Calidad alta (90)
curl -X POST -F "image=@image.jpg" -F "quality=90" \
  http://localhost:8080/compress -o high_quality.jpg

# Calidad media (60)
curl -X POST -F "image=@image.jpg" -F "quality=60" \
  http://localhost:8080/compress -o medium_quality.jpg

# Calidad baja (30)
curl -X POST -F "image=@image.jpg" -F "quality=30" \
  http://localhost:8080/compress -o low_quality.jpg
```

**Compresión en diferentes formatos:**
```bash
# Convertir a PNG
curl -X POST -F "image=@image.jpg" -F "format=png" \
  http://localhost:8080/compress -o converted.png

# Convertir a WebP
curl -X POST -F "image=@image.jpg" -F "format=webp" \
  http://localhost:8080/compress -o converted.webp
```

## 📋 Próximas mejoras

- [ ] Soporte para más formatos de imagen
- [ ] Compresión con diferentes algoritmos
- [ ] Redimensionamiento de imágenes
- [ ] Filtros y efectos
- [ ] Cache de imágenes procesadas
- [ ] Métricas y monitoreo
- [ ] Tests unitarios e integración
- [ ] API de métricas y estadísticas
- [ ] Soporte para procesamiento asíncrono
