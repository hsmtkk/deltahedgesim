package option

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Schema struct {
	Symbol      string
	SymbolName  string
	Year        int
	Month       int
	Quantity    int
	BoughtPrice int
}

const dataPath = "data/option.yaml"

func (s Schema) Save() error {
	bs, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	if err := os.WriteFile(dataPath, bs, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dataPath, err)
	}
	return nil
}

func Load() (Schema, error) {
	bs, err := os.ReadFile(dataPath)
	if err != nil {
		return Schema{}, fmt.Errorf("failed to read file %s: %w", dataPath, err)
	}
	decoded := Schema{}
	if err := yaml.Unmarshal(bs, &decoded); err != nil {
		return Schema{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return decoded, nil
}
