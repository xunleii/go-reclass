package goreclass

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
func BuildInventory(filepath string) (*Inventory, error) {
	return nil, nil
}
