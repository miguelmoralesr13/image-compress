//go:build ignore
// +build ignore

package main

import (
	"log"

	"github.com/swaggo/swag/gen"
)

func main() {
	// Generar documentación de Swagger
	err := gen.New().Build(&gen.Config{
		SearchDir:          "./",
		OutputDir:          "./docs",
		MainAPIFile:        "image_compress.go",
		PropNamingStrategy: "snakecase",
		OutputTypes:        []string{"go", "json", "yaml"},
		ParseDependency:    true,
		ParseInternal:      false,
		MarkdownFilesDir:   "",
		GeneratedTime:      false,
	})
	if err != nil {
		log.Fatal(err)
	}
}
