package goreclass

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
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

	for _, class := range inventory.Classes {
		nodePath := filepath.Join(filepath.Dir(firstNode), class+".yml")
		// FIXME: class loop here
		subInventory, err := BuildInventory(nodePath)
		if err != nil {
			return nil, err
		}

		_ = mergo.Merge(subInventory, inventory)
		inventory = *subInventory
	}
	inventory.Classes = append(inventory.Classes, classFromFilename(firstNode))

	err = resolveReferences(&inventory, inventory.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve references: %w", err)
	}

	removeOverrideToken(inventory.Parameters)
	return &inventory, nil
}

func removeOverrideToken(node map[string]interface{}) {
	for key, rawfield := range node {
		if strings.HasPrefix(key, "~") {
			node[strings.TrimPrefix(key, "~")] = rawfield
			delete(node, key)
		}

		if next, isNode := rawfield.(map[string]interface{}); isNode {
			removeOverrideToken(next)
		}
	}
}

func resolveReferences(inventory *Inventory, node map[string]interface{}) error {
	for key, rawfield := range node {
		switch field := rawfield.(type) {
		case string:
			var err error
			node[key], err = resolveReferenceString(inventory, field)
			if err != nil {
				return err
			}

		case map[string]interface{}:
			if err := resolveReferences(inventory, field); err != nil {
				return err
			}
		}
	}
	return nil
}

var rxFullRef = regexp.MustCompile(`^\$\{((:?\w+:?)+)\}$`)
var rxRef = regexp.MustCompile(`\$\{((:?\w+:?)+)\}`)

func resolveReferenceString(inventory *Inventory, ref string) (interface{}, error) {
	var lastErr error

	if rxFullRef.MatchString(ref) {
		bytes := rxFullRef.FindString(ref)
		str := strings.TrimPrefix(strings.TrimSuffix(bytes, "}"), "${")
		resolved, err := resolveReference(inventory.Parameters, strings.Split(str, ":"))
		if err != nil {
			return nil, err
		}

		return resolved, nil
	}

	res := rxRef.ReplaceAllFunc([]byte(ref), func(bytes []byte) []byte {
		str := strings.TrimPrefix(strings.TrimSuffix(string(bytes), "}"), "${")
		resolved, err := resolveReference(inventory.Parameters, strings.Split(str, ":"))
		if err != nil {
			lastErr = err
			return nil
		}

		if _, isString := resolved.(string); isString {
			return []byte(fmt.Sprint(resolved))
		}

		lastErr = fmt.Errorf("string with reference must refer only strings")
		return nil
	})
	return string(res), lastErr
}

func resolveReference(node interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return node, nil
	}

	switch node.(type) {
	case map[string]interface{}:
		next, exists := node.(map[string]interface{})[path[0]]
		if !exists {
			return nil, fmt.Errorf("invalid reference")
		}
		return resolveReference(next, path[1:])

	default:
		return nil, fmt.Errorf("invalid reference")
	}

	return nil, nil
}

func classFromFilename(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), ".yml")
}
