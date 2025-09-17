# Etapa de construcción
FROM golang:1.25.1-alpine AS builder

# Instalar dependencias del sistema necesarias para compilar
RUN apk add --no-cache git ca-certificates tzdata

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Instalar swag CLI (solo para generar docs si es necesario)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa de ejecución
FROM alpine:latest

# Instalar dependencias del sistema necesarias para la ejecución
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root para seguridad
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el binario desde la etapa de construcción
COPY --from=builder /app/main .

# Los archivos de documentación se generan automáticamente en el contenedor

# Configurar permisos
RUN chown -R appuser:appgroup /app

# Cambiar al usuario no-root
USER appuser

# Exponer puerto
EXPOSE 8080

# Variables de entorno por defecto
ENV PORT=8080
ENV MAX_IMAGE_SIZE=33554432
ENV MAX_BATCH_SIZE=10
ENV LOG_LEVEL=info

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Comando para ejecutar la aplicación
CMD ["./main"]
