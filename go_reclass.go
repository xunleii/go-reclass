package goreclass

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type (
	// Inventory represents the reclass-like inventory, built by
	// BuildInventory.
	Inventory struct {
		Classes    []string `yaml:"classes"`
		Parameters map[string]interface{}
	}
)

// BuildInventory builds a reclass-like inventory using the first node file
// given as parameter.
func BuildInventory(firstNode string) (*Inventory, error) {
	file, err := os.Open(firstNode)
	if err != nil {
		return nil, fmt.Errorf("failed to read first node '%s': %w", firstNode, err)
	}

	var inventory Inventory
	err = yaml.NewDecoder(file).Decode(&inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall first node '%s': %w", firstNode, err)
	}

	inventory.Classes = append(inventory.Classes, classFromFilename(firstNode))

	return &inventory, nil
}

func classFromFilename(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), ".yml")
}