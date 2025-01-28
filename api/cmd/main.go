package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

const (
	templatePath = "cmd/template.yaml"
	outputPath   = "gen.api.yaml"
)

var componentPaths = []string{"api/base.yaml", "api/user.yaml", "api/admin.yaml"}

func main() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	templateDoc := loadOpenAPIDoc(loader, templatePath)
	templateDoc.Paths = &openapi3.Paths{}

	for _, path := range componentPaths {
		log.Printf("Merging file: %s", path)
		mergePaths(templateDoc, loadOpenAPIDoc(loader, path))
	}

	templateDoc.InternalizeRefs(context.Background(), nil)
	saveMergedSpec(templateDoc, outputPath)
	fmt.Printf("Merged OpenAPI spec saved to %s\n", outputPath)
}

func loadOpenAPIDoc(loader *openapi3.Loader, path string) *openapi3.T {
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load file %s: %v", path, err)
	}
	return doc
}

func mergePaths(target, source *openapi3.T) {
	for path, pathItem := range source.Paths.Map() {
		if _, exists := target.Paths.Map()[path]; exists {
			log.Printf("Conflict: Path %s already exists in target, overwriting", path)
		}
		target.Paths.Set(path, pathItem)
	}
}

func saveMergedSpec(doc *openapi3.T, outputPath string) {
	// Struct is required to marshall the spec to yaml with fields in correct order
	spec := struct {
		Openapi    string               `yaml:"openapi"`
		Info       *openapi3.Info       `yaml:"info"`
		Paths      *openapi3.Paths      `yaml:"paths"`
		Components *openapi3.Components `yaml:"components,omitempty"`
	}{
		Openapi:    doc.OpenAPI,
		Info:       doc.Info,
		Paths:      doc.Paths,
		Components: doc.Components,
	}

	data, err := yaml.Marshal(spec)
	if err != nil {
		log.Fatalf("Failed to marshal merged spec: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Fatalf("Failed to write merged spec to %s: %v", outputPath, err)
	}
}
